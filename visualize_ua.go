package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/woothee/woothee-go"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"regexp"
	_ "sort"
	_ "strconv"
	"database/sql"
	"log"
)

var (
	simplify     = flag.Bool("s", false, "simplify output")
	rxLogPattern = regexp.MustCompile(`\[(.+)\].*USER-AGENT:(.+) SCREEN-SIZE:`)
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n[OPTIONS]\n", os.Args[0])
	flag.PrintDefaults()
}

type dataMapType map[string]map[string]int
type flareType struct {
	name string
	count int
	children []flareType
}

var (
	unitSet    = make(map[string]struct{}) // Use like Set. value is dummy.
	osMap      = make(dataMapType)
	brwsrPCMap = make(dataMapType)
	brwsrSPMap = make(dataMapType)
	brwsrMPMap = make(dataMapType)
	brwsrELMap = make(dataMapType)
	flare      = &flareType{"os", 0, []flareType{}}
)

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
		// force add version if IE
		if browser == "Internet Explorer" {
			browser = fmt.Sprintf("%s %s", "Internet Explorer", version)
		}
	} else if version != woothee.ValueUnknown {
		browser = fmt.Sprintf("%s %s", browser, version)
	}
	return browser
}

//func putMap(dataMap dataMapType, key string, unit string) {
//	if subMap, ok1 := dataMap[key]; ok1 {
//		if cnt, ok2 := subMap[unit]; ok2 {
//			subMap[unit] = cnt + 1
//		} else {
//			subMap[unit] = 1
//		}
//		dataMap[key] = subMap
//	} else {
//		subMap := make(map[string]int)
//		subMap[unit] = 1
//		dataMap[key] = subMap
//	}
//}
//
//func sortedList(m map[string]interface{}) []string {
//	l := make([]string, len(m))
//	i := 0
//	for unit := range m {
//		l[i] = unit
//		i = i + 1
//	}
//	sort.Strings(l)
//	return l
//}
//
//func sortedList2(m map[string]map[string]int) []string {
//	l := make([]string, len(m))
//	i := 0
//	for unit := range m {
//		l[i] = unit
//		i = i + 1
//	}
//	sort.Strings(l)
//	return l
//}
//
//func printHeader(unitList []string) {
//	// Header
//	fmt.Print("|Name/Keyword|Sum|Avg|")
//	delim := "|:---|---:|---:|"
//	for _, unit := range unitList {
//		fmt.Printf("%s|", unit)
//		delim = fmt.Sprintf("%s---:|", delim)
//	}
//	fmt.Println("")
//	fmt.Println(delim)
//}
//
//func initKeywordTotalMap(unitList []string) map[string]int {
//	cntMap := make(map[string]int)
//	for _, unit := range unitList {
//		cntMap[unit] = 0
//	}
//	return cntMap
//}
//
//func printData(dataMap map[string]map[string]int, unitList []string) {
//	dataList := sortedList2(dataMap)
//
//	// Calc major total
//	majorTotal := 0
//	for _, data := range dataList {
//		for _, unit := range unitList {
//			if cnt, ok := dataMap[data][unit]; ok {
//				majorTotal = majorTotal + cnt
//			}
//		}
//	}
//
//	// Output indivisual row
//	unitTotalMap := initKeywordTotalMap(unitList)
//	for _, data := range dataList {
//		buff := ""
//		lineTotal := 0
//		for _, unit := range unitList {
//			if cnt, ok := dataMap[data][unit]; ok {
//				lineTotal = lineTotal + cnt
//				unitTotalMap[unit] = unitTotalMap[unit] + cnt
//				buff = fmt.Sprintf("%s%d|", buff, cnt)
//			} else {
//				buff = fmt.Sprintf("%s0|", buff)
//			}
//		}
//		avg := strconv.FormatFloat(float64(lineTotal)/float64(majorTotal)*100, 'f', 1, 64)
//		fmt.Printf("|%s|%d|%s%%|%s\n", data, lineTotal, avg, buff)
//	}
//
//	// Output majorTotal and unitTotal
//	fmt.Printf("|**Sum**|%d|100%%|", majorTotal)
//	for _, unit := range unitList {
//		fmt.Printf("%d|", unitTotalMap[unit])
//	}
//	fmt.Println("")
//}
//
//func print() {
//	unitList := sortedList(unitSet)
//
//	fmt.Print("# OS\n\n")
//	printHeader(unitList)
//	printData(osMap, unitList)
//
//	fmt.Print("# ブラウザ(PC)\n\n")
//	printHeader(unitList)
//	printData(brwsrPCMap, unitList)
//	fmt.Println("")
//
//	fmt.Print("# ブラウザ(スマートフォン)\n\n")
//	printHeader(unitList)
//	printData(brwsrSPMap, unitList)
//	fmt.Println("")
//
//	fmt.Print("# ブラウザ(ガラケー)\n\n")
//	printHeader(unitList)
//	printData(brwsrMPMap, unitList)
//	fmt.Println("")
//
//	fmt.Print("# ブラウザ(その他)\n\n")
//	printHeader(unitList)
//	printData(brwsrELMap, unitList)
//}
//
//func printGraph() error {
//	unitList := sortedList(unitSet)
//
//	f, err := os.Create("data.js")
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//	w := bufio.NewWriter(f)
//
//	template := `
//	{
//		"label": "%s",
//		"value": %d,
//		"color": "#2383c1"
//	},`
//
//	for _, unit := range unitList {
//		fmt.Fprintf(w, "var os_%s = [", unit)
//		for k, v := range osMap {
//			if c, ok := v[unit]; ok {
//				fmt.Fprintf(w, template, k, c)
//			}
//		}
//		fmt.Fprintf(w, "];\n")
//	}
//	w.Flush()
//
//	return nil
//}

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
	  unit      text,
	  name      text,
	  category  text,
	  os        text,
	  osversion text,
	  platform  text,
	  version   text,
	  vendor    text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

