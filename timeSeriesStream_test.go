package tsdb

import (
	"testing"
)

/*
func TestAppendAndRead(t *testing.T) {
	s := NewSeries()
	for _, p := range TestData {
		err := s.Append(p.Timestamp, p.Value, 1)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i, p := range TestData {
		timestamp, value, err := s.Read()
		if err != nil {
			t.Fatal(err)
		}
		if p.Timestamp != timestamp || p.Value != value {
			t.Fatalf("No.%d get (%v,%v),want (%v,%v)\n", i, timestamp, value, p.Timestamp, p.Value)
		}
	}

	s.Reset()
	if len(s.Bs.Stream) != 0 {
		t.Fatal("Reset failed!")
	}
}
*/

func TestTsAppendAndRead(t *testing.T) {
	s := NewSeries(nil)
	for _, p := range TestData {
		err := s.Append(p.Timestamp, p.Value, 1)
		if err != nil {
			t.Fatal(err)
		}
	}

	out, err := ReadValues(s.Bs.Stream, 1440583000, 1440591000, len(TestData))
	if err != nil {
		t.Fatal(err)
	}

	if len(out) != len(TestData) {
		t.Fatalf("length of data, want %d,get %d", len(TestData), len(out))
	}

	for i, p := range TestData {
		if p.Timestamp != out[i].Timestamp || p.Value != out[i].Value {
			t.Errorf("wrong result")
		}
	}

	// read at most 20 points
	out2, err := ReadValues(s.Bs.Stream, 1440583000, 1440591000, 20)
	if err != nil {
		t.Fatal(err)
	}

	if len(out2) != 20 {
		t.Fatal("wrong length of result")
	}

	for i := 0; i < 20; i++ {
		if TestData[i].Timestamp != out2[i].Timestamp || TestData[i].Value != out2[i].Value {
			t.Errorf("wrong result")
		}
	}

	// read points in a time range
	out3, err := ReadValues(s.Bs.Stream, 1440583260, 1440583741, len(TestData))
	if err != nil {
		t.Fatal(err)
	}

	if len(out3) != 9 {
		t.Fatal("wrong result")
	}
}

func benchmarkTimeSeriesStreamAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := NewSeries(nil)
		for j := 0; j < len(TestData); j++ {
			s.Append(TestData[j].Timestamp, TestData[j].Value, 1)
		}
	}
}

func benchmarkTimeSeriesStreamRead(b *testing.B) {
	s := NewSeries(nil)
	for j := 0; j < len(TestData); j++ {
		s.Append(TestData[j].Timestamp, TestData[j].Value, 1)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ReadValues(s.Bs.Stream, 1440583000, 1440591000, len(TestData))
	}
}
