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
	}
	return Holiday{}, false
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

		location := time.Now().Location()
		date, err := time.ParseInLocation("2006-01-02", parts[0], location)
		if err != nil {
			date, err = time.ParseInLocation("01-02", parts[0], location)
		}
		color := parts[1]
		message := strings.Join(parts[2:], " ")

		holidays = append(holidays, Holiday{
			Date:    date.Format("2006-01-02"),
			Color:   color,
			Message: message,
		})
	}

	return holidays, scanner.Err()
}
