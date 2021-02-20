package xfs

import "testing"

func TestCreatePrjDir(t *testing.T) {
	var dir = "steve"
	if err := createPrjDir(dir); err != nil {
		t.Fatal(err)
	}
}
