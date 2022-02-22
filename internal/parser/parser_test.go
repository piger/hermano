package parser

import (
	"io"
	"os"
	"testing"
)

func TestParsePage(t *testing.T) {
	fh, err := os.Open("test/page.html")
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()

	contents, err := io.ReadAll(fh)
	if err != nil {
		t.Fatal(err)
	}

	products, err := ParsePage(contents)
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 7 {
		t.Fatalf("products counts should be 7, is %d", len(products))
	}

	p0 := products[0]
	if p0.Name != "Herman Miller Aeron Chair Fully Loaded (Size B)" {
		t.Fatalf("name doesn't match: %q", p0.Name)
	}

	if p0.Price != "€ 595.00" {
		t.Fatalf("price doesn't match: %q", p0.Price)
	}

	if p0.Link != "https://usedaeronireland.ie/product/herman-miller-aeron-chair-fully-loaded-size-b/" {
		t.Fatalf("link doesn't match: %q", p0.Link)
	}

	if p0.Available {
		t.Fatalf("available doesn't match: %v", p0.Available)
	}
}
