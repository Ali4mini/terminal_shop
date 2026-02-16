package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Origin      string  `json:"origin"`
	Roast       string  `json:"roast"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

func StartMockServer() {
	products := []Product{
		{
			ID:          "1",
			Name:        "Midnight Blend",
			Origin:      "Ethiopia / Brazil",
			Roast:       "Dark",
			Price:       18.50,
			Description: "Notes of dark chocolate and toasted marshmallow. Perfect for late nights.",
		},
		{
			ID:          "2",
			Name:        "Golden Hour",
			Origin:      "Colombia",
			Roast:       "Light",
			Price:       22.00,
			Description: "Bright acidity with citrus notes and a honey-like sweetness.",
		},
		{
			ID:          "3",
			Name:        "Velvet Espresso",
			Origin:      "Guatemala",
			Roast:       "Medium",
			Price:       20.00,
			Description: "Smooth body with a nutty finish and hints of red apple.",
		},
		{
			ID:          "4",
			Name:        "Cloud Nine",
			Origin:      "Costa Rica",
			Roast:       "Medium-Light",
			Price:       24.50,
			Description: "Floral aroma with a clean, tea-like body and jasmine finish.",
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	})

	// Add a dummy checkout endpoint for later
	mux.HandleFunc("POST /checkout", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, "Order received!")
	})

	fmt.Printf("Coffee API listening on port 9991 \n")
	http.ListenAndServe("localhost:9991", mux)
}
