package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/woothee/woothee-go"
	"io"
	"log"
	"os"
	"regexp"
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

			_, err = stmt.Exec(unit, result.Name, browserShortName, result.Category, result.Os, osShortName, result.OsVersion, result.Type, result.Version, result.Vendor)
			if err != nil {
				log.Fatal(err)
			}
			stmt.Close()
		}
	}

	tx.Commit()

	if err = outputHtml(db); err != nil {
		return err
	}
	if err = outputD3pieJson(db); err != nil {
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
