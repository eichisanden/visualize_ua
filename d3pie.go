package main

// Output JSON File for d3pie chart
// http://d3pie.org/

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"math/rand"
)

type PieDataMap map[string]PieData

type PieData struct {
	SortOrder string      `json:"sort_order"`
	Content   []PieDetail `json:"content"`
}

type PieDetail struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	Color string `json:"color"`
}

type ResultSet struct {
	Name  string
	Count int
}

var colors = []string{
	"#1f77b4", "#aec7e8",
	"#ff7f0e", "#ffbb78",
	"#2ca02c", "#98df8a",
	"#d62728", "#ff9896",
	"#9467bd", "#c5b0d5",
	"#8c564b", "#c49c94",
	"#e377c2", "#f7b6d2",
	"#7f7f7f", "#c7c7c7",
	"#bcbd22", "#dbdb8d",
	"#17becf", "#9edae5",
}

func outputD3pieJson(db *sql.DB) error {

	fp, err := os.Create("./static/data/pie.js")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)

	// ----------------------------------
	// Main data map
	// ----------------------------------
	var pMap = make(PieDataMap)

	// ----------------------------------
	// OS all
	// ----------------------------------
	var p = PieData{SortOrder: "value-desc", Content: []PieDetail{}}

	var query = "select os_short_name,count(*) as cnt from log group by os_short_name"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var result = ResultSet{}
		rows.Scan(&result.Name, &result.Count)
		p.Content = append(p.Content, PieDetail{result.Name, result.Count, colors[rand.Intn(19)]})
	}
	rows.Close()

	pMap["o_all"] = p

	// ----------------------------------
	// Browser all
	// ----------------------------------
	p = PieData{SortOrder: "value-desc", Content: []PieDetail{}}

	query = "select browser_short_name,count(*) as cnt from log group by browser_short_name"
	rows, err = db.Query(query)
	if err != nil {
		return err
	}

	for rows.Next() {
		var result = ResultSet{}
		rows.Scan(&result.Name, &result.Count)
		p.Content = append(p.Content, PieDetail{result.Name, result.Count, colors[rand.Intn(19)]})
	}
	rows.Close()

	pMap["b_all"] = p

	// Output JSON File
	b, err := json.Marshal(pMap)
	if err != nil {
		return err
	}
	fmt.Fprint(w, "var data_all=", string(b))
	w.Flush()
	fp.Close()

	return nil
}
