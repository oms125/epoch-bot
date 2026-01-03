package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/oms125/epoch-bot/bot"
	"github.com/oms125/epoch-bot/game"
)

var (
	Bot *bot.Bot
	Game *game.Game
)

func main() {
	//Init database
	Game = game.New()
	Game.InitTables()
	defer Game.DB.Close()

	//Init bot
	Bot = bot.New(Game)
	Bot.InitCommands()

	//Start session
	log.Println("Starting bot session...")
	err := Bot.Session.Open()
	if err != nil { 
		log.Fatal(err) 
	} else {
		log.Println("Bot session started")
	}
	defer Bot.Session.Close()

	log.Println("Bot running...")
 	c := make(chan os.Signal, 1)
 	signal.Notify(c, os.Interrupt)
 	<-c
}