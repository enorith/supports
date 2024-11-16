package dbutil

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/enorith/supports/carbon"
	jsoniter "github.com/json-iterator/go"
)

type Datetime struct {
	carbon.Carbon
}

func (c Datetime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, c.GetDateTimeString())), nil
}

func (c *Datetime) UnmarshalJSON(row []byte) error {
	s := strings.Trim(string(row), "\"")
	if s == "" || s == "0001-01-01 00:00:00" {
		return nil
	}
	cb, e := carbon.Parse(s, carbon.Timezone)
	if e == nil {
		*c = Datetime{Carbon: cb}
	}
	return e
}

func (c *Datetime) ScanInput(data []byte) error {
	return c.UnmarshalJSON(data)
}

func (c Datetime) OfDate() Date {
	return Date(c)
}

type Date struct {
	carbon.Carbon
}

func (d *Date) UnmarshalJSON(row []byte) error {
	s := strings.Trim(string(row), "\"")
	if s == "" || s == "0000-00-00" {
		return nil
	}
	t, e := time.Parse("2006-01-02", s)
	if e != nil {
		return e
	}
	*d = Date{Carbon: carbon.New(t)}
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.GetDateString())), nil
}

func (d Date) Value() (driver.Value, error) {
	return d.GetDateString(), nil
}

func (d *Date) ScanInput(data []byte) error {
	return d.UnmarshalJSON(data)
}

type WithTimestamps struct {
	CreatedAt Datetime  `gorm:"column:created_at;autoCreateTime;type:timestamp null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;type:timestamp null" json:"updated_at"`
}

type SliceString []string

// Scan assigns a value from a database driver.
// The src value will be of one of the following types:
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
//	nil - for NULL values
//
// An error should be returned if the value cannot be stored
// without loss of information.
//
// Reference types such as []byte are only valid until the next call to Scan
// and should not be retained. Their underlying memory is owned by the driver.
// If retention is necessary, copy their values before the next call to Scan.
func (ss *SliceString) Scan(src any) error {
	var val string
	if s, ok := src.(string); ok {
		val = s
	}

	if s, ok := src.([]byte); ok {
		val = string(s)
	}
	if val == "" {
		*ss = make(SliceString, 0)
	} else {
		*ss = strings.Split(val, ",")
	}

	return nil
}

func (ss SliceString) Value() (driver.Value, error) {
	return strings.Join(ss, ","), nil
}

type SliceInt []int64

func (si *SliceInt) Scan(src any) error {
	var val string
	if s, ok := src.(string); ok {
		val = s
	}
	if s, ok := src.([]byte); ok {
		val = string(s)
	}

	if val == "" {
		*si = make(SliceInt, 0)
	} else {
		strs := strings.Split(val, ",")
		sl := make(SliceInt, 0)
		for _, s := range strs {
			if s == "" {
				continue
			}
			i, _ := strconv.ParseInt(s, 10, 64)
			sl = append(sl, i)
		}
		*si = sl
	}

	return nil
}

func (si SliceInt) Value() (driver.Value, error) {
	sl := make([]string, len(si))

	for i, v := range si {
		sl[i] = strconv.FormatInt(v, 10)
	}

	return strings.Join(sl, ","), nil
}

type JsonObjString map[string]interface{}

// Scan assigns a value from a database driver.
// The src value will be of one of the following types:
//
//	int64
//	float64
//	bool
//	[]byte
//	string
//	time.Time
//	nil - for NULL values
//
// An error should be returned if the value cannot be stored
// without loss of information.
//
// Reference types such as []byte are only valid until the next call to Scan
// and should not be retained. Their underlying memory is owned by the driver.
// If retention is necessary, copy their values before the next call to Scan.
func (js *JsonObjString) Scan(src any) error {
	if src == nil {
		return nil
	}
	var val []byte
	if s, ok := src.(string); ok {
		val = []byte(s)
	}

	if s, ok := src.([]byte); ok {
		val = s
	}
	return jsoniter.Unmarshal(val, js)
}

func (js *JsonObjString) ScanInput(data []byte) error {
	if data == nil {
		return nil
	}

	return jsoniter.Unmarshal(data, js)
}

func (js JsonObjString) Value() (driver.Value, error) {
	if js == nil {
		return nil, nil
	}

	return jsoniter.Marshal(js)
}
