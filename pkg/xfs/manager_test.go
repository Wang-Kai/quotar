package xfs

import (
	"fmt"
	"testing"
)

func TestReadMappingInfo(t *testing.T) {
	manager := NewPrjManager()
	if err := manager.readMappingInfo(); err != nil {
		t.Error(err)
	}

	for _, item := range manager.Items {
		fmt.Printf("%v\n", item)
	}
}
