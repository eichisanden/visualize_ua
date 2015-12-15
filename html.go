package main

import (
	"bufio"
	"database/sql"
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

	var tableHead = `
	<table class="table table-striped table-condensed">
	<caption>ブラウザ利用比率</caption>
`

	fmt.Fprintln(w, tableHead)
	fmt.Fprint(w, `<tr><th></th><th><a href="javascript:void(0)" onclick="openD3pieChart('b_all');">すべて</a></th>`)
	for rowsUnit.Next() {
		var unit string
		rowsUnit.Scan(&unit)
		fmt.Fprintf(w, `<th><a href="javascript:void(0)" onclick="openD3pieChart('b_%s');">%s</a></th>`, unit, unit)
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
		fmt.Fprintf(w, `<tr><td><a href="javascript:void(0)" onclick="openD3pieChart('b_%s');">%s</a></td>`, browserShortName, browserShortName)

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
			fmt.Fprintf(w, `<td class="right">%d</td>`, osAllCount)
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
				fmt.Fprintf(w, `<td class="right">%d</td>`, unitOscount)
			}
			rowsUnitBrowserCount.Close()
		}
		fmt.Fprintln(w, "</tr>")
		rowsUnit.Close()
	}
	fmt.Fprintln(w, "</table>")
	return nil
}

func outputHtmlOs(db *sql.DB, w *bufio.Writer) error {
	var qUnit = `select distinct unit from log order by unit`
	rowsUnit, err := db.Query(qUnit)
	if err != nil {
		return err
	}
	defer rowsUnit.Close()

	var tableHead = `
	<table class="table table-striped table-condensed">
	<caption>OS利用比率</caption>
`

	fmt.Fprintln(w, tableHead)
	fmt.Fprint(w, `<tr><th></th><th><a href="javascript:void(0)" onclick="openD3pieChart('o_all');">すべて</a></th>`)
	for rowsUnit.Next() {
		var unit string
		rowsUnit.Scan(&unit)
		fmt.Fprintf(w, `<th><a href="javascript:void(0)" onclick="openD3pieChart('b_%s');">%s</a></th>`, unit, unit)
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
		fmt.Fprintf(w, `<tr><td><a href="javascript:void(0)" onclick="openD3pieChart('o_%s');">%s</a></td>`, osShortName, osShortName)

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
			fmt.Fprintf(w, `<td class="right">%d</td>`, osAllCount)
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
				fmt.Fprintf(w, `<td class="right">%d</td>`, unitOscount)
			}
			rowsUnitOsCount.Close()
		}
		fmt.Fprintln(w, "</tr>")
		rowsUnit.Close()
	}
	fmt.Fprintln(w, "</table>")
	return nil
}

func outputHtml(db *sql.DB) error {
	fp, err := os.Create("./static/index.html")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)

	var htmlHeader = `<html lang="ja">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="css/style.css">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" integrity="sha384-1q8mTJOASx8j1Au+a5WDVnPi2lkFfwwEAa8hDDdjZlpLegxhjVME1fgjWPGmkzs7" crossorigin="anonymous">
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css" integrity="sha384-fLW2N01lMqjakBkx3l/M9EahuwpSfeNvV63J5ezn3uZzapT0u7EYsXMjQV+0En5r" crossorigin="anonymous">
</head>
<body>
<div class="container">
`

	fmt.Fprintln(w, htmlHeader)

	err = outputHtmlBrowser(db, w)
	if err != nil {
		return err
	}

	err = outputHtmlOs(db, w)
	if err != nil {
		return err
	}

	var htmlFooter = `<!-- Modal -->
<div class="modal fade" id="myModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
	<div class="modal-dialog modal-lg" role="document">
		<div class="modal-content">
			<div class="modal-header">
				<button type="button" class="close" data-dismiss="modal" aria-label="Close">
					<span aria-hidden="true">&times;</span>
				</button>
				<h4 class="modal-title" id="myModalLabel">Modal title</h4>
			</div>
			<div class="modal-body">
				<div id="pieChart"></div>
			</div>
			<div class="modal-footer">
				<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
			</div>
		</div>
	</div>
</div>
<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.3/jquery.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha384-0mSbJDEHialfmuBBQP6A4Qrprq5OVfW37PRR3j5ELqxss1yVqOtnepnHVP9aJ7xS" crossorigin="anonymous"></script>
<script src="http://cdnjs.cloudflare.com/ajax/libs/d3/3.4.4/d3.min.js"></script>
<script src="js/d3pie.min.js"></script>
<script src="data/pie.js"></script>
<script src="js/index.js"></script>
</div>
</body>
</html>
`

	fmt.Fprintln(w, htmlFooter)

	w.Flush()
	fp.Close()
	return nil
}
