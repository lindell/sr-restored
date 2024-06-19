package scraper

import (
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var durationParsing = regexp.MustCompile(`^(\d+)(?::(\d+))? minuter$`)

func parseDuring(str string) (time.Duration, error) {
	sections := durationParsing.FindStringSubmatch(str)
	if sections == nil {
		return time.Duration(0), errors.New("could not parse duration")
	}

	minutes, _ := strconv.Atoi(sections[1])
	seconds, _ := strconv.Atoi(sections[2])

	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
