package holiday

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Holidays []Holiday

// Match attempts to match a given time with a holiday.
func (hs Holidays) Match(t time.Time) (Holiday, bool) {
	for _, h := range hs {
		if h.Date == t.Format("2006-01-02") {
			return h, true
		}
		if strings.TrimPrefix(h.Date, "0000-") == t.Format("01-02") {
			return h, true
		}
		if strings.TrimPrefix(h.Date, "0000-00-") == t.Format("02") {
			return h, true
		}
	}
	return Holiday{}, false
}

// Prefix a note if the holiday matches a given date.
func (hs Holidays) Prefix(t time.Time, note string) string {
	if h, ok := hs.Match(t); ok {
		note = h.Message + "\n\n" + note
	}
	return note
}

type Holiday struct {
	Date    string
	Color   string
	Message string
}

func Load(lists []string) Holidays {
	var holidays []Holiday
	for _, l := range lists {
		f, err := os.Open(os.ExpandEnv(l))
		if err != nil {
			log.Println(err)
			continue
		}
		defer f.Close()

		h, err := parse(f)
		if err != nil {
			log.Printf("failed parsing %v: %v\n", l, err)
		}
		holidays = append(holidays, h...)
	}
	return holidays
}

func parse(r io.Reader) ([]Holiday, error) {
	var holidays []Holiday
	scanner := bufio.NewScanner(r)
	for i := 1; scanner.Scan(); i++ {
		text := scanner.Text()
		if text == "" {
			continue
		}

		parts := strings.Split(text, " ")
		if len(parts) < 2 {
			return nil, fmt.Errorf(
				"line %v: not enough fields",
				i,
			)
		}

		date, err := parseDate(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid date %v: %v", parts[0], err)
		}
		color := parts[1]
		message := strings.Join(parts[2:], " ")

		holidays = append(holidays, Holiday{
			Date:    date,
			Color:   color,
			Message: message,
		})
	}

	return holidays, scanner.Err()
}

func parseDate(s string) (string, error) {
	location := time.Now().Location()
	date, err := time.ParseInLocation("2006-01-02", s, location)
	if err == nil {
		return date.Format("2006-01-02"), nil
	}
	date, err = time.ParseInLocation("01-02", s, location)
	if err == nil {
		// Year will be 0000.
		return date.Format("2006-01-02"), nil
	}

	date, err = time.ParseInLocation("02", s, location)
	// Year and month will be 0000-00.
	d := strings.ReplaceAll(
		date.Format("2006-01-02"),
		"0000-01",
		"0000-00",
	)
	return d, err
}
