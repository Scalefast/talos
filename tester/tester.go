package tester

import (
	"testing"
)

type Event struct {
	id      int
	message string
}

// Standart tester
type StandardTester struct {
	*testing.T
}

func NewTester(t *testing.T) *StandardTester {
	return &StandardTester{t}
}

var (
	GivenMessage = Event{0, "Given: %s"}
	WhenMessage  = Event{1, "When:  %s"}
	ThenMessage  = Event{2, "Then:  %s"}
	AndMessage   = Event{3, "And:   %s"}
	ErrorMessage = Event{4, "Error: %s\n\t    Expected: %s\n\t    Actual:   %s"}
)

func (t StandardTester) Given(s string) {
	t.Logf(GivenMessage.message, s)
}
func (t StandardTester) When(s string) {
	t.Logf(WhenMessage.message, s)
}
func (t StandardTester) Then(s string) {
	t.Logf(ThenMessage.message, s)
}
func (t StandardTester) And(s string) {
	t.Logf(AndMessage.message, s)
}
func (t StandardTester) Error(msg, expected, reality string) {
	t.Fatalf(ErrorMessage.message, msg, expected, reality)
}
