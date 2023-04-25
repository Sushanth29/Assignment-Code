package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Deal represents the lightning deal.
type Deal struct {
	ID             int       `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	ActualPrice    float64   `json:"actual_price,omitempty"`
	FinalPrice     float64   `json:"final_price,omitempty"`
	TotalUnits     int       `json:"total_units,omitempty"`
	AvailableUnits int       `json:"available_units,omitempty"`
	ExpiryTime     time.Time `json:"expiry_time,omitempty"`
}

// Order represents a customer's order.
type Order struct {
	ID         int       `json:"id,omitempty"`
	DealID     int       `json:"deal_id,omitempty"`
	UserID     int       `json:"user_id,omitempty"`
	Quantity   int       `json:"quantity,omitempty"`
	TotalPrice float64   `json:"total_price,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

// DB stores the database connection.
type DB struct {
	db *sql.DB
}

// createDealHandler creates a new lightning deal.
func (db *DB) createDealHandler(w http.ResponseWriter, r *http.Request) {
	var deal Deal
	err := json.NewDecoder(r.Body).Decode(&deal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the expiry time of the deal.
	deal.ExpiryTime = time.Now().Add(12 * time.Hour)

	// Insert the deal into the database.
	stmt, err := db.db.Prepare("INSERT INTO deals(name, actual_price, final_price, total_units, available_units, expiry_time) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(deal.Name, deal.ActualPrice, deal.FinalPrice, deal.TotalUnits, deal.AvailableUnits, deal.ExpiryTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ID of the deal and return it.
	id, _ := result.LastInsertId()
	deal.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deal)
}

// updateDealHandler updates an existing lightning deal.
func (db *DB) updateDealHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID not specified", http.StatusBadRequest)
		return
	}

	var deal Deal
	err := json.NewDecoder(r.Body).Decode(&deal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the deal in the database.
	stmt, err := db.db.Prepare("UPDATE deals SET name=?, actual_price=?, final_price=?, total_units=?, available_units=?, expiry_time=? WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(deal.Name, deal.ActualPrice, deal.FinalPrice, deal.TotalUnits, deal.AvailableUnits, deal.ExpiryTime, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        }
}