//	data := make(map[string] struct{})
//	data["name"] = "all"
//	data["count"] = 0
//	data["children"] = make([]map[string] struct{})

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

			//os := getOs(result.Os, result.OsVersion)
			//browser := getBrowser(result.Name, result.Version)

			stmt, err := tx.Prepare("insert into log(unit,name,category,os,osversion,platform,version,vendor) values(?,?,?,?,?,?,?,?)")
			if err != nil {
				return err
			}
			defer stmt.Close()

			_, err = stmt.Exec(unit, result.Name, result.Category, result.Os, result.OsVersion, result.Type, result.Version, result.Vendor)
			if err != nil {
				log.Fatal(err)
			}

//			unitSet[unit] = struct{}{} // Set dummy value
//
//			putMap(osMap, os, unit)
//
//			switch result.Category {
//			case "pc":
//				putMap(brwsrPCMap, browser, unit)
//			case "smartphone":
//				putMap(brwsrSPMap, browser, unit)
//			case "mobilephone":
//				putMap(brwsrMPMap, browser, unit)
//			case "appliance":
//			case "crawler":
//			case "misc":
//			case "UNKNOWN":
//				putMap(brwsrELMap, browser, unit)
//			}

//			data["count"] = data["count"] + 1
//			for _,  d := range data["children"] {
//				if d["name"] == unit {
//					d["count"] = d["count"] + 1
//				}
//			}

		// 	var isFoundCategory bool = false
		// 	for _, savedCategory := range flare.children {
		// 		if savedCategory.name == result.Category {
		// 			isFoundCategory = true
		// 			for _, savedOs := range savedCategory.children {
		// 				if savedOs.name == os {
		// 					var notFound = true
		// 					for _, savedBrowser := range savedOs.children {
		// 						if savedBrowser.name == browser {
		// 							savedBrowser.count = savedBrowser.count + 1
		// 							notFound = false
		// 						}
		// 					}
		// 					if notFound {
		// 						var savedBrowser = flareType{}
		// 						savedBrowser.name = browser
		// 						savedBrowser.count = 1
		// 						savedOs.children = append(savedOs.children, savedBrowser)
		// 					}
		// 				}
		// 			}
		// 		}	
		// 	}
			
		// 	if !isFoundCategory {
		// 		category.children
		// 		categoryElement := flareType{result.Category, 0, []flareType{}}
		// 		osElement := flareType{os, 0, []flareType{}}
		// 		brElement := flareType{browser, 1, []flareType{}}
		// 		osElement.children = append(osElement.children, osElement.children)
		// 		l := append(flare.children, nil)
		// 	}
		}
	}

	tx.Commit()

	fp, err := os.Create("graph.json")
	if err != nil {
		return err
	}
	w := bufio.NewWriter(fp)
	fmt.Fprintf(w, `{"name": "all", "count": %d, "children": [`, allCnt)
