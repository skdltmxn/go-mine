package server

import (
	"encoding/hex"
	"log"

	"github.com/skdltmxn/go-mine/net"
	"github.com/skdltmxn/go-mine/net/packet"
)

type GamePlayer struct {
	name string
	eid  int32
}

type GameServer struct {
	sessMap map[*net.Session]*GamePlayer
}

func NewGameServer() *GameServer {
	g := &GameServer{make(map[*net.Session]*GamePlayer)}
	go g.waitForDataFromLoginServer()
	return g
}

func (g *GameServer) Dispatch(sess *net.Session, p *packet.Packet) bool {
	if sess.State() != net.SessionStateGame {
		return false
	}

	switch p.Id() {
	case 0x05:
		g.saveClientSetting(sess, p)
	case 0x0b:
		g.handlePluginMessage(p)
	default:
		log.Printf("[GAME] Unknown packet ID: %d / %+v", p.Id(), hex.EncodeToString(p.Data()))
	}

	return true
}

func (g *GameServer) saveClientSetting(sess *net.Session, p *packet.Packet) {
	r := packet.NewReader(p)
	locale, _ := r.ReadString()
	distance, _ := r.ReadByte()
	chatMode, _ := r.ReadVarint()
	chatColor, _ := r.ReadBoolean()
	displaySkinParts, _ := r.ReadUbyte()
	mainHand, _ := r.ReadVarint()

	log.Printf("[GAME] Client setting locale: %s / view distance: %d / chat: %d with color(%+v) / skin: %x / hand: %d",
		locale,
		distance,
		chatMode,
		chatColor,
		displaySkinParts,
		mainHand,
	)

	// TODO: actually save the settings
}

func (g *GameServer) handlePluginMessage(p *packet.Packet) {
	r := packet.NewReader(p)
	channel, _ := r.ReadString()

	if channel == "minecraft:brand" {
		data, _ := r.ReadString()
		log.Printf("[GAME] Plugin message ident: %s / data: %s", channel, data)
	}
}

func (g *GameServer) sendServerBrand(sess *net.Session) {
	p := packet.NewPacket(0x19)
	w := packet.NewWriter(p)

	w.WriteString("minecraft:brand")
	w.WriteString("go-mine")

	sess.SendPacket(p)
}

func (g *GameServer) waitForDataFromLoginServer() {
	ch := getTunnelReceiver()
	select {
	case data := <-ch:
		g.sessMap[data.sess] = &GamePlayer{
			data.name,
			data.eid,
		}
		go g.sendServerBrand(data.sess)
	}
}
