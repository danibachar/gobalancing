package gobalancing_test

import (
	"github.com/danibachar/gobalancing"
	. "github.com/danibachar/gobalancing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type server struct {
	name   string
	weight float64
}

var _ = Describe("Smooth Weighted RR Public API Test", func() {

	var servers = [...]server{
		{name: "server1", weight: 10.0},
		{name: "server2", weight: 1.0},
		{name: "server3", weight: 4.5},
	}
	var lb *SmoothWeightedRR

	validateServerAdded := func(index int) {
		var server = servers[index]
		var serversMap = lb.All()
		Expect(serversMap[server.name]).To(Equal(server.weight))
	}

	validateAllServersAdded := func() {
		for i := 0; i < len(servers); i++ {
			validateServerAdded(i)
		}
	}

	validateEmptyLBState := func() {
		Expect(len(lb.All())).To(Equal((0)))
	}

	addServer := func(index int) server {
		var server = servers[index]
		err := lb.Add(server.name, server.weight)
		if err != nil {
			Fail(err.Error())
		}
		return server
	}

	addAllServers := func() {
		for i := 0; i < len(servers); i++ {
			_ = addServer(i)
		}
	}

	BeforeEach(func() {
		lb = gobalancing.NewSWRR()
	})

	When("lb is created", func() {
		It("its state is new", func() {
			validateEmptyLBState()
		})
	})

	When("removing all items", func() {
		It("should return an empty map", func() {
			addAllServers()
			validateAllServersAdded()
			lb.RemoveAll()
			validateEmptyLBState()
		})
	})

	When("a single item is added", func() {
		BeforeEach(func() {
			lb = gobalancing.NewSWRR()
		})

		It("should return contain it", func() {
			_ = addServer(0)
			validateServerAdded(0)
		})

		It("should return it all the time", func() {
			var server = addServer(0)
			for i := 0; i < 10; i++ {
				Expect(lb.Next()).To(Equal(server.name))
			}
		})

		It("should not allow to add it again", func() {
			var server = addServer(0)
			var err = lb.Add(server.name, server.weight)
			Expect(err).ToNot(BeNil())
		})

		It("should allow to update it", func() {
			var server = addServer(0)
			validateServerAdded(0)

			var newWeight = 60.0
			var err = lb.Update(server.name, newWeight)

			var serversMap = lb.All()
			Expect(err).To(BeNil())
			Expect(serversMap[server.name]).To(Equal(newWeight))
		})
	})

	When("load balancing", func() {
		It("should balance between the servers by weight", func() {
			addAllServers()
			results := make(map[string]int)
			for i := 0; i < 100; i++ {
				s := lb.Next().(string)
				results[s]++
			}
			for i := 0; i < len(servers); i++ {

			}
		})
	})
})

// func testSWRR_Add(t *testing.T) {

// }

// func TestSWRR_Next(t *testing.T) {
// 	lb := &SmoothWeightedRR{}

// 	lb.Add("server1", 5)
// 	lb.Add("server2", 2)
// 	lb.Add("server3", 3)

// results := make(map[string]int)

// 	for i := 0; i < 100; i++ {
// 		s := lb.Next().(string)
// 		results[s]++
// 	}

// 	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
// 		t.Error("the algorithm is wrong")
// 	}

// 	lb.Reset()
// 	results = make(map[string]int)

// 	for i := 0; i < 100; i++ {
// 		s := lb.Next().(string)
// 		results[s]++
// 	}

// 	if results["server1"] != 50 || results["server2"] != 20 || results["server3"] != 30 {
// 		t.Error("the algorithm is wrong")
// 	}

// 	lb.RemoveAll()
// 	lb.Add("server1", 7)
// 	lb.Add("server2", 9)
// 	lb.Add("server3", 13)

// 	results = make(map[string]int)

// 	for i := 0; i < 29000; i++ {
// 		s := lb.Next().(string)
// 		results[s]++
// 	}

// 	if results["server1"] != 7000 || results["server2"] != 9000 || results["server3"] != 13000 {
// 		t.Error("the algorithm is wrong")
// 	}
// }
