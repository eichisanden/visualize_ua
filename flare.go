package main
import (
	"database/sql"
	"os"
	"bufio"
	"encoding/json"
	"fmt"
)

//func outputFullJson(db *sql.DB, allCnt int) error {
//	fp, err := os.Create("graph.json")
//	if err != nil {
//		return err
//	}
//	w := bufio.NewWriter(fp)
//	fmt.Fprintf(w, `{"name": "all", "count": %d, "children": [`, allCnt)
//
//	var query = "select unit,count(*) as cnt from log group by unit"
//
//	rows, err := db.Query(query)
//	if err != nil {
//		return err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var unit string
//		var cnt int
//		rows.Scan(&unit, &cnt)
//		fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, unit, cnt)
//
//		var query = "select category,count(*) as cnt from log where unit=? group by category"
//
//		rows2, err := db.Query(query, unit)
//		if err != nil {
//			return err
//		}
//		defer rows2.Close()
//
//		var i2 int
//		for i2 = 0; rows2.Next(); i2++ {
//			var category string
//			var cnt2 int
//			rows2.Scan(&category, &cnt2)
//			fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, category, cnt2)
//
//			var query = "select os,count(*) as cnt from log where unit=? and category=? group by os"
//
//			rows3, err := db.Query(query, unit, category)
//			if err != nil {
//				return err
//			}
//			defer rows3.Close()
//
//			var i3 int
//			for i3 = 0; rows3.Next(); i3++ {
//				var os string
//				var cnt4 int
//				rows3.Scan(&os, &cnt4)
//				fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, os, cnt4)
//
//				var query = "select os_version,count(*) as cnt from log where unit=? and category=? and os = ? group by os_version"
//
//				rows4, err := db.Query(query, unit, category, os)
//				if err != nil {
//					return err
//				}
//				defer rows4.Close()
//
//				var i4 int
//				for i4 = 0; rows4.Next(); i4++ {
//					var osversion string
//					var cnt5 int
//					rows4.Scan(&osversion, &cnt5)
//					fmt.Fprintf(w, `{"browser": "%s", "count": %d, "children": [`, osversion, cnt5)
//
//					var query = "select browser,count(*) as cnt from log where unit=? and category=? and os=? and os_version=? group by browser"
//
//					rows5, err := db.Query(query, unit, category, os, osversion)
//					if err != nil {
//						return err
//					}
//					defer rows5.Close()
//
//					var i5 int
//					for i5 = 0; rows5.Next(); i5++ {
//						var browser string
//						var cnt6 int
//						rows5.Scan(&browser, &cnt6)
//						fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, browser, cnt6)
//
//						var query = "select version,count(*) as cnt from log where unit=? and category=? and os=? and os_version=? and browser=? group by version"
//
//						rows6, err := db.Query(query, unit, category, os, osversion, browser)
//						if err != nil {
//							return err
//						}
//						defer rows6.Close()
//
//						d6 := []string{}
//						var i6 int
//						for i6 := 0; rows6.Next(); i6++ {
//							var version string
//							var cnt7 int
//							rows6.Scan(&version, &cnt7)
//							d6 = append(d6, fmt.Sprintf(`{"name": "%s", "count": %d}`, version, cnt7))
//						}
//						fmt.Fprintln(w, strings.Join(d6, ","))
//						fmt.Fprint(w, "]}")
//						if i6 > 1 {
//							fmt.Fprintln(w, ",")
//						}
//					}
//					fmt.Fprint(w, "]}")
//					if i5 > 1 {
//						fmt.Fprintln(w, ",")
//					}
//				}
//				fmt.Fprint(w, "]}")
//				if i4 > 1 {
//					fmt.Fprintln(w, ",")
//				}
//			}
//			fmt.Fprint(w, "]}")
//			if i3 > 1 {
//				fmt.Fprintln(w, ",")
//			}
//		}
//		fmt.Fprint(w, "]}")
//		if i2 > 1 {
//			fmt.Fprintln(w, ",")
//		}
//	}
//	fmt.Fprintln(w, "]}")
//
//	w.Flush()
//	fp.Close()
//
//	return nil
//}

type Data struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Children []Data `json:"children"`
}

func OutputOsJson(db *sql.DB, allCnt int) error {
	fp, err := os.Create("./json/os.json")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)

	d := Data{Name: "", Count: allCnt, Children: []Data{}}

	var query = "select os_short_name,count(*) as cnt from log group by os_short_name"

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var osShortName string
		var cnt int
		rows.Scan(&osShortName, &cnt)
		d.Children = append(d.Children, Data{Name: osShortName, Count: cnt, Children: []Data{}})
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}

	fmt.Fprint(w, string(b))
	w.Flush()
	fp.Close()

	return nil
}
