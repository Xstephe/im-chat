package bitmap

import (
	"testing"
)

func TestBitmap_Set(t *testing.T) {
	b := NewBitmap(5)
	b.Set("pppp")
	b.Set("2222")
	b.Set("pppp")
	b.Set("ccc")

	for _, bit := range b.bits {
		t.Logf("%b,%v", bit, bit)
	}
}
