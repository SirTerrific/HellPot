package config

import "testing"

// The Docker deployment behind a reverse proxy relies on this variable; if the
// env key mapping breaks, the limit silently falls back to 10.
func TestMaxConnsPerIPFromEnv(t *testing.T) {
	t.Setenv("HELLPOT_PERFORMANCE_MAX__CONNS__PER__IP", "100")
	associateExportedVariables()
	if MaxConnsPerIP != 100 {
		t.Fatalf("MaxConnsPerIP = %d, want 100", MaxConnsPerIP)
	}
}
