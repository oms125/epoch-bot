package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/oms125/epoch-bot/bot"
	"github.com/oms125/epoch-bot/database"
)

var (
	Bot *bot.Bot
)

func main() {
	db := database.Init()

	Bot = bot.New(db)

	Bot.InitCommands()

	err := Bot.Session.Open()
	if err != nil { log.Fatal(err) }
	defer Bot.Session.Close()

	log.Println("Bot running...")
 	c := make(chan os.Signal, 1)
 	signal.Notify(c, os.Interrupt)
 	<-c
}