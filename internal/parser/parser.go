package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/anaskhan96/soup"
)

type Product struct {
	Name      string
	Price     string
	Link      string
	Available bool
}

func ParsePage(contents []byte) ([]Product, error) {
	var result []Product

	doc := soup.HTMLParse(string(contents))
	if doc.Error != nil {
		return result, doc.Error
	}

	products := doc.FindAll("li", "class", "product")
	for _, product := range products {
		h2 := product.Find("h2")
		if h2.Error != nil {
			log.Printf("cannot find title (h2) in product")
			continue
		}

		title := h2.Text()

		a := product.Find("a", "class", "woocommerce-LoopProduct-link")
		if a.Error != nil {
			log.Printf("cannot find link for %s", title)
			continue
		}

		link := a.Attrs()["href"]

		bdi := product.Find("bdi")
		if bdi.Error != nil {
			log.Printf("cannot find price for %s", title)
			continue
		}

		price := bdi.Text()
		price = strings.Trim(price, "\u00a0")

		currency := bdi.Find("span").Text()
		price = fmt.Sprintf("%s %s", currency, price)

		available := true

		for _, span := range product.FindAll("span") {
			if class, ok := span.Attrs()["class"]; ok {
				if strings.Contains(class, "soldout") {
					available = false
					break
				}
			}
		}

		p := Product{
			Name:      title,
			Price:     price,
			Link:      link,
			Available: available,
		}
		result = append(result, p)
	}

	return result, nil
}
