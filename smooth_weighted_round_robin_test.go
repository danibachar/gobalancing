package gobalancing_test

import (
	"math"
	"math/rand"
	"time"

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
	// Global Vars
	var servers []server
	var smoothTestingServers []server
	var roundRobinServers []server
	var lb *SmoothWeightedRR

	// Helpers
	randFloat := func() float64 {
		minWeight := 0.0
		maxWeight := 20.0
		return minWeight + rand.Float64()*(maxWeight-minWeight)
	}
	initState := func() {
		lb = gobalancing.NewSWRR()
		smoothTestingServers = []server{
			{name: "server1", weight: 5},
			{name: "server2", weight: 1},
			{name: "server3", weight: 1},
		}
		rand.Seed(time.Now().UnixNano())
		servers = []server{
			{name: "server1", weight: randFloat()},
			{name: "server2", weight: randFloat()},
			{name: "server3", weight: randFloat()},
		}

		roundRobinServers = []server{
			{name: "server1", weight: math.SmallestNonzeroFloat64},
			{name: "server2", weight: math.SmallestNonzeroFloat64},
			{name: "server3", weight: math.SmallestNonzeroFloat64},
		}

		rand.Shuffle(len(servers), func(i, j int) { servers[i], servers[j] = servers[j], servers[i] })
	}

	addServer := func(index int) server {
		server := servers[index]
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

	getTotalServesWeight := func(servers []server) float64 {
		var weight float64
		for _, v := range servers {
			weight += v.weight
		}
		return weight
	}

	// Validations
	validateServerAdded := func(index int) {
		server := servers[index]
		serversMap := lb.All()
		Expect(serversMap[server.name]).To(Equal(server.weight))
	}

	validateAllServersAdded := func() {
		for i := 0; i < len(servers); i++ {
			validateServerAdded(i)
		}
	}

	validateEmptyLBState := func() {
		Expect(len(lb.All())).To(Equal((0)))
		for i := 0; i < 100; i++ {
			Expect(lb.Next()).To(BeNil())
		}
	}

	validateLoabBalancing := func(rounds int, servers []server) {
		totalServersWeight := getTotalServesWeight(servers)

		results := make(map[string]float64)

		for i := 0; i < rounds; i++ {
			s := lb.Next().(string)
			results[s]++
		}

		for i := 0; i < len(servers); i++ {
			s := servers[i]
			expectedWeight := float64(rounds) * s.weight / totalServersWeight
			// Expected weight should be +-1 from the weight counted
			Expect(results[s.name]).Should(Or(Equal(math.Ceil(expectedWeight)), Equal(math.Floor(expectedWeight))))
		}
	}

	BeforeEach(func() {
		initState()
	})

	When("lb is created", func() {
		It("its state is new", func() {
			validateEmptyLBState()
		})

		It("should not update, i.e fail", func() {
			server := servers[0]
			newWeight := 60.0
			err := lb.Update(server.name, newWeight)
			Expect(err).ToNot(BeNil())
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

	When("a updating an item", func() {
		It("should fail on sending nil", func() {
			addAllServers()
			err := lb.Update(nil, 100)
			Expect(err).ToNot(BeNil())
		})

		It("should fail on adding negative weight", func() {
			addAllServers()
			err := lb.Update(servers[0], -100)
			Expect(err).ToNot(BeNil())
		})
	})

	When("a single item is added", func() {
		It("should fail on adding nil interface", func() {
			err := lb.Add(nil, 100)
			Expect(err).ToNot(BeNil())
			validateEmptyLBState()
		})

		It("should fail on adding negative weight", func() {
			err := lb.Add(servers[0], -100)
			Expect(err).ToNot(BeNil())
			validateEmptyLBState()
		})

		It("should return contain it", func() {
			_ = addServer(0)
			validateServerAdded(0)
		})

		It("should return it all the time", func() {
			server := addServer(0)
			for i := 0; i < 10; i++ {
				Expect(lb.Next()).To(Equal(server.name))
			}
		})

		It("should not allow to add it again", func() {
			server := addServer(0)
			err := lb.Add(server.name, server.weight)
			Expect(err).ToNot(BeNil())
			Expect(len(lb.All())).To(Equal(1))
		})

		It("should allow to update it", func() {
			server := addServer(0)
			validateServerAdded(0)

			newWeight := 60.0
			err := lb.Update(server.name, newWeight)
			Expect(err).To(BeNil())

			serversMap := lb.All()
			Expect(serversMap[server.name]).To(Equal(newWeight))
		})
	})


	When("load balancing with all weights equal 0 - i.e Round Robin", func() {
		It("should balance equaly", func() {
			for _, s := range roundRobinServers {
				err := lb.Add(s.name, s.weight)
				Expect(err).To(BeNil())
			}

			for i := 0; i < 100 ;i++ {
				for _, s := range roundRobinServers {
					Expect(lb.Next().(string)).To(Equal(s.name))
				}
			}
		})
	})

	When("load balancing", func() {
		It("should balance between the servers by weight", func() {
			addAllServers()
			validateLoabBalancing(100, servers)
		})

		It("should smooth balance", func() {
			for _, s := range smoothTestingServers {
				err := lb.Add(s.name, s.weight)
				Expect(err).To(BeNil())
			}

			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[2].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
			Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))

			validateLoabBalancing(100, smoothTestingServers)
		})

		When("reseting", func() {
			It("should reset the state", func() {
				addAllServers()
				validateLoabBalancing(100, servers)
				lb.Reset()
				validateLoabBalancing(500, servers)
			})

			It("should reset smooth balance", func() {
				for _, s := range smoothTestingServers {
					err := lb.Add(s.name, s.weight)
					Expect(err).To(BeNil())
				}

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))

				lb.Reset()

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[2].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))

				validateLoabBalancing(100, smoothTestingServers)

			})
		})

		When("updating weight while balancing", func() {
			It("should accommodate the update", func() {
				for _, s := range smoothTestingServers {
					err := lb.Add(s.name, s.weight)
					Expect(err).To(BeNil())
				}

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))

				newServer := server{name: "server4", weight: 1}
				smoothTestingServers = append(smoothTestingServers, newServer)
				err := lb.Add(newServer.name, newServer.weight)
				Expect(err).To(BeNil())

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[2].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[3].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))

				validateLoabBalancing(100, smoothTestingServers)
			})
		})

		When("adding item while balancing", func() {
			It("should accommodate the addition", func() {
				for _, s := range smoothTestingServers {
					err := lb.Add(s.name, s.weight)
					Expect(err).To(BeNil())
				}

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))

				smoothTestingServers[1].weight = 3
				err := lb.Update(smoothTestingServers[1].name, smoothTestingServers[1].weight)
				Expect(err).To(BeNil())

				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[2].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[2].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[1].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))
				Expect(lb.Next().(string)).To(Equal(smoothTestingServers[0].name))

				validateLoabBalancing(100, smoothTestingServers)
			})
		})
	})
})
