package basemoji

import (
	"testing"
)

func TestRoundtrip(t *testing.T) {
	tcs := []string{
		"",
		"this is some test case",
		"اللُّغَة العَرَبِيّة هي أكثرُ اللغاتِ السامية تحدثاً، وإحدى أكثر اللغات انتشاراً في العالم، يتحدثُها أكثرُ من",
		"漢語，又稱中文、唐話[2]，或被視為一個語族，或被視為隸屬於漢藏語系漢語族之一種語言。",
	}
	e := NewEncoding(StdEncoding)
	for i, tc := range tcs {
		b, err := e.DecodeString(e.EncodeToString([]byte(tc)))
		if err != nil {
			t.Errorf("TestRoundtrip %d: error %v", i, err)
		}
		if string(b) != tc {
			t.Errorf("TestRoundtrip %d: results\nexpected len(%d): %s\nreceived len(%d): %s\n", i, len(tc), tc, len(string(b)), string(b))
		}
	}
}
