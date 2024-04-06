package client

import (
	"log"
	"log/slog"
)

func ParseArguments(args []string) (file, address, token string) {
	var skip bool
	for i, arg := range args {
		if skip {
			skip = false
			continue
		}

		if arg == "-addr" || arg == "-a" {
			address = args[i+1]
			skip = true
			continue
		}
		if arg == "-file" || arg == "-f" {
			file = args[i+1]
			skip = true
			continue
		}
		if arg == "-token" || arg == "-t" {
			token = args[i+1]
			skip = true
			continue
		}
	}
	if token == "" {
		slog.Info("You didnt set token, continuing anyway")
	}
	if address == "" {
		log.Fatalf("missing address ( -addr | -a )")
	}
	if file == "" {
		log.Fatalf("missing file ( -file | -f )")
	}
	return
}
func ArgumentsGuide() string {
	return "Client:\n  - (-file | -f) path to file|folder\n  - (-addr | -a) set address to send\n  - (-token | -t) set token to authenticate with server (optional)\n  - example: 'sharer client -f myfile.txt -a 127.0.0.1:9000 -t myauthtoken`\n"
}
