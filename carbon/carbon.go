package carbon

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

var (
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
	DefaultDateFormat     = "2006-01-02"
	WeekStartDay          = time.Monday
	Timezone              = time.Local
)

type Carbon struct {
	t time.Time
}

func (c *Carbon) MarshalJSON() ([]byte, error) {
	year, month, day := c.t.Date()
	hour := c.t.Hour()
	minute := c.t.Minute()
	second := c.t.Second()
	millisecond := c.t.UnixNano() / int64(time.Millisecond)

	return json.Marshal(map[string]interface{}{
		"year":      year,
		"month":     month,
		"day":       day,
		"hour":      hour,
		"minute":    minute,
		"second":    second,
		"time":      millisecond,
		"timezone":  c.t.Location().String(),
		"datetime":  c.GetDateTimeString(),
		"timestamp": c.GetTimestamp(),
	})
}

func (c *Carbon) String() string {
	return c.t.String()
}

func (c *Carbon) GetTime() time.Time {
	return c.t
}

func (c *Carbon) GetTimestamp() int64 {
	return c.t.Unix()
}

func (c *Carbon) Format(format ...string) string {
	var f string
	if len(format) < 1 {
		f = DefaultDateTimeFormat
	} else {
		f = format[0]
	}

	return c.t.Format(f)
}

func (c *Carbon) GetDateTimeString() string {
	return c.Format(DefaultDateTimeFormat)
}

func (c *Carbon) GetDateString() string {
	return c.Format(DefaultDateFormat)
}

//////////////////////////
// Minute
/////////////////////////
func (c *Carbon) AddMinutes(minutes int) *Carbon {
	return c.Add(time.Duration(minutes) * time.Minute)
}

func (c *Carbon) AddMinute() *Carbon {
	return c.Add(time.Minute)
}

func (c *Carbon) StartOfMinute() *Carbon {
	c.t = c.t.Truncate(time.Minute)
	return c
}

func (c *Carbon) EndOfMinute() *Carbon {
	return c.StartOfMinute().Add(time.Minute - time.Nanosecond)
}

//////////////////////////
// Day
/////////////////////////

// AddDays modify day
func (c *Carbon) AddDays(days int) *Carbon {
	return c.Add(time.Duration(days) * 24 * time.Hour)
}

// AddDay add one day
func (c *Carbon) AddDay() *Carbon {
	return c.Add(24 * time.Hour)
}

// StartOfDay
func (c *Carbon) StartOfDay() *Carbon {
	year, month, day := c.t.Date()

	c.t = time.Date(year, month, day, 0, 0, 0, 0, c.t.Location())
	return c
}

// EndOfDay
func (c *Carbon) EndOfDay() *Carbon {
	year, month, day := c.t.Date()

	c.t = time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), c.t.Location())
	return c
}

//////////////////////////
// Week
/////////////////////////

func (c *Carbon) AddWeeks(weeks int) *Carbon {
	return c.AddDays(7 * weeks)
}

func (c *Carbon) AddWeek() *Carbon {
	return c.AddWeeks(1)
}

func (c *Carbon) StartOfWeek() *Carbon {
	t := c.StartOfDay()
	weekday := int(t.t.Weekday())

	if WeekStartDay != time.Sunday {
		weekStartDayInt := int(WeekStartDay)

		if weekday < weekStartDayInt {
			weekday = weekday + 7 - weekStartDayInt
		} else {
			weekday = weekday - weekStartDayInt
		}
	}
	c.t = c.t.AddDate(0, 0, -weekday)
	return c
}

func (c *Carbon) EndOfWeek() *Carbon {
	c.StartOfWeek()
	c.t = c.t.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return c
}

//////////////////////////
// Others
/////////////////////////
func (c *Carbon) Add(sec time.Duration) *Carbon {
	c.t = c.t.Add(sec)
	return c
}

func (c *Carbon) Clone() *Carbon {
	return NewCarbon(c.t)
}

func (c *Carbon) Scan(src interface{}) (e error) {
	if bv, ok := src.([]byte); ok {
		c.t, e = time.Parse(DefaultDateTimeFormat, string(bv))
		return
	}

	if sv, ok := src.(string); ok {
		c.t, e = time.Parse(DefaultDateTimeFormat, sv)
		return
	}

	if ti, ok := src.(time.Time); ok {
		c.t = ti
		return
	}

	return
}

func (c Carbon) Value() (driver.Value, error) {
	return c.GetDateTimeString(), nil
}

////////////////////////
//  Public functions  //
////////////////////////

func NowCarbon(tz ...*time.Location) *Carbon {
	return NewCarbon(time.Now(), tz...)
}

func NewCarbon(t time.Time, tz ...*time.Location) *Carbon {
	if len(tz) > 0 && tz[0] != nil {
		t.In(tz[0])
	} else {
		t.In(Timezone)

	}

	return &Carbon{t}
}

func Parse(value string, tz *time.Location, layout ...string) (*Carbon, error) {
	if len(layout) < 1 {
		return nil, errors.New("no layout input")
	}

	for _, v := range layout {
		ti, err := time.Parse(v, value)
		if err == nil {
			return NewCarbon(ti, tz), nil
		}
	}
	return nil, errors.New("can not parse value to carbon")
}

func TodayCarbon() *Carbon {
	return NowCarbon().StartOfDay()
}

func TomorrowCarbon() *Carbon {
	return NowCarbon().AddDay().StartOfDay()
}
