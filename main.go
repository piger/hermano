package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piger/hermano/internal/config"
	"github.com/piger/hermano/internal/parser"
)

var (
	configFilename = flag.String("config", "hermano.toml", "Path to configuration file")
	interval       = flag.Duration("interval", 30*time.Minute, "Polling interval")
)

const storeURL = "https://usedaeronireland.ie/used-herman-miller-aeron-chairs/"

func run() error {
	conf, err := config.ReadConfig(*configFilename)
	if err != nil {
		return fmt.Errorf("error reading configuration file %q: %s", *configFilename, err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	t := time.NewTicker(*interval)
	defer t.Stop()

	ignored := make(map[string]struct{})
	for _, ig := range conf.Ignored {
		ignored[ig] = struct{}{}
	}

	fh, err := os.Open("page.html")
	if err != nil {
		return err
	}
	defer fh.Close()

	contents, err := io.ReadAll(fh)
	if err != nil {
		return err
	}

	products, err := parser.ParsePage(contents)
	if err != nil {
		return err
	}

	for i, product := range products {
		if _, ok := ignored[product.Name]; ok {
			continue
		}
		fmt.Printf("%d: %s (%s) available=%v - %s\n", i+1, product.Name, product.Price, product.Available, product.Link)
	}

	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
