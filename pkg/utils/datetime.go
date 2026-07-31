package utils

import (
	"fmt"
	"time"
)

const (
	DateFormat        = "2006-01-02"
	DateTimeFormat    = "2006-01-02 15:04:05"
	DateFormatJSON    = "2006-01-02T15:04:05Z07:00"
	ExcelDateFormat   = "02-Jan-2006"
)

func FormatDate(date time.Time) string {
	return date.Format(DateFormat)
}

func FormatDateTime(date time.Time) string {
	return date.Format(DateTimeFormat)
}

func ParseDate(s string) (time.Time, error) {
	return time.Parse(DateFormat, s)
}

func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(DateTimeFormat, s)
}

func FormatDateExcel(serial int) string {
	t := ExcelSerialToDate(serial)
	return t.Format(ExcelDateFormat)
}

func ExcelSerialToDate(serial int) time.Time {
	// Excel serial date: 1 = Jan 1, 1900 (with a leap year bug)
	excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	return excelEpoch.AddDate(0, 0, serial)
}

func DateToExcelSerial(date time.Time) int {
	excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
	days := int(date.Sub(excelEpoch).Hours() / 24)
	return days
}

func ToDateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func DaysBetween(a, b time.Time) int {
	return int(ToDateOnly(b).Sub(ToDateOnly(a)).Hours() / 24)
}

func ParseDateWithFormats(s string) (time.Time, error) {
	formats := []string{
		DateFormat,
		DateTimeFormat,
		time.RFC3339,
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}
