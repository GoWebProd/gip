package fasttime

import (
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	laepoch = 946684800 + 86400*(31+29)

	daysPer400y = 365*400 + 97
	daysPer100y = 365*100 + 24
	daysPer4y   = 365*4 + 1
)

var daysInMonth = [...]int64{31, 30, 31, 30, 31, 31, 30, 31, 30, 31, 31, 29}

/* Originally released by the musl project (http://www.musl-libc.org/) under the
 * MIT license. Taken from the file /src/time/__secs_to_tm.c */
func ParseDateTime(timestamp int64) (year int64, month int64, day int64, hour int64, minute int64, seconds int64) {
	secs := timestamp - laepoch
	days := secs / 86400

	remsecs := secs % 86400
	if remsecs < 0 {
		remsecs += 86400
		days--
	}

	wday := (3 + days) % 7
	if wday < 0 {
		wday += 7
	}

	qcCycles := days / daysPer400y
	remdays := days % daysPer400y
	if remdays < 0 {
		remdays += daysPer400y
		qcCycles--
	}

	cCycles := remdays / daysPer100y
	if cCycles == 4 {
		cCycles--
	}
	remdays -= cCycles * daysPer100y

	qCycles := remdays / daysPer4y
	if qCycles == 25 {
		qCycles--
	}
	remdays -= qCycles * daysPer4y

	remyears := remdays / 365
	if remyears == 4 {
		remyears--
	}
	remdays -= remyears * 365

	var leap int64
	if remyears == 0 && (qCycles > 0 || cCycles == 0) {
		leap = 1
	}

	yday := remdays + 31 + 28 + leap
	if yday >= 365+leap {
		yday -= 365 + leap
	}

	years := remyears + 4*qCycles + 100*cCycles + 400*qcCycles

	var months int64
	for months = 0; daysInMonth[months] <= remdays; months++ {
		remdays -= daysInMonth[months]
	}

	tmYear := years + 2000
	tmMon := months + 2
	if tmMon >= 12 {
		tmMon -= 12
		tmYear++
	}
	tmMday := remdays + 1

	tmHour := remsecs / 3600
	tmMin := remsecs / 60 % 60
	tmSec := remsecs % 60

	return tmYear, tmMon, tmMday, tmHour, tmMin, tmSec
}

func ParseDate(timestamp int64) (year int64, month int64, day int64) {
	year, month, day, _, _, _ = ParseDateTime(timestamp)
	return
}

func FormatTimestampToDate(timestamp int64) string {
	return FormatDate(ParseDate(timestamp))
}

func FormatTimestampToDateTime(timestamp int64) string {
	return FormatDateTime(ParseDateTime(timestamp))
}

func FormatDate(year int64, month int64, day int64) string {
	month += 1

	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(year)))
	if month < 10 {
		sb.WriteString("0")
	}
	sb.WriteString(strconv.Itoa(int(month)))
	if day < 10 {
		sb.WriteString("0")
	}
	sb.WriteString(strconv.Itoa(int(day)))

	return sb.String()
}

func FormatDateTime(year int64, month int64, day int64, hour int64, minute int64, second int64) string {
	month++

	sb := strings.Builder{}

	sb.WriteString(strconv.Itoa(int(year)))

	if month < 10 {
		sb.WriteString("0")
	}

	sb.WriteString(strconv.Itoa(int(month)))

	if day < 10 {
		sb.WriteString("0")
	}

	sb.WriteString(strconv.Itoa(int(day)))

	if hour < 10 {
		sb.WriteString("0")
	}

	sb.WriteString(strconv.Itoa(int(hour)))

	if minute < 10 {
		sb.WriteString("0")
	}

	sb.WriteString(strconv.Itoa(int(minute)))

	if second < 10 {
		sb.WriteString("0")
	}

	sb.WriteString(strconv.Itoa(int(second)))

	return sb.String()
}

func Now() int64 {
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	if err != nil {
		return time.Now().Unix()
	}
	return tv.Sec
}

func NowNano() int64 {
	var tv syscall.Timeval
	err := syscall.Gettimeofday(&tv)
	if err != nil {
		return time.Now().UnixNano()
	}
	return tv.Nano()
}

func Round(timestamp int64, period int64) int64 {
	delta := timestamp % period
	result := timestamp - delta
	// if delta >= (period >> 1) {
	// 	result += period
	// }

	return result
}