//	var query = `
//		select id,unit,name,category,os,osVersion,platform,version,vendor from log
//	`
	var query = `
		select unit,count(*) as cnt from log group by unit
	`
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var unit      string
		var cnt       int
		rows.Scan(&unit, &cnt)
		fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, unit, cnt)

		var query = `
		select category,count(*) as cnt from log where unit=? group by category
		`
		rows2, err := db.Query(query, unit)
		if err != nil {
			return err
		}
		defer rows2.Close()

		for rows2.Next() {
			var category string
			var cnt2 int
			rows2.Scan(&category, &cnt2)
			fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, category, cnt2)

			var query = `
				select os,count(*) as cnt from log where unit=? and category=? group by os
			`
			rows3, err := db.Query(query, unit, category)
			if err != nil {
				return err
			}
			defer rows3.Close()

			for rows3.Next() {
				var os string
				var cnt4 int
				rows3.Scan(&os, &cnt4)
				fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, os, cnt4)

				var query = `
					select osversion,count(*) as cnt from log where unit=? and category=? and os = ? group by osversion
				`
				rows4, err := db.Query(query, unit, category, os)
				if err != nil {
					return err
				}
				defer rows4.Close()

				for rows4.Next() {
					var osversion string
					var cnt5 int
					rows4.Scan(&osversion, &cnt5)
					fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, osversion, cnt5)

					var query = `
					select name,count(*) as cnt from log where unit=? and category=? and os=? and osversion=? group by name
					`

					rows5, err := db.Query(query, unit, category, os, osversion)
					if err != nil {
						return err
					}
					defer rows5.Close()

					for rows5.Next() {
						var name string
						var cnt6 int
						rows5.Scan(&name, &cnt6)
						fmt.Fprintf(w, `{"name": "%s", "count": %d, "children": [`, name, cnt6)

						var query = `
							select version,count(*) as cnt from log where unit=? and category=? and os=? and osversion=? and name=? group by version
						`

						rows6, err := db.Query(query, unit, category, os, osversion, name)
						if err != nil {
							return err
						}
						defer rows6.Close()

						for rows6.Next() {
							var version string
							var cnt7 int
							rows6.Scan(&version, &cnt7)
							fmt.Fprintf(w, `{"name": "%s", "count": %d}`, version, cnt7)
						}
						fmt.Fprintln(w, "]},")
					}
					fmt.Fprintln(w, "]},")
				}
				fmt.Fprintln(w, "]},")
			}
			fmt.Fprintln(w, "]},")
		}
		fmt.Fprintln(w, "]},")
	}
	fmt.Fprintln(w, "]}")

	w.Flush()
	fp.Close()

//	print()
//	if err := printGraph(); err != nil {
//		return err
//	}

	return nil
}

func main() {
	if err := process(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
