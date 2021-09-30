package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"event-bot-wpp/src/auth"
	"event-bot-wpp/src/bot"
	"event-bot-wpp/src/event"

	"github.com/Rhymen/go-whatsapp"
)

func main() {
	now := uint64(time.Now().Unix())
	//create new WhatsApp connection
	wac, err := whatsapp.NewConn(5 * time.Second)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return
	}
	wac.AddHandler(&bot.WaHandler{
		C:         wac,
		Event:     event.Event{},
		StartTime: now,
	})

	err = auth.Login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return
	}

	//verifies phone connectivity
	pong, err := wac.AdminTest()

	if !pong || err != nil {
		log.Fatalf("error pinging in: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	<-time.After(3 * time.Second)
	//SAVE INFO
	// HERE
	//BEFORE
	//EXITING

	//Disconnect safe
	fmt.Println("Shutting down now.")
	session, err := wac.Disconnect()
	if err != nil {
		log.Fatalf("error disconnecting: %v\n", err)
	}
	if err := auth.WriteSession(session); err != nil {
		log.Fatalf("error saving session: %v", err)
	}

}
