package main

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
	"text/tabwriter"
)

type Bucket struct {
	index int
	from  int
	to    int
}

func (bucket Bucket) String() string {
	return fmt.Sprintf("{from: %d, to: %d, index: %d}", bucket.from, bucket.to, bucket.index)
}

type Buckets struct {
	max int
	all []Bucket
}

func NewBuckets(weights []int) *Buckets {
	buckets := make([]Bucket, len(weights))
	for i := 0; i < len(weights); i++ {
		buckets[i] = Bucket{index: i, to: weights[i]}
	}
	slices.SortFunc(buckets, func(a, b Bucket) int {
		if a.to < b.to {
			return -1
		}
		if a.to == b.to {
			return 0
		}
		return 1
	})
	for i := 1; i < len(buckets); i++ {
		buckets[i].to = buckets[i].to + buckets[i-1].to
		buckets[i].from = buckets[i-1].to
	}
	max := buckets[len(buckets)-1].to
	return &Buckets{
		max: max,
		all: buckets,
	}
}

func bucketsCmp(bucket Bucket, value int) int {
	if bucket.from <= value && value < bucket.to {
		return 0
	}
	if value < bucket.from {
		return 1
	}
	return -1
}

func (buckets *Buckets) Samples(size int) []int {
	results := make([]int, size)
	for i := 0; i < size; i++ {
		bucketIndex, ok := slices.BinarySearchFunc(buckets.all, rand.Intn(buckets.max), bucketsCmp)
		if !ok {
			panic("bucket does not contains value")
		}
		results[i] = buckets.all[bucketIndex].index
	}
	return results
}

func main() {
	for _, weights := range [][]int{
		{1, 1},
		{1, 1, 1, 1, 1, 1},
		{10, 10, 10, 10, 10, 10},
		{1, 2},
		{1, 1, 8},
		{1, 2, 7},
	} {
		buckets := NewBuckets(weights)
		freq := make(map[int]int)
		checks := 10000
		size := 100
		for i := 0; i < checks; i++ {
			for _, value := range buckets.Samples(size) {
				freq[value] += 1
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.Debug)
		fmt.Println("===================================")
		fmt.Fprintf(w, "index\tweight\tavg[per %d]\n", size)
		for key, value := range freq {
			fmt.Fprintf(w, "%d\t%d\t%.2f\n", key, weights[key], float64(value)/float64(checks))
		}
		w.Flush()
	}
}
