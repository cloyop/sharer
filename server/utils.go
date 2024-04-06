package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func fileIsExists(fn string) (bool, error) {
	if _, err := os.Stat(fn); err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	return true, nil
}
func makeSessionToken() string {
	s := fmt.Sprintf("%d", time.Now().Unix())
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}
func ParseArguments(args []string) (size int, port, token string, noAuth bool) {
	var skip bool
	port = "9000"
	size = 20 * 1_000_000
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}
		if arg == "-size" || arg == "-s" {
			sizeStr := args[i+1]
			siz, err := strconv.Atoi(sizeStr)
			if err != nil {
				fmt.Println(err)
				log.Fatal("invalid size")
			}
			size = siz * 1_000_000
			skip = true
			continue
		}
		if arg == "-port" || arg == "-p" {
			port = args[i+1]
			skip = true
			continue
		}
		if arg == "-token" || arg == "-t" {
			token = args[i+1]
			skip = true
			continue
		}
		if arg == "-unsecure" || arg == "-u" {
			noAuth = true
			continue
		}
	}
	if !noAuth && token == "" {
		token = makeSessionToken()
	}
	return
}
func ArgumentsGuide() string {
	return "Sharer [ client | server ] \n Server:\n  - (-port | -p) port to listen (default 9000) \n  - (-size | -s) size (MB) allow to recieve per request (default 20MB)\n  - (-token | -t) token to allow client share to you (optional)\n  - (-unsecure | -u) to not use token and allow everyone share to you\n  - example: 'sharer server -p 4000 -s 30 -t mycustomtoken"
}
