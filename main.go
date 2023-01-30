package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var revision = "unknown"

var (
	ts   = flag.String("r", "", "timestamp to convert to date")
	date = flag.String("date", "", "date to convert to timestamp")

	version = flag.Bool("version", false, "print version and exit")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Converts a timestamp to a date and vice versa.\n")

		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -r 1674898343\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -date 01.01.2018\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -date \"01.01.2018, 16:32:15\"\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -date \"01.01.2018, 16:32:15 +0400\"\n", os.Args[0])

		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Printf("ts version %s\n", revision)
		os.Exit(0)
	}

	err := run(*ts, *date, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run is the main entry point for the program.
func run(ts string, date string, w io.Writer) error {
	// convert timestamp to date
	if ts != "" {
		i, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid timestamp: %s", ts)
		}
		_, err = fmt.Fprintln(w, time.Unix(i, 0).UTC().Format(time.RFC1123))
		if err != nil {
			return fmt.Errorf("failed to write output: %s", err)
		}
		return nil
	}

	// get timestamp for current date
	if date == "" {
		_, err := fmt.Fprintln(w, time.Now().Unix())
		if err != nil {
			return fmt.Errorf("failed to write output: %s", err)
		}
		return nil
	}

	// parse and convert date to timestamp
	dt := strings.Split(date, ",")
	if len(dt) == 1 {
		// only date
		d, err := time.Parse("02.01.2006", dt[0])
		if err != nil {
			return fmt.Errorf("invalid date: %s", date)
		}
		_, err = fmt.Fprintln(w, d.Unix())
		if err != nil {
			return fmt.Errorf("failed to write output: %s", err)
		}
		return nil
	}

	// have date and time
	var (
		d   time.Time
		err error
	)
	tt := strings.Split(strings.TrimSpace(dt[1]), " ")
	if len(tt) == 2 {
		if strings.HasPrefix(tt[1], "+") || strings.HasPrefix(tt[1], "-") {
			// numeric time zone
			d, err = time.Parse("02.01.2006, 15:04:05 -0700", date)
		} else {
			// named time zone or GMT+XX
			// fallback to UTC
			tt[1] = "UTC"
			d, err = time.Parse("02.01.2006, 15:04:05 MST", date)
		}
	} else {
		// no time zone
		d, err = time.Parse("02.01.2006, 15:04:05", date)
	}
	if err != nil {
		return fmt.Errorf("invalid date: %s", date)
	}

	_, err = fmt.Fprintln(w, d.Unix())
	if err != nil {
		return fmt.Errorf("failed to write output: %s", err)
	}
	return nil
}
