package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/woothee/woothee-go"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	_ "sort"
	_ "strconv"
	"strings"
	"runtime"
	"os/exec"
	"time"
	"text/template"
)

var (
	simplify     = flag.Bool("s", false, "simplify output")
	rxLogPattern = regexp.MustCompile(`\[(.+)\].*USER-AGENT:(.+) SCREEN-SIZE:`)
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n[OPTIONS]\n", os.Args[0])
	flag.PrintDefaults()
}

func getOs(os, osVersion string) string {
	if *simplify {
		if os == "Android" {
			// omit patch version ex)4.5.3 -> 4.5
			var rx = regexp.MustCompile(`^(\d{1,2}\.\d{1,2})`)
			if mVersion := rx.FindStringSubmatch(osVersion); mVersion != nil {
				os = fmt.Sprintf("Android %s", mVersion[1])
			}
		}
	} else if osVersion != woothee.ValueUnknown {
		os = fmt.Sprintf("%s %s", os, osVersion)
	}
	return os
}

func getBrowser(browser, version string) string {
	if *simplify {
		if browser == "Internet Explorer" {
			browser = fmt.Sprintf("%s%s", "IE", version)
		}
		switch browser {
		case "Chrome":
		case "Firefox":
		case "Safari":
		case "Opera":
		case "Sleipnir":
			return browser
		default:
			if version == woothee.ValueUnknown {
				return browser
			}
			return fmt.Sprintf("%s %s", browser, version)
		}

	} else if version != woothee.ValueUnknown {
		browser = fmt.Sprintf("%s %s", browser, version)
	}
	return browser
}

func outputHtmlBrowser(db *sql.DB, w *bufio.Writer) error {
	var qUnit = `select distinct unit from log order by unit`
	rowsUnit, err := db.Query(qUnit)
	if err != nil {
		return err
	}
	defer rowsUnit.Close()

	fmt.Fprintln(w, `<html lang="ja"><head><meta charset="utf-8"><link rel="stylesheet" href="style.css"></head><body>`)
	fmt.Fprintln(w, `<table>`)
	fmt.Fprint(w, "<tr><th></th><th>all</th>")
	for rowsUnit.Next() {
		var unit string
		rowsUnit.Scan(&unit)
		fmt.Fprintf(w, "<th>%s</th>", unit)
	}
	fmt.Fprintln(w, "</tr>")

	var qBrowser = `select distinct browser_short_name from log order by browser_short_name`
	rowsBrowser, err := db.Query(qBrowser)
	if err != nil {
		return err
	}
	defer rowsBrowser.Close()

	for rowsBrowser.Next() {
		var browserShortName string
		rowsBrowser.Scan(&browserShortName)
		fmt.Fprintf(w, "<tr><td>%s</td>", browserShortName)

		var qUnit2 = `select distinct unit from log order by unit`
		rowsUnit2, err := db.Query(qUnit2)
		if err != nil {
			return err
		}

		var qBrowserAll = `select count(*) count from log where browser_short_name = ?`
		rowsBrowserAll, err := db.Query(qBrowserAll, browserShortName)
		if err != nil {
			return err
		}
		for rowsBrowserAll.Next() {
			var osAllCount int
			rowsBrowserAll.Scan(&osAllCount)
			fmt.Fprintf(w, "<td>%d</td>", osAllCount)
		}
		rowsBrowserAll.Close()

		for rowsUnit2.Next() {
			var unit string
			rowsUnit2.Scan(&unit)
			var qBrowserUnit = `select count(*) count from log where unit = ? and browser_short_name = ?`
			rowsUnitBrowserCount, err := db.Query(qBrowserUnit, unit, browserShortName)
			if err != nil {
				return err
			}
			for rowsUnitBrowserCount.Next() {
				var unitOscount int
				rowsUnitBrowserCount.Scan(&unitOscount)
				fmt.Fprintf(w, "<td>%d</td>", unitOscount)
			}
			rowsUnitBrowserCount.Close()
		}
		fmt.Fprintln(w, "</tr>")
		rowsUnit.Close()
	}
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "</body></html>")
	return nil
}

