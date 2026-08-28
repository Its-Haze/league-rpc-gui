package types

import "testing"

func TestValidGameMode(t *testing.T) {
	for _, gm := range GameModes() {
		if !ValidGameMode(gm) {
			t.Errorf("ValidGameMode(%q) = false, want true", gm)
		}
	}
	if ValidGameMode("NOT_A_MODE") {
		t.Error("ValidGameMode(\"NOT_A_MODE\") = true, want false")
	}
	if ValidGameMode("") {
		t.Error("ValidGameMode(\"\") = true, want false")
	}
}
