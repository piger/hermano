package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/piger/hermano/internal/config"
	"github.com/piger/hermano/internal/notify"
	"github.com/piger/hermano/internal/parser"
)

var (
	configFilename = flag.String("config", "hermano.toml", "Path to configuration file")
	interval       = flag.Duration("interval", 30*time.Minute, "Polling interval")
)

// The URL of the store selling products.
const storeURL = "https://usedaeronireland.ie/used-herman-miller-aeron-chairs/"

// fetchProducts download the remote sale page and extract all sale information; it returns
// an array of Product objects describing each product on sale.
func fetchProducts() ([]parser.Product, error) {
	var result []parser.Product

	resp, err := http.Get(storeURL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	products, err := parser.ParsePage(contents)
	if err != nil {
		return result, err
	}

	result = append(result, products...)
	return result, nil
}

// checkPage downloads the sales page from the website and print any offer that is found.
// This function should be called at the desired interval.
func checkPage(conf *config.Config, ignored map[string]struct{}) error {
	log.Printf("checking for offers")

	products, err := fetchProducts()
	if err != nil {
		return err
	}

	for _, product := range products {
		if _, ok := ignored[product.Name]; ok {
			continue
		}
		msg := fmt.Sprintf("%s (%s) available=%v - %s", product.Name, product.Price, product.Available, product.Link)
		fmt.Println(msg)

		if product.Available && conf.CanNotify() {
			if err := notify.Notify(conf, msg); err != nil {
				log.Printf("error sending notification: %s", err)
			}
		}
	}
	fmt.Println("")

	return nil
}

func run() error {
	conf, err := config.ReadConfig(*configFilename)
	if err != nil {
		return fmt.Errorf("error reading configuration file %q: %s", *configFilename, err)
	}

	// set up signal handler and ticker
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	t := time.NewTicker(*interval)
	defer t.Stop()

	// build a map of products to ignore, for easy lookup.
	ignored := make(map[string]struct{})
	for _, ig := range conf.Ignored {
		ignored[ig] = struct{}{}
	}

	// check the page initially and then every $interval seconds.
	if err := checkPage(conf, ignored); err != nil {
		log.Println(err)
	}

Loop:
	for {
		select {
		case <-t.C:
			if err := checkPage(conf, ignored); err != nil {
				log.Println(err)
			}

		case s := <-sig:
			log.Printf("signal received: %s", s)
			break Loop
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}
