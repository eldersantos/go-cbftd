package cbftd

import "testing"

func checkInitialized(bh *ByteHistogram, t *testing.T) {
	for _, c := range bh.Count {
		if c > 0 {
			t.Errorf("expected Count value to be 0; got %v instead", c)
		}
	}
}

func TestNewByteHistogram(t *testing.T) {
	bh := NewByteHistogram()

	checkInitialized(bh, t)
}

func TestTrainByteHistogram(t *testing.T) {
	bh := NewByteHistogram()

	bh.Train("./testdata/")

	t.Log(bh.String())
}
