package exec

import "testing"

func TestMergeEnv(t *testing.T) {
	t.Parallel()

	env1 := []string{"A=123", "B=456"}
	env2 := []string{"A=789", "C=11111"}
	merged := mergeEnv(env1, env2)

	if len(merged) != 3 {
		t.Fatal(`len(merged) != 3`)
	}
	if merged[0] != "A=123" {
		t.Error(`merged[0] != "A=123"`)
	}
	if merged[1] != "B=456" {
		t.Error(`merged[1] != "B=456"`)
	}
	if merged[2] != "C=11111" {
		t.Error(`merged[2] != "C=11111"`)
	}
}
