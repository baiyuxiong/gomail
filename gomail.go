package main

import (
	"flag"
	"github.com/baiyuxiong/gomail/config"
	"github.com/baiyuxiong/gomail/mail"
	"github.com/baiyuxiong/gomail/http"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println(config.VERSION)
		os.Exit(0)
	}

	//parse config file
	config.Parse(*cfg)

	go http.Start()

	go mail.Start()


	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("gomail quit.")
	os.Exit(0)
}