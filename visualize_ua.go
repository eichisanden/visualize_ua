package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/woothee/woothee-go"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
)

var (
	simplify     = flag.Bool("s", false, "simplify output")
	rxLogPattern = regexp.MustCompile(`\[(.+)\].*USER-AGENT:(.+) SCREEN-SIZE:`)
)

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n\n[OPTIONS]\n", os.Args[0])
	flag.PrintDefaults()
}

type unitSetType map[string]struct{} // Use like Set. value is dummy.
type dataMapType map[string]map[string]int
type flare struct {
	name string
	count int
	children []flare
}

var (
	unitSet    = make(unitSetType)
	osMap      = make(dataMapType)
	brwsrPCMap = make(dataMapType)
	brwsrSPMap = make(dataMapType)
	brwsrMPMap = make(dataMapType)
	brwsrELMap = make(dataMapType)
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

func putMap(dataMap dataMapType, key string, unit string) {
	if subMap, ok1 := dataMap[key]; ok1 {
		if cnt, ok2 := subMap[unit]; ok2 {
			subMap[unit] = cnt + 1
		} else {
			subMap[unit] = 1
		}
		dataMap[key] = subMap
	} else {
		subMap := make(map[string]int)
		subMap[unit] = 1
		dataMap[key] = subMap
	}
}

func sortedList(m map[string]struct{}) []string {
	l := make([]string, len(m))
	i := 0
	for unit := range m {
		l[i] = unit
		i = i + 1
	}
	sort.Strings(l)
	return l
}

func sortedList2(m map[string]map[string]int) []string {
	l := make([]string, len(m))
	i := 0
	for unit := range m {
		l[i] = unit
		i = i + 1
	}
	sort.Strings(l)
	return l
}

func printHeader(unitList []string) {
	// Header
	fmt.Print("|Name/Keyword|Sum|Avg|")
	delim := "|:---|---:|---:|"
	for _, unit := range unitList {
		fmt.Printf("%s|", unit)
		delim = fmt.Sprintf("%s---:|", delim)
	}
	fmt.Println("")
	fmt.Println(delim)
}

func initKeywordTotalMap(unitList []string) map[string]int {
	cntMap := make(map[string]int)
	for _, unit := range unitList {
		cntMap[unit] = 0
	}
	return cntMap
}

func printData(dataMap map[string]map[string]int, unitList []string) {
	dataList := sortedList2(dataMap)

	// Calc major total
	majorTotal := 0
	for _, data := range dataList {
		for _, unit := range unitList {
			if cnt, ok := dataMap[data][unit]; ok {
				majorTotal = majorTotal + cnt
			}
		}
	}

	// Output indivisual row
	unitTotalMap := initKeywordTotalMap(unitList)
	for _, data := range dataList {
		buff := ""
		lineTotal := 0
		for _, unit := range unitList {
			if cnt, ok := dataMap[data][unit]; ok {
				lineTotal = lineTotal + cnt
				unitTotalMap[unit] = unitTotalMap[unit] + cnt
				buff = fmt.Sprintf("%s%d|", buff, cnt)
			} else {
				buff = fmt.Sprintf("%s0|", buff)
			}
		}
		avg := strconv.FormatFloat(float64(lineTotal)/float64(majorTotal)*100, 'f', 1, 64)
		fmt.Printf("|%s|%d|%s%%|%s\n", data, lineTotal, avg, buff)
	}

	// Output majorTotal and unitTotal
	fmt.Printf("|**Sum**|%d|100%%|", majorTotal)
	for _, unit := range unitList {
		fmt.Printf("%d|", unitTotalMap[unit])
	}
	fmt.Println("")
}

func print() {
	unitList := sortedList(unitSet)

	fmt.Print("# OS\n\n")
	printHeader(unitList)
	printData(osMap, unitList)

	fmt.Print("# ブラウザ(PC)\n\n")
	printHeader(unitList)
	printData(brwsrPCMap, unitList)
	fmt.Println("")

	fmt.Print("# ブラウザ(スマートフォン)\n\n")
	printHeader(unitList)
	printData(brwsrSPMap, unitList)
	fmt.Println("")

	fmt.Print("# ブラウザ(ガラケー)\n\n")
	printHeader(unitList)
	printData(brwsrMPMap, unitList)
	fmt.Println("")

	fmt.Print("# ブラウザ(その他)\n\n")
	printHeader(unitList)
	printData(brwsrELMap, unitList)
}

func printGraph() error {
	unitList := sortedList(unitSet)

	f, err := os.Create("data.js")
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	template := `
	{
		"label": "%s",
		"value": %d,
		"color": "#2383c1"
	},`

	for _, unit := range unitList {
		fmt.Fprintf(w, "var os_%s = [", unit)
		for k, v := range osMap {
			if c, ok := v[unit]; ok {
				fmt.Fprintf(w, template, k, c)
			}
		}
		fmt.Fprintf(w, "];\n")
	}
	w.Flush()

	return nil
}

func process() error {
	flag.Usage = usage
	flag.Parse()

	r := bufio.NewReader(os.Stdin)

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if matches := rxLogPattern.FindStringSubmatch(line); matches != nil {
			unit := matches[1]
			ua := matches[2]

			result, err := woothee.Parse(ua)
			if err != nil {
				//log.Fatalf("Cound not parse '%s' : %s", ua, err)
				continue // Ignore err because error occurs when unknow ua come.
			}

			os := getOs(result.Os, result.OsVersion)
			browser := getBrowser(result.Name, result.Version)

			if _, ok := unitSet[unit]; !ok {
				unitSet[unit] = struct{}{}
			}

			putMap(osMap, os, unit)

			switch result.Category {
			case "pc":
				putMap(brwsrPCMap, browser, unit)
			case "smartphone":
				putMap(brwsrSPMap, browser, unit)
			case "mobilephone":
				putMap(brwsrMPMap, browser, unit)
			case "appliance":
			case "crawler":
			case "misc":
			case "UNKNOWN":
				putMap(brwsrELMap, browser, unit)
			}
		}
	}

	print()
	if err := printGraph(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := process(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
