package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const localDateTimeFormat = "2006-01-02 15:04:05"

// LocalTime 自定义时间类型，JSON序列化时不返回时区信息
type LocalTime time.Time

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	if tTime.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, tTime.Format(localDateTimeFormat))), nil
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == `""` || string(data) == "null" {
		return nil
	}

	str := string(data)
	str = str[1 : len(str)-1] // 去掉引号

	parsedTime, err := time.ParseInLocation(localDateTimeFormat, str, time.Local)
	if err != nil {
		return err
	}

	*t = LocalTime(parsedTime)
	return nil
}

func (t LocalTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	if tTime.IsZero() {
		return nil, nil
	}
	return tTime, nil
}

func (t *LocalTime) Scan(value interface{}) error {
	if value == nil {
		*t = LocalTime(time.Time{})
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = LocalTime(v)
	case []byte:
		tTime, err := time.ParseInLocation(localDateTimeFormat, string(v), time.Local)
		if err != nil {
			return err
		}
		*t = LocalTime(tTime)
	case string:
		tTime, err := time.ParseInLocation(localDateTimeFormat, v, time.Local)
		if err != nil {
			return err
		}
		*t = LocalTime(tTime)
	default:
		return fmt.Errorf("cannot scan type %T into LocalTime", value)
	}

	return nil
}

// Time 转换为标准 time.Time
func (t LocalTime) Time() time.Time {
	return time.Time(t)
}

// Now 获取当前时间的 LocalTime
func Now() LocalTime {
	return LocalTime(time.Now())
}
