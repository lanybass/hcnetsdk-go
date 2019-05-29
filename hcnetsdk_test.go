package hcnetsdk

import (
	"testing"
)

func TestHCNetSDK_Init(t *testing.T) {
	r := Init()
	if r != 1 {
		t.Errorf("init error: %v\n", r)
	}
}