package carbon_test

import (
	"testing"

	"github.com/enorith/supports/carbon"
)

func TestParse(t *testing.T) {
	c, e := carbon.Parse("2019-01-01T00:02:00Z", nil)
	if e != nil {
		t.Error(e)
	}
	t.Log(c.GetDateTimeString())
}

func TestOffset(t *testing.T) {
	c := carbon.Now()
	t.Log(c.AddHours(1))
	t.Log(c.AddMinutes(2))
	t.Log(c.StartOfHour(), c.EndOfHour())
}
