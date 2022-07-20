package main

import (
	"flag"

	"github.com/xraimanx/timemanagementtgbot/tg"
)

func main() {
	var token string
	const usage = "token of your telegram bot"
	flag.StringVar(&token, "token", "", usage)
	flag.StringVar(&token, "t", "", usage+" (shorthand)")
	flag.Parse()
	tg.Run(token)
}
