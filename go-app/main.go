package main

import (
	"database/sql"
	"fmt"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"log"
	_ "time"
)

type Currency struct {
	Cur_ID           int     `json:"Cur_ID"`
	Date             string  `json:"Date"`
	Cur_Abbreviation string  `json:"Cur_Abbreviation"`
	Cur_Scale        int     `json:"Cur_Scale"`
	Cur_Name         string  `json:"Cur_Name"`
	Cur_OfficialRate float64 `json:"Cur_OfficialRate"`
}

func fetchRates() ([]Currency, error) {
	client := resty.New()
	resp, err := client.R().SetResult(&[]Currency{}).Get("https://api.nbrb.by/exrates/rates?periodicity=0")
	if err != nil {
		return nil, err
	}

	rates := *resp.Result().(*[]Currency)
	return rates, nil
}

func saveToDB(db *sql.DB, rates []Currency) error {
	query := `INSERT INTO currency_rates (cur_id, date, cur_abbreviation, cur_scale, cur_name, cur_official_rate) 
              VALUES (?, ?, ?, ?, ?, ?)`

	for _, rate := range rates {
		_, err := db.Exec(query, rate.Cur_ID, rate.Date, rate.Cur_Abbreviation, rate.Cur_Scale, rate.Cur_Name, rate.Cur_OfficialRate)
		if err != nil {
			return err
		}
	}

	return nil
}

func processing() {
	dsn := "user:password@tcp(localhost:3306)/currency_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	rates, err := fetchRates()
	if err != nil {
		log.Fatal(err)
	}

	err = saveToDB(db, rates)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("The data has been successfully saved to the database.")
}

func main() {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", processing)
	if err != nil {
		return
	}
	c.Start()
}
