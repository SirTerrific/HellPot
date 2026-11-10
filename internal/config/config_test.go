package config

import (
	"reflect"
	"testing"
)

// The Docker deployment behind a reverse proxy relies on this variable; if the
// env key mapping breaks, the limit silently falls back to 10.
func TestMaxConnsPerIPFromEnv(t *testing.T) {
	t.Setenv("HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP", "100")
	associateExportedVariables()
	if MaxConnsPerIP != 100 {
		t.Fatalf("MaxConnsPerIP = %d, want 100", MaxConnsPerIP)
	}
}

// List options given as environment variables must not come out empty:
// an empty path list serves 404 everywhere.
func TestListsFromEnv(t *testing.T) {
	setDefaults()
	t.Setenv("HELLPOT_HTTP_ROUTER_PATHS", "a,b")
	t.Setenv("HELLPOT_HTTP_UAGENT__STRING__BLACKLIST", "x,y")
	associateExportedVariables()
	if want := []string{"a", "b"}; !reflect.DeepEqual(Paths, want) {
		t.Fatalf("Paths = %q, want %q", Paths, want)
	}
	if want := []string{"x", "y"}; !reflect.DeepEqual(UseragentBlacklistMatchers, want) {
		t.Fatalf("UseragentBlacklistMatchers = %q, want %q", UseragentBlacklistMatchers, want)
	}
}
