package timex

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"
)

// DateTimeLayout 是前后端约定的秒级日期时间格式。
const DateTimeLayout = "2006-01-02 15:04:05"

// DateTime 统一数据库时间和 JSON 时间格式，避免接口直接暴露 RFC3339。
type DateTime time.Time

// Now 返回当前本地时间的 DateTime 表示。
func Now() DateTime {
	return DateTime(time.Now())
}

// MarshalJSON 将时间输出为 DateTimeLayout，零值输出为空字符串以兼容表单展示。
func (d DateTime) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(strconv.Quote(t.Local().Format(DateTimeLayout))), nil
}

// UnmarshalJSON 按 DateTimeLayout 解析前端提交的时间字符串。
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

// Value 实现 driver.Valuer，零值时间写入数据库时按 NULL 处理。
func (d DateTime) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// Scan 实现 sql.Scanner，兼容 PostgreSQL 时间、字符串和毫秒时间戳输入。
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
	// 数据库驱动或外部数据可能返回 RFC3339Nano，保留兜底解析以兼容历史数据。
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return err
	}
	*d = DateTime(parsed)
	return nil
}
