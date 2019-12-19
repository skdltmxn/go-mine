package server

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"log"

	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/net/crypto"
	"github.com/skdltmxn/go-mine/net/packet"
)

type GameContext struct {
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
	name    string
}

type LoginServer struct {
	sessMap map[*net.Session]*GameContext
}

func NewLoginServer() *LoginServer {
	return &LoginServer{make(map[*net.Session]*GameContext)}
}

func (d *LoginServer) Dispatch(sess *net.Session, p *packet.Packet) bool {
	if sess.State() != net.SessionStateLogin {
		return false
	}

	log.Printf("LoginServer got packet %d / %+v", p.Id(), hex.EncodeToString(p.Data()))

	packetId := p.Id()
	if packetId == 0 {
		r := packet.NewReader(p)
		name, err := r.ReadString()
		if err != nil {
			log.Println("err:", err)
			sess.Close()
			return true
		}

		rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 1024)
		if err != nil {
			log.Println("rsa err:", err)
			sess.Close()
			return true
		}

		rsaPubKey := &rsaPrivKey.PublicKey
		rawPubKey, _ := x509.MarshalPKIXPublicKey(rsaPubKey)
		token := make([]byte, 4)
		rand.Read(token)

		res := packet.NewPacket(1)
		w := packet.NewWriter(res)

		w.WriteString("")
		w.WriteVarint(len(rawPubKey))
		w.Write(rawPubKey)
		w.WriteVarint(len(token))
		w.Write(token)

		sess.SendData(res.Raw())
		d.sessMap[sess] = &GameContext{
			rsaPrivKey,
			rsaPubKey,
			name,
		}
	} else if packetId == 1 {
		r := packet.NewReader(p)
		sharedSecretLength, err := r.ReadVarint()
		if err != nil {
			log.Print("Shared secret len error", err)
			sess.Close()
			return true
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
			return true
		}

		// RSA key pair is no longer used
		ctx.privKey = nil
		ctx.pubKey = nil

		block, err := aes.NewCipher(plainSecret)
		if err != nil {
			log.Printf("aes.NewCipher failed: %+v", err)
			sess.Close()
			return true
		}

		encrypter := crypto.NewCFB8Encrypter(block, plainSecret)
		decrypter := crypto.NewCFB8Decrypter(block, plainSecret)
		sess.SetCryptor(encrypter, decrypter)

		// Login Success
		successPacket := packet.NewPacket(2)
		w := packet.NewWriter(successPacket)
		w.WriteString(authResult.Id)
		w.WriteString(authResult.Name)

		sess.SendData(successPacket.Raw())

		// Join Game
		joinGamePacket := packet.NewPacket(0x26)
		w = packet.NewWriter(joinGamePacket)
		w.WriteInt(0)
		w.WriteUbyte(1)
		w.WriteInt(0)
		w.WriteLong(0xabbccddeeff0102)
		w.WriteUbyte(0)
		w.WriteString("default")
		w.WriteVarint(32)
		w.WriteBool(false)
		w.WriteBool(true)

		sess.SendData(joinGamePacket.Raw())
	}

	return true
}

func decryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		log.Println("RSA decrypt error:", err)
	}
	return plaintext
}