func outputHtmlOs(db *sql.DB, w *bufio.Writer) error {
	var qUnit = `select distinct unit from log order by unit`
	rowsUnit, err := db.Query(qUnit)
	if err != nil {
		return err
	}
	defer rowsUnit.Close()

	fmt.Fprintln(w, `<html lang="ja"><head><meta charset="utf-8"><link rel="stylesheet" href="style.css"></head><body>`)
	fmt.Fprintln(w, `<table>`)
	fmt.Fprint(w, "<tr><th></th><th>all</th>")
	for rowsUnit.Next() {
		var unit string
		rowsUnit.Scan(&unit)
		fmt.Fprintf(w, "<th>%s</th>", unit)
	}
	fmt.Fprintln(w, "</tr>")

	var qOs = `select distinct os_short_name from log order by os_short_name`
	rowsOs, err := db.Query(qOs)
	if err != nil {
		return err
	}
	defer rowsOs.Close()

	for rowsOs.Next() {
		var osShortName string
		rowsOs.Scan(&osShortName)
		fmt.Fprintf(w, "<tr><td>%s</td>", osShortName)

		var qUnit2 = `select distinct unit from log order by unit`
		rowsUnit2, err := db.Query(qUnit2)
		if err != nil {
			return err
		}

		var qOsAll = `select count(*) count from log where os_short_name = ?`
		rowsOsAll, err := db.Query(qOsAll, osShortName)
		if err != nil {
			return err
		}
		for rowsOsAll.Next() {
			var osAllCount int
			rowsOsAll.Scan(&osAllCount)
			fmt.Fprintf(w, "<td>%d</td>", osAllCount)
		}
		rowsOsAll.Close()

		for rowsUnit2.Next() {
			var unit string
			rowsUnit2.Scan(&unit)
			var qOsUnit = `select count(*) count from log where unit = ? and os_short_name = ?`
			rowsUnitOsCount, err := db.Query(qOsUnit, unit, osShortName)
			if err != nil {
				return err
			}
			for rowsUnitOsCount.Next() {
				var unitOscount int
				rowsUnitOsCount.Scan(&unitOscount)
				fmt.Fprintf(w, "<td>%d</td>", unitOscount)
			}
			rowsUnitOsCount.Close()
		}
		fmt.Fprintln(w, "</tr>")
		rowsUnit.Close()
	}
	fmt.Fprintln(w, "</table>")
	fmt.Fprintln(w, "</body></html>")
	return nil
}

func outputHtml(db *sql.DB) error {
	fp, err := os.Create("./static/index.html")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)

	err = outputHtmlOs(db, w)
	if err != nil {
		return err
	}
	err = outputHtmlBrowser(db, w)
	if err != nil {
		return err
	}
	w.Flush()
	fp.Close()
	return nil
}

