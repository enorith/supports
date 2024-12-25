package carbon_test

import (
	"testing"
	"time"

	"github.com/enorith/supports/carbon"
)

type TS time.Time

func TestParse(t *testing.T) {
	c, e := carbon.Parse("2019-01-01 00:02:00", nil)
	if e != nil {
		t.Error(e)
	}
	t.Log(c.GetTime())
}

func TestOffset(t *testing.T) {
	c := carbon.Now()
	t.Log(c.AddHours(1))
	t.Log(c.AddMinutes(2))
	t.Log("===========")
	t.Log(c.StartOfHour(), c.EndOfHour())
	t.Log(c.StartOfDay(), c.EndOfDay())
	t.Log(c.StartOfMonth(), c.EndOfMonth())
	t.Log(c.StartOfYear(), c.EndOfYear())
	t.Log(c.StartOfWeek(), c.EndOfWeek())
}

func TestScan(t *testing.T) {
	var c carbon.Carbon

	e := c.Scan("2019-01-01 ")
	if e != nil {
		t.Error(e)
	}

	t.Log(c)
}
