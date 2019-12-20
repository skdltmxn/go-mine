package server

import (
	"crypto/aes"
	"crypto/rand"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"log"
	mrand "math/rand"

	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/net/crypto"
	"github.com/skdltmxn/go-mine/net/packet"
)

type LoginPlayer struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	name    string
	uuid    string
}

type LoginServer struct {
	sessMap map[*net.Session]*LoginPlayer
	tunnel  chan<- *DataTunnel
}

func NewLoginServer() *LoginServer {
	return &LoginServer{
		make(map[*net.Session]*LoginPlayer),
		getTunnelSender(),
	}
}

func (d *LoginServer) Dispatch(sess *net.Session, p *packet.Packet) bool {
	if sess.State() != net.SessionStateLogin {
		return false
	}

	switch p.Id() {
	case 0:
		d.requestEncryption(sess, p)
	case 1:
		if d.authenticate(sess, p) {
			d.loginSuccess(sess)
			d.joinGame(sess)
		}
	default:
		log.Printf("[LOGIN] Unknown packet ID: %d / %s", p.Id(), hex.EncodeToString(p.Data()))
	}

	return true
}

func (d *LoginServer) loginSuccess(sess *net.Session) {
	ctx := d.sessMap[sess]

	p := packet.NewPacket(2)
	w := packet.NewWriter(p)
	w.WriteString(ctx.uuid)
	w.WriteString(ctx.name)
	sess.SendPacket(p)
}

func (d *LoginServer) joinGame(sess *net.Session) {
	newEid := getNextEntityId()
	d.tunnel <- newDataTunnel(sess, d.sessMap[sess].name, newEid)

	joinGamePacket := packet.NewPacket(0x26)
	w := packet.NewWriter(joinGamePacket)

	w.WriteInt(newEid) // entity id
	w.WriteUbyte(GameModeCreative)
	w.WriteInt(GameDimensionOverworld)
	w.WriteLong(mrand.Int63())
	w.WriteUbyte(0)
	w.WriteString(GameLevelDefault)
	w.WriteVarint(32)  // render distance
	w.WriteBool(false) // reduced debug info
	w.WriteBool(true)  // enable respawn screen

	sess.SetState(net.SessionStateGame)
	sess.SendPacket(joinGamePacket)
}

func (d *LoginServer) requestEncryption(sess *net.Session, p *packet.Packet) {
	r := packet.NewReader(p)
	name, err := r.ReadString()
	if err != nil {
		log.Println("err:", err)
		sess.Close()
		return
	}

	rsaPrivKey, err := rsa.GenerateKey(crand.Reader, 1024)
	if err != nil {
		log.Println("rsa err:", err)
		sess.Close()
		return
	}

	rsaPubKey := &rsaPrivKey.PublicKey
	rawPubKey, _ := x509.MarshalPKIXPublicKey(rsaPubKey)
	token := make([]byte, 4)
	crand.Read(token)

	d.sessMap[sess] = &LoginPlayer{
		rsaPrivKey,
		rsaPubKey,
		name,
		"",
	}

	req := packet.NewPacket(1)
	w := packet.NewWriter(req)

	w.WriteString("")
	w.WriteVarint(len(rawPubKey))
	w.Write(rawPubKey)
	w.WriteVarint(len(token))
	w.Write(token)

	sess.SendPacket(req)
}

func (d *LoginServer) authenticate(sess *net.Session, p *packet.Packet) bool {
	r := packet.NewReader(p)
	sharedSecretLength, err := r.ReadVarint()
	if err != nil {
		log.Print("Shared secret len error", err)
		sess.Close()
		return false
	}

	sharedSecret := make([]byte, sharedSecretLength)
	r.Read(sharedSecret)

	tokenLength, _ := r.ReadVarint()
	token := make([]byte, tokenLength)
	r.Read(token)

	ctx := d.sessMap[sess]
	plainSecret := decryptWithPrivateKey(sharedSecret, ctx.privKey)
	token = decryptWithPrivateKey(token, ctx.privKey)

	rawPubKey, _ := x509.MarshalPKIXPublicKey(ctx.pubKey)
	authResult := auth("", ctx.name, plainSecret, rawPubKey)
	if authResult == nil {
		sess.Close()
		return false
	}

	ctx.uuid = authResult.Id

	// RSA key pair is no longer used
	ctx.privKey = nil
	ctx.pubKey = nil

	block, err := aes.NewCipher(plainSecret)
	if err != nil {
		log.Printf("aes.NewCipher failed: %+v", err)
		sess.Close()
		return false
	}

	encrypter := crypto.NewCFB8Encrypter(block, plainSecret)
	decrypter := crypto.NewCFB8Decrypter(block, plainSecret)
	sess.SetCryptor(encrypter, decrypter)

	return true
}

func decryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		log.Println("RSA decrypt error:", err)
	}
	return plaintext
}
