package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"
)

func MakeSessionToken() string {
	s := fmt.Sprintf("%d%s", time.Now().Unix(), time.Now())
	rs := md5.Sum([]byte(s))
	return hex.EncodeToString(rs[:])
}

func GetIpAddr() (string, error) {
	s, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	return strings.Split(s[1].String(), "/")[0], nil
}
