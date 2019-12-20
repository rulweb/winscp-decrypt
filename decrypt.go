package main

import (
	"strconv"
)

const (
	PwMagic = 0xA3
	PwFlag  = 0xFF
)

func decrypt(host, username, password string) string {
	key := username + host
	var passBytes []byte
	for i := 0; i < len(password); i++ {
		val, _ := strconv.ParseInt(string(password[i]), 16, 8)
		passBytes = append(passBytes, byte(val))
	}
	var flag byte
	flag, passBytes = decNextChar(passBytes)
	var length byte = 0
	if flag == PwFlag {
		_, passBytes = decNextChar(passBytes)

		length, passBytes = decNextChar(passBytes)
	} else {
		length = flag
	}
	toBeDeleted, passBytes := decNextChar(passBytes)
	passBytes = passBytes[toBeDeleted*2:]

	clearPass := ""
	var (
		i   byte
		val byte
	)
	for i = 0; i < length; i++ {
		val, passBytes = decNextChar(passBytes)
		clearPass += string(val)
	}

	if flag == PwFlag {
		clearPass = clearPass[len(key):]
	}
	return clearPass
}

func decNextChar(passBytes []byte) (byte, []byte) {
	if len(passBytes) <= 0 {
		return 0, passBytes
	}
	a := passBytes[0]
	b := passBytes[1]
	passBytes = passBytes[2:]
	return ^(((a << 4) + b) ^ PwMagic) & 0xff, passBytes
}
