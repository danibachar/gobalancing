package gobalancing_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/danibachar/gobalancing/smooth_weighted_round_robin"
)

type server struct {
	name string
	weight float64
}

servers = [...]server{
	server{name: "server1", weight: 10.0},
	server{name: "server2", weight: 1.0},
	server{name: "server3", weight: 4.5},
}

func testSWRR_Add(t *testing.T) { 

}

func TestSWRR_Next(t *testing.T) {
	lb := &SmoothWeightedRR{}

	lb.Add("server1", 5)
	lb.Add("server2", 2)
	lb.Add("server3", 3)

	results := make(map[string]int)

	for i := 0; i < 100; i++ {
		s := lb.Next().(string)
		results[s]++
	}

	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
		t.Error("the algorithm is wrong")
	}

	lb.Reset()
	results = make(map[string]int)

	for i := 0; i < 100; i++ {
		s := lb.Next().(string)
		results[s]++
	}

	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
		t.Error("the algorithm is wrong")
	}

	lb.RemoveAll()
	lb.Add("server1", 7)
	lb.Add("server2", 9)
	lb.Add("server3", 13)

	results = make(map[string]int)

	for i := 0; i < 29000; i++ {
		s := lb.Next().(string)
		results[s]++
	}

	if results["server1"] != 7000 || results["server2"] != 9000 || results["server3"] != 13000 {
		t.Error("the algorithm is wrong")
	}
}
