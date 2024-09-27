package main

import (
	"camp-mgr/app/campmgr/internal/server"
	"github.com/teamgram/marmota/pkg/commands"
)

func main() {
	commands.Run(server.New())
}
