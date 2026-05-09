package timex

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04:05"

type DateTime time.Time

func Now() DateTime {
	return DateTime(time.Now())
}

func (d DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(t.Local().Format(DateTimeLayout))), nil
}

func (d *DateTime) UnmarshalJSON(data []byte) error {
	value, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}
	if value == "" {
		*d = DateTime(time.Time{})
		return nil
	}
	parsed, err := time.ParseInLocation(DateTimeLayout, value, time.Local)
	if err != nil {
		return err
	}
	*d = DateTime(parsed)
	return nil
}

func (d DateTime) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (d *DateTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*d = DateTime(time.Time{})
	case time.Time:
		*d = DateTime(v)
	case string:
		return d.scanString(v)
	case []byte:
		return d.scanString(string(v))
	case int64:
		*d = DateTime(time.UnixMilli(v))
	case int:
		*d = DateTime(time.UnixMilli(int64(v)))
	default:
		return fmt.Errorf("unsupported DateTime value %T", value)
	}
	return nil
}

func (d *DateTime) scanString(value string) error {
	if value == "" {
		*d = DateTime(time.Time{})
		return nil
	}
	if parsed, err := time.ParseInLocation(DateTimeLayout, value, time.Local); err == nil {
		*d = DateTime(parsed)
		return nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return err
	}
	*d = DateTime(parsed)
	return nil
}
