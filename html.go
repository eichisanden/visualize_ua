package main
import (
	"database/sql"
	"bufio"
	"fmt"
	"os"
)

func outputHtmlBrowser(db *sql.DB, w *bufio.Writer) error {
	var qUnit = `select distinct unit from log order by unit`
	rowsUnit, err := db.Query(qUnit)
	if err != nil {
		return err
	}
	defer rowsUnit.Close()

	fmt.Fprintln(w, `<html lang="ja"><head><meta charset="utf-8"><link rel="stylesheet" href="css/style.css"></head><body>`)
	fmt.Fprintln(w, `<table>`)
	fmt.Fprint(w, `<tr><th></th><th><a target="_blank" href="pie.html?target=b_all">all</a></th>`)
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
	fmt.Fprint(w, `<tr><th></th><th><a target="_blank" href="pie.html?target=o_all">all</a></th>`)
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
