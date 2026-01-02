package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/oms125/epoch-bot/bot"
	"github.com/oms125/epoch-bot/commands"
)

func main() {
	log.Println("Starting Bot...")

	commands.Init()

	err := bot.Session.Open()
	if err != nil { log.Fatal(err) }
	defer bot.Session.Close()

	log.Println("Bot running...")
 	c := make(chan os.Signal, 1)
 	signal.Notify(c, os.Interrupt)
 	<-c
}