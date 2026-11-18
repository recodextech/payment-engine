package services

import (
	"testing"

	"github.com/recodextech/api-definitions/events"
)

func row(code string, end int64, remaining int) events.SubscriptionEvent {
	var r events.SubscriptionEvent
	r.Payload.Plan.Code, r.Payload.PeriodEnd, r.Payload.RemainingJobs = code, end, remaining
	return r
}

func TestCurrentTakesFromSoonestEndingRow(t *testing.T) {
	cur := current{Plan: row("BASIC", 300, 7), Addons: []events.SubscriptionEvent{row("ADDON_A", 400, 5), row("ADDON_B", 200, 2)}}
	var order []string
	for _, r := range cur.rows() {
		order = append(order, r.Payload.Plan.Code)
	}
	if got := order[0] + "," + order[1] + "," + order[2]; got != "ADDON_B,BASIC,ADDON_A" {
		t.Fatalf("rows() order = %s", got)
	}
	if cur.remaining() != 14 {
		t.Fatalf("remaining() = %d, want 14", cur.remaining())
	}
}
