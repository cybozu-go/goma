package average

import "testing"

func TestAverage(t *testing.T) {
	f, err := construct(nil)
	if err != nil {
		t.Fatal(err)
	}
	if f.Put(0) != 0 {
		t.Error("non-zero average")
	}
}
