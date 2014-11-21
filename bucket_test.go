package ccache

import (
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type BucketTests struct {
}

func Tests_Bucket(t *testing.T) {
	Expectify(new(BucketTests), t)
}

func (_ *BucketTests) GetMissFromBucket() {
	bucket := testBucket()
	Expect(bucket.get("invalid")).To.Equal(nil)
}

func (_ *BucketTests) GetHitFromBucket() {
	bucket := testBucket()
	item := bucket.get("power")
	assertValue(item, "9000")
}

func (_ *BucketTests) DeleteItemFromBucket() {
	bucket := testBucket()
	bucket.delete("power")
	Expect(bucket.get("power")).To.Equal(nil)
}

func (_ *BucketTests) SetsANewBucketItem() {
	bucket := testBucket()
	item, new, d := bucket.set("spice", TestValue("flow"), time.Minute)
	assertValue(item, "flow")
	item = bucket.get("spice")
	assertValue(item, "flow")
	Expect(new).To.Equal(true)
	Expect(d).To.Equal(1)
}

func (_ *BucketTests) SetsAnExistingItem() {
	bucket := testBucket()
	item, new, d := bucket.set("power", TestValue("9002"), time.Minute)
	assertValue(item, "9002")
	item = bucket.get("power")
	assertValue(item, "9002")
	Expect(new).To.Equal(false)
	Expect(d).To.Equal(0)
}

func (_ *BucketTests) ReplaceDoesNothingIfKeyDoesNotExist() {
	bucket := testBucket()
	Expect(bucket.replace("power", TestValue("9002"))).To.Equal(false)
	Expect(bucket.get("power")).To.Equal(nil)
}

func (_ *BucketTests) ReplaceReplacesThevalue() {
	bucket := testBucket()
	item, _, _ := bucket.set("power", TestValue("9002"), time.Minute)
	Expect(bucket.replace("power", TestValue("9004"))).To.Equal(true)
	Expect(item.Value().(string)).To.Equal("9004")
	Expect(bucket.get("power").Value().(string)).To.Equal("9004")
	//not sure how to test that the TTL hasn't changed sort of a sleep..
}

func testBucket() *bucket {
	b := &bucket{lookup: make(map[string]*Item)}
	b.lookup["power"] = &Item{
		key:   "power",
		value: TestValue("9000"),
	}
	return b
}

func assertValue(item *Item, expected string) {
	value := item.value.(TestValue)
	Expect(value).To.Equal(TestValue(expected))
}

type TestValue string

func (v TestValue) Expires() time.Time {
	return time.Now()
}
