package main

// Output JSON File for d3pie chart
// http://d3pie.org/

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

type PieDataMap map[string]*PieData

type PieData struct {
	SortOrder string      `json:"sortOrder"`
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

func makeD3pieJson(db *sql.DB, keyColumn, unit string) (*PieData, error) {
	var p = PieData{SortOrder: "value-desc", Content: []PieDetail{}}

	var where string = ""
	if unit != "" {
		where = fmt.Sprintf("WHERE unit = '%s'", unit)
	}

	var query = fmt.Sprintf("SELECT %s, count(*) AS count FROM log %s GROUP BY %s", keyColumn, where, keyColumn)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var result = ResultSet{}
		rows.Scan(&result.Name, &result.Count)
		p.Content = append(p.Content, PieDetail{result.Name, result.Count, colors[i % 20]})
	}

	return &p, nil
}

func makeD3pieUnitJson(db *sql.DB, pMap PieDataMap) error {
	var query = "SELECT DISTINCT unit FROM log"
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for i := 0; rows.Next(); i++ {
		var unit string
		rows.Scan(&unit)

		pMap[fmt.Sprintf("o_%s", unit)], err = makeD3pieJson(db, "os_short_name", unit)
		if err != nil {
			return err
		}

		pMap[fmt.Sprintf("b_%s", unit)], err = makeD3pieJson(db, "browser_short_name", unit)
		if err != nil {
			return err
		}
	}

	return nil
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
	pMap["o_all"], err = makeD3pieJson(db, "os_short_name", "")
	if err != nil {
		return err
	}

	// ----------------------------------
	// Browser all
	// ----------------------------------
	pMap["b_all"], err = makeD3pieJson(db, "browser_short_name", "")
	if err != nil {
		return err
	}

	// ----------------------------------
	// OS and browser by unit
	// ----------------------------------
	err = makeD3pieUnitJson(db, pMap)
	if err != nil {
		return err
	}

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
