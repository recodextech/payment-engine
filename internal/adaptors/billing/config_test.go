package billing

import "testing"

func TestReference(t *testing.T) {
	c := &Config{ReferencePrefix: "AA"}
	if got := c.Reference(1111111); got != "AA1111111" {
		t.Fatalf("Reference(1111111) = %s", got)
	}
}

func TestParseReference(t *testing.T) {
	c := &Config{ReferencePrefix: "AA"}
	cases := map[string]int64{"AA1111111": 1111111, "aa1111111": 1111111, "AA-111 1111": 1111111, "AA1": 1}
	for in, want := range cases {
		got, ok := c.ParseReference(in)
		if !ok || got != want {
			t.Errorf("ParseReference(%q) = %d, %v; want %d", in, got, ok, want)
		}
	}
	for _, in := range []string{"", "AA", "AA000", "FF1111111", "AA12B", "1111111", "AA99999999999999999999"} {
		if got, ok := c.ParseReference(in); ok {
			t.Errorf("ParseReference(%q) = %d; want invalid", in, got)
		}
	}
}
