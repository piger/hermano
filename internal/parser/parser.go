package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/anaskhan96/soup"
)

// Product is a product from the web shop; it's just a data container describing
// each chair on sale.
type Product struct {
	Name      string
	Price     string
	Link      string
	Available bool
}

// ParsePage analyse the HTML contents of the store's page and extracts all the information
// about each chair on sale.
func ParsePage(contents []byte) ([]Product, error) {
	var result []Product

	doc := soup.HTMLParse(string(contents))
	if doc.Error != nil {
		return result, doc.Error
	}

	for _, product := range doc.FindAll("li", "class", "product") {
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

		link, ok := a.Attrs()["href"]
		if !ok {
			log.Printf("cannot find href in <a> for %s", title)
			continue
		}

		bdi := product.Find("bdi")
		if bdi.Error != nil {
			log.Printf("cannot find price for %s", title)
			continue
		}

		price := bdi.Text()
		price = strings.Trim(price, "\u00a0")

		var currency string
		currency_span := bdi.Find("span")
		if currency_span.Error != nil {
			// assume euro if we fail to parse this span
			currency = "€"
		} else {
			currency = currency_span.Text()
		}
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