func outputFullJson(db *sql.DB, allCnt int) error {
	fp, err := os.Create("graph.json")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)
	fmt.Fprintf(w, `{"name": "all", "count": %d, "children": [`, allCnt)

	var query = "select unit,count(*) as cnt from log group by unit"

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var unit string
		var cnt int
		rows.Scan(&unit, &cnt)
		fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, unit, cnt)

		var query = "select category,count(*) as cnt from log where unit=? group by category"

		rows2, err := db.Query(query, unit)
		if err != nil {
			return err
		}
		defer rows2.Close()

		var i2 int
		for i2 = 0; rows2.Next(); i2++ {
			var category string
			var cnt2 int
			rows2.Scan(&category, &cnt2)
			fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, category, cnt2)

			var query = "select os,count(*) as cnt from log where unit=? and category=? group by os"

			rows3, err := db.Query(query, unit, category)
			if err != nil {
				return err
			}
			defer rows3.Close()

			var i3 int
			for i3 = 0; rows3.Next(); i3++ {
				var os string
				var cnt4 int
				rows3.Scan(&os, &cnt4)
				fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, os, cnt4)

				var query = "select os_version,count(*) as cnt from log where unit=? and category=? and os = ? group by os_version"

				rows4, err := db.Query(query, unit, category, os)
				if err != nil {
					return err
				}
				defer rows4.Close()

				var i4 int
				for i4 = 0; rows4.Next(); i4++ {
					var osversion string
					var cnt5 int
					rows4.Scan(&osversion, &cnt5)
					fmt.Fprintf(w, `{"browser": "%s", "count": %d, "children": [`, osversion, cnt5)

					var query = "select browser,count(*) as cnt from log where unit=? and category=? and os=? and os_version=? group by browser"

					rows5, err := db.Query(query, unit, category, os, osversion)
					if err != nil {
						return err
					}
					defer rows5.Close()

					var i5 int
					for i5 = 0; rows5.Next(); i5++ {
						var browser string
						var cnt6 int
						rows5.Scan(&browser, &cnt6)
						fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, browser, cnt6)

						var query = "select version,count(*) as cnt from log where unit=? and category=? and os=? and os_version=? and browser=? group by version"

						rows6, err := db.Query(query, unit, category, os, osversion, browser)
						if err != nil {
							return err
						}
						defer rows6.Close()

						d6 := []string{}
						var i6 int
						for i6 := 0; rows6.Next(); i6++ {
							var version string
							var cnt7 int
							rows6.Scan(&version, &cnt7)
							d6 = append(d6, fmt.Sprintf(`{"name": "%s", "count": %d}`, version, cnt7))
						}
						fmt.Fprintln(w, strings.Join(d6, ","))
						fmt.Fprint(w, "]}")
						if i6 > 1 {
							fmt.Fprintln(w, ",")
						}
					}
					fmt.Fprint(w, "]}")
					if i5 > 1 {
						fmt.Fprintln(w, ",")
					}
				}
				fmt.Fprint(w, "]}")
				if i4 > 1 {
					fmt.Fprintln(w, ",")
				}
			}
			fmt.Fprint(w, "]}")
			if i3 > 1 {
				fmt.Fprintln(w, ",")
			}
		}
		fmt.Fprint(w, "]}")
		if i2 > 1 {
			fmt.Fprintln(w, ",")
		}
	}
	fmt.Fprintln(w, "]}")

	w.Flush()
	fp.Close()

	return nil
}

type Data struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Children []Data `json:"children"`
}

func outputOsJson(db *sql.DB, allCnt int) error {
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

func process() error {
	flag.Usage = usage
	flag.Parse()

	r := bufio.NewReader(os.Stdin)

	os.Remove("./ua.db")
	db, err := sql.Open("sqlite3", "./ua.db")
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := `
	create table log (
	  id integer not null primary key,
	  unit                  text,
	  browser               text,
	  browser_short_name    text,
	  category              text,
	  os                    text,
	  os_short_name         text,
	  os_version            text,
	  platform              text,
	  version               text,
	  vendor                text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	allCnt := 0
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if matches := rxLogPattern.FindStringSubmatch(line); matches != nil {

			allCnt = allCnt + 1
			unit := matches[1]
			ua := matches[2]

			result, err := woothee.Parse(ua)
			if err != nil {
				//log.Fatalf("Cound not parse '%s' : %s", ua, err)
				continue // Ignore err because error occurs when unknow ua come.
			}

			osShortName := getOs(result.Os, result.OsVersion)
			browserShortName := getBrowser(result.Name, result.Version)

			stmt, err := tx.Prepare("insert into log(unit,browser,browser_short_name,category,os,os_short_name,os_version,platform,version,vendor) values(?,?,?,?,?,?,?,?,?,?)")
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(unit, result.Name, browserShortName, result.Category, result.Os, osShortName, result.OsVersion, result.Type, result.Version, result.Vendor)
			if err != nil {
				log.Fatal(err)
			}

		}
	}

	tx.Commit()

	//outputFullJson(db, allCnt)
	outputOsJson(db, allCnt)
	err = outputHtml(db)
	if err != nil {
		return nil
	}
	http.Handle("/", http.FileServer(http.Dir("static")))
	var url = "http://localhost:4000"

	go func() {
		if waitServer(url) && startBrowser(url) {
			log.Printf("A browser window should open. If not, please visit %s", url)
		} else {
			log.Printf("Please open your web browser and visit %s", url)
		}
	}()
	http.ListenAndServe(":4000", nil)
	return nil
}

func handler(
	w http.ResponseWriter,
	r *http.Request) {

	t, _ := template.ParseFiles("index.html")
	t.Execute(w, make(map[string]string))
}

func waitServer(url string) bool {
	tries := 20
	for tries > 0 {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
		tries--
	}
	return false
}

func startBrowser(url string) bool {
	// try to start the browser
	var args []string
	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default:
		args = []string{"xdg-open"}
	}
	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start() == nil
}

func main() {
	if err := process(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
