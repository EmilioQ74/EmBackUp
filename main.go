package main

import (
	"log"

	"github.com/EmilioQ74/EmBackUp/config"
	"github.com/EmilioQ74/EmBackUp/gui"
)

func main() {
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	gui.Start()
}

// CLI
// func main() {
// 	ctx, cancel := signal.NotifyContext(context.Background(),
// 		os.Interrupt,
// 		syscall.SIGTERM,
// 	)
// 	defer cancel()

// 	cmd.ExecuteContext(ctx)
// }
