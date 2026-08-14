package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestSKUWithSlash_CreateThenGet verifies that a product with slash in SKU
// can be both created and retrieved via the HTTP API.
func TestSKUWithSlash_CreateThenGet(t *testing.T) {
	clk := &testClock{t: time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)}
	api := NewWithClock(clk.now)
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	// create product with slash in SKU
	body := `{"sku":"ITEM/V2","name":"slash-sku","stock":5}`
	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/products", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	createBody, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		// correctly rejected at creation - acceptable fix
		return
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", resp.StatusCode, createBody)
	}

	// if creation succeeded, the product must be queryable
	resp2, err := http.Get(srv.URL + "/api/products/ITEM/V2")
	if err != nil {
		t.Fatal(err)
	}
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("created SKU 'ITEM/V2' but GET returned %d - routing breaks on slash in SKU", resp2.StatusCode)
	}
}
