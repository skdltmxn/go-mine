package server

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"hash"
	"log"
	"net/http"
	"strings"
)

type hashContext struct {
	h hash.Hash
}

func (ctx *hashContext) update(data []byte) {
	ctx.h.Write(data)
}

func (ctx *hashContext) digest() string {
	return generateHash(ctx.h.Sum(nil))
}

type authResult struct {
	Id         string           `json:"id"`
	Name       string           `json:"name"`
	Properties []authProperties `json:"properties"`
}

type authProperties struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

func auth(server, name string, sharedSecret, pubkey []byte) *authResult {
	hash := &hashContext{sha1.New()}
	hash.update([]byte(server))
	hash.update(sharedSecret)
	hash.update(pubkey)
	digest := hash.digest()

	url := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s", name, digest)
	res, err := http.Get(url)
	if err != nil {
		log.Print("Auth failed ", err)
		return nil
	}

	defer res.Body.Close()

	var result authResult
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		log.Print("JSON decode failed ", err)
		return nil
	}

	if len(result.Id) != 32 {
		return nil
	}

	result.Id = fmt.Sprintf("%s-%s-%s-%s-%s",
		result.Id[:8],
		result.Id[8:12],
		result.Id[12:16],
		result.Id[16:20],
		result.Id[20:],
	)

	return &result
}

func generateHash(hash []byte) string {
	// Check for negative hashes
	negative := (hash[0] & 0x80) == 0x80
	if negative {
		hash = twosComplement(hash)
	}

	// Trim away zeroes
	res := strings.TrimLeft(fmt.Sprintf("%x", hash), "0")
	if negative {
		res = "-" + res
	}

	return res
}

// little endian
func twosComplement(p []byte) []byte {
	carry := true
	for i := len(p) - 1; i >= 0; i-- {
		p[i] = byte(^p[i])
		if carry {
			carry = p[i] == 0xff
			p[i]++
		}
	}
	return p
}
