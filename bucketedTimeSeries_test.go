package tsdb

import (
	"testing"
)

var (
	bucketId     uint32 = 7
	timeSeriesId uint32 = 100
)

func TestBucketedTimeSeriesReset(t *testing.T) {
	//var b BucketedTimeSeries
	b := NewBucketedTimeSeries()
	b.Reset(5)
	if b.queriedBucketsAgo_ != 255 || len(b.blocks_) != 5 {
		t.Fatal("Reset() failed!")
	}
}

func TestBucketedTimeSeriesPutAndGet(t *testing.T) {
	inputData1 := TestData[:5]
	inputData2 := TestData[5:10]
	inputBlock1 := WriteValues(inputData1)
	inputBlock2 := WriteValues(inputData2)
	dataDirectory := DataDirectory_Test

	// var b BucketedTimeSeries
	b := NewBucketedTimeSeries()
	b.Reset(5)
	storage := NewBueketStorage(5, 1, dataDirectory)

	// input no.7 bucket
	for _, dp := range inputData1 {
		err := b.Put(7, timeSeriesId, dp, storage, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Input into a old bucket
	err := b.Put(6, timeSeriesId, TestData[0], storage, nil)
	if err == nil || err.Error() != "Invalid bucket number!" {
		t.Fatal("Wrong err message when input into a old bucket!")
	}

	// Test SetQueried()
	b.SetQueried()
	if b.queriedBucketsAgo_ != 0 {
		t.Fatal("SetQueried() failed!")
	}

	// input no.8 bucket
	var category uint16 = 2
	for _, dp := range inputData2 {
		err := b.Put(8, timeSeriesId, dp, storage, &category)
		if err != nil {
			t.Fatal(err)
		}
	}

	// get null bucket
	out1, err := b.Get(3, 4, storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(out1) != 0 {
		t.Fatal("Output data not null!")
	}

	// get one bucket
	out2, err := b.Get(8, 8, storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(out2) != 1 {
		t.Fatal("Length of output data is wrong!")
	}
	query := out2[0]
	if query.Count != inputBlock2.Count || IsEqualByteSlice(query.Data,
		inputBlock2.Data) == false {
		t.Fatal("input and output data mot match!")
	}

	// get two bucket
	out3, err := b.Get(0, 100, storage)
	if err != nil {
		t.Fatal(err)
	}
	if len(out3) != 2 {
		t.Fatal("Length of output data is wrong!")
	}
	query = out3[0]
	if query.Count != inputBlock1.Count || IsEqualByteSlice(query.Data,
		inputBlock1.Data) == false {
		t.Fatal("input and output data mot match!")
	}

	query = out3[1]
	if query.Count != inputBlock2.Count || IsEqualByteSlice(query.Data,
		inputBlock2.Data) == false {
		t.Fatal("input and output data mot match!")
	}
}

func TestBucketedTimeSeriesSetCurretBucket(t *testing.T) {
	dataDirectory := DataDirectory_Test
	inputData1 := TestData[:5]

	// var b BucketedTimeSeries
	b := NewBucketedTimeSeries()
	b.Reset(5)
	storage := NewBueketStorage(5, 1, dataDirectory)
	for _, dp := range inputData1 {
		err := b.Put(7, timeSeriesId, dp, storage, nil)
		if err != nil {
			t.Fatal(err)
		}
	}

	b.SetCurrentBucket(8, timeSeriesId, storage)
	if b.current_ != 8 {
		t.Fatal("b.current_ not match!")
	}

}

func TestBucketedTimeSeriesSetDataBlock(t *testing.T) {
	// var b BucketedTimeSeries
	b := NewBucketedTimeSeries()
	b.Reset(5)

	b.SetDataBlock(8, 5, 100)
	if b.current_ != 9 || b.blocks_[3] != 100 {
		t.Fatal("SetDataBlock() failed!")
	}
}

func benchmarkBucketedTimeSeriesPut(b *testing.B) {
	bucket := NewBucketedTimeSeries()
	bucket.Reset(5)
	dp := TestData[0]
	storage := NewBueketStorage(5, 1, dataDirectory)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bucket.Put(7, 100, dp, storage, nil)
	}
}
