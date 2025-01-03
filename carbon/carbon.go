package carbon

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	DefaultDateTimeFormat = "2006-01-02 15:04:05"
	DefaultDateFormat     = "2006-01-02"
	WeekStartDay          = time.Monday
	Timezone              = time.Local
)

var (
	ParseLoyouts = []string{
		DefaultDateTimeFormat,
		DefaultDateFormat,
		time.RFC3339,
	}
	mu sync.Mutex
)

func AddParseLoyouts(layouts ...string) {
	mu.Lock()
	defer mu.Unlock()
	ParseLoyouts = append(ParseLoyouts, layouts...)
}

type Carbon struct {
	t     time.Time
	Valid bool
}

func (c Carbon) MarshalJSON() ([]byte, error) {
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

func (c Carbon) String() string {
	return c.t.String()
}

func (c Carbon) GetTime() time.Time {
	return c.t
}

func (c Carbon) GetTimestamp() int64 {
	return c.t.Unix()
}

func (c Carbon) Format(format ...string) string {
	var f string
	if len(format) < 1 {
		f = DefaultDateTimeFormat
	} else {
		f = format[0]
	}

	return c.t.Format(f)
}

func (c Carbon) GetDateTimeString() string {
	return c.Format(DefaultDateTimeFormat)
}

func (c Carbon) GetDateString() string {
	return c.Format(DefaultDateFormat)
}

// ////////////////////////
// Minute
// ///////////////////////
func (c Carbon) AddMinutes(minutes int) Carbon {
	return c.Add(time.Duration(minutes) * time.Minute)
}

func (c Carbon) AddMinute() Carbon {
	return c.Add(time.Minute)
}

func (c Carbon) StartOfMinute() Carbon {
	c.t = c.t.Truncate(time.Minute)
	return c
}

func (c Carbon) EndOfMinute() Carbon {
	return c.StartOfMinute().Add(time.Minute - time.Nanosecond)
}

// ////////////////////////
// Hours
// ///////////////////////
func (c Carbon) AddHours(hours int) Carbon {
	return c.Add(time.Duration(hours) * time.Hour)
}

func (c Carbon) AddHour() Carbon {
	return c.Add(time.Hour)
}

func (c Carbon) StartOfHour() Carbon {
	c.t = c.t.Truncate(time.Hour)
	return c
}

func (c Carbon) EndOfHour() Carbon {
	return c.StartOfHour().Add(time.Hour - time.Nanosecond)
}

//////////////////////////
// Day
/////////////////////////

// AddDays modify day
func (c Carbon) AddDays(days int) Carbon {
	return c.Add(time.Duration(days) * 24 * time.Hour)
}

// AddDay add one day
func (c Carbon) AddDay() Carbon {
	return c.Add(24 * time.Hour)
}

// StartOfDay
func (c Carbon) StartOfDay() Carbon {
	year, month, day := c.t.Date()

	c.t = time.Date(year, month, day, 0, 0, 0, 0, c.t.Location())
	return c
}

// EndOfDay
func (c Carbon) EndOfDay() Carbon {
	year, month, day := c.t.Date()

	c.t = time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), c.t.Location())
	return c
}

//////////////////////////
// Week
/////////////////////////

func (c Carbon) AddWeeks(weeks int) Carbon {
	return c.AddDays(7 * weeks)
}

func (c Carbon) AddWeek() Carbon {
	return c.AddWeeks(1)
}

func (c Carbon) StartOfWeek() Carbon {
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
	c.t = t.t.AddDate(0, 0, -weekday)
	return c
}

func (c Carbon) EndOfWeek() Carbon {

	c.t = c.StartOfWeek().t.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return c
}

func (c Carbon) AddMonths(months int) Carbon {
	return c.Add(time.Duration(months) * 30 * 24 * time.Hour)
}

func (c Carbon) StartOfMonth() Carbon {
	year, month, _ := c.t.Date()
	c.t = time.Date(year, month, 1, 0, 0, 0, 0, c.t.Location())
	return c
}

func (c Carbon) EndOfMonth() Carbon {
	c.t = c.StartOfMonth().t.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return c
}

func (c Carbon) AddYears(years int) Carbon {
	return c.Add(time.Duration(years) * 365 * 24 * time.Hour)
}

func (c Carbon) StartOfYear() Carbon {
	year, _, _ := c.t.Date()
	c.t = time.Date(year, 1, 1, 0, 0, 0, 0, c.t.Location())
	return c
}

func (c Carbon) EndOfYear() Carbon {

	c.t = c.StartOfYear().t.AddDate(1, 0, 0).Add(-time.Nanosecond)
	return c
}

// ////////////////////////
// Others
// ///////////////////////
func (c Carbon) Add(sec time.Duration) Carbon {
	c.t = c.t.Add(sec)
	return c
}

func (c Carbon) Clone() Carbon {
	return New(c.t)
}

func (c *Carbon) Scan(src interface{}) (e error) {
	if bv, ok := src.([]byte); ok {
		*c, e = Parse(string(bv), Timezone)
		return
	}

	if sv, ok := src.(string); ok {
		*c, e = Parse(sv, Timezone)
		return
	}

	if ti, ok := src.(time.Time); ok {
		c.t = ti
		return
	}

	return
}

func (c *Carbon) ScanInput(data []byte) error {
	return c.Scan(data)
}

func (c Carbon) Value() (driver.Value, error) {
	return c.GetDateTimeString(), nil
}

////////////////////////
//  Public functions  //
////////////////////////

func Now(tz ...*time.Location) Carbon {
	return New(time.Now(), tz...)
}

func New(t time.Time, tz ...*time.Location) Carbon {
	if len(tz) > 0 && tz[0] != nil {
		t.In(tz[0])
	} else {
		t.In(Timezone)

	}

	return Carbon{t: t, Valid: true}
}

func Parse(value string, tz *time.Location, layout ...string) (Carbon, error) {
	value = strings.TrimSpace(value)
	if len(layout) < 1 {
		layout = ParseLoyouts
	}
	if tz == nil {
		tz = Timezone
	}

	for _, v := range layout {
		ti, err := time.ParseInLocation(v, value, tz)
		if err == nil {
			return New(ti, tz), nil
		}
	}

	return Carbon{}, errors.New("can not parse value to carbon")
}

func Today() Carbon {
	return Now().StartOfDay()
}

func Tomorrow() Carbon {
	return Now().AddDay().StartOfDay()
}
