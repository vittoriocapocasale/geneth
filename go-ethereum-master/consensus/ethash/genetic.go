package ethash

import (
	"math/bits"
	"math/rand"
	"sort"
	"time"
)

type Individual struct {
	Genotype []uint8
	Fitness  uint64
}

type GA struct {
	PopSize       int
	OffspringSize int
	MutRate       float32
	ParentHash    []byte
	MinerAddr     []byte
	Generations   uint64
	Population    []*Individual
	Best          *Individual
	Random        *rand.Rand
}

//Greate a new GA structure. It still needs initialization
func NewGA(popSize int, offspringSize int, mutRate float32, parentHash []byte, miner []byte) *GA {
	seed := rand.NewSource(time.Now().UnixNano())
	generator := rand.New(seed)
	return &GA{
		PopSize:       popSize,
		OffspringSize: offspringSize,
		MutRate:       mutRate,
		ParentHash:    parentHash,
		MinerAddr:     miner,
		Generations:   0,
		Population:    make([]*Individual, popSize),
		Best:          &Individual{[]byte{}, 0},
		Random:        generator}
}

/*
func main() {
	//nonce, _ := hex.DecodeString("d3ee432b4fb3d26b")
	prevBlockHash, _ := hex.DecodeString("44bca881b07a6a09f83b130798072441705d9a665c5ac8bdf2f39a3cdf3bee29")
	ga := NewGA(50, 300, 0.05, prevBlockHash)
	ga.Init()
	for ga.Best.Fitness < 225 && ga.Generations < 100000 {
		ga.Evolve()
	}
	fmt.Println(ga.Best.Fitness, ga.Generations)
}
*/

//Initialize a GA. Can be called multiple times in order to start from scratch
func (self *GA) Init() {
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	self.Random = rnd
	self.Generations = 0
	self.Best = &Individual{[]byte{}, 0}
	self.Population = []*Individual{}
	for i := 0; i < self.PopSize; i++ {
		individual := self.RandomIndividual()
		if individual.Fitness > self.Best.Fitness {
			self.Best = individual
		}
		self.Population = append(self.Population, individual)
	}
}

//performs a single evolution step
func (self *GA) Evolve() {
	self.Generations++
	nextGen := []*Individual{}
	for i := 0; i < self.OffspringSize; i++ {
		child := self.Reproduce()
		if child.Fitness > self.Best.Fitness {
			self.Best = child
			//fmt.Println("best", best)
		}
		nextGen = append(nextGen, child)
	}
	//nextGen = append(nextGen, self.Population...)
	sort.Slice(nextGen, func(i, j int) bool {
		return nextGen[i].Fitness > nextGen[j].Fitness
	})
	self.Population = nextGen[:self.PopSize]

}

//binary tournement selection
func (self *GA) TournamentSelection() *Individual {
	candidate1 := self.Population[self.Random.Intn(len(self.Population))]
	candidate2 := self.Population[self.Random.Intn(len(self.Population))]
	if candidate2.Fitness > candidate1.Fitness {
		candidate1 = candidate2
	}
	return candidate1
}

//generates a random individual. Useful in the initialization phase
func (self *GA) RandomIndividual() *Individual {
	genotype := make([]byte, 8)
	self.Random.Read(genotype)
	return &Individual{genotype, self.Fitness(genotype)}
}

//wrapper function
func (self *GA) Reproduce() *Individual {
	/*if self.MutRate > self.Random.Float32() {
		return self.Mutation()
	} else {
		return self.Crossover()
	}*/
	return self.Crossover2()
}

//child inherits bits that are common to both parents, the others are selected randomly. Genes may mutate 
func (self *GA) Crossover4() *Individual {
	masks := [8]uint8{0b00000001, 0b00000010, 0b00000100, 0b00001000, 0b00010000, 0b00100000, 0b01000000, 0b10000000}
	father := self.TournamentSelection()
	mother := self.TournamentSelection()
	offspring := []uint8{}
	for i := 0; i < len(father.Genotype); i++ {
		octet := father.Genotype[i] & mother.Genotype[i]
		for j := 0; j < len(masks); j++ {
			//se il figlio ha bit a zero ma uno tra padre e madre lo aveva a 1, metti il bit a 1 con probabilità 50%
			if (octet&masks[j] == 0) && ((father.Genotype[i]&masks[j]) > 0 || (mother.Genotype[i]&masks[j]) > 0) {
				if self.Random.Intn(2) == 1 {
					octet += masks[j]
				}
			} else {
				//mutation with probability mutRate
				if self.MutRate > self.Random.Float32() {
					if octet&masks[j] == 0 {
						octet += masks[j]
					} else {
						octet -= masks[j]
					}
				}
			}
		}
		offspring = append(offspring, octet)
	}
	return &Individual{offspring, self.Fitness(offspring)}
}

//one cut crossover
func (self *GA) Crossover3() *Individual {
	masks := [8]uint8{0b00000001, 0b00000010, 0b00000100, 0b00001000, 0b00010000, 0b00100000, 0b01000000, 0b10000000}
	father := self.TournamentSelection()
	mother := self.TournamentSelection()
	cut := self.Random.Intn(len(father.Genotype))
	genotype := append(father.Genotype[:cut], mother.Genotype[cut:]...)
	offspring := &Individual{genotype, self.Fitness(genotype)}
	for i := 0; i < len(offspring.Genotype); i++ {
		for j := 0; j < len(masks); j++ {
			//mutation with probability mutRate
			if self.MutRate > self.Random.Float32() {
				if offspring.Genotype[i]&masks[j] == 0 {
					offspring.Genotype[i] += masks[j]
				} else {
					offspring.Genotype[i] -= masks[j]
				}
			}
		}
	}
	return offspring
}

//hinerits blocks of 1 byte randomly from one parent.  Genes may mutate 
func (self *GA) Crossover2() *Individual {
	father := self.TournamentSelection()
	mother := self.TournamentSelection()
	masks := [8]uint8{0b00000001, 0b00000010, 0b00000100, 0b00001000, 0b00010000, 0b00100000, 0b01000000, 0b10000000}
	offspring := []uint8{}
	for i := 0; i < len(father.Genotype); i++ {
		if self.Random.Intn(2) == 1 {
			offspring = append(offspring, father.Genotype[i])
		} else {
			offspring = append(offspring, mother.Genotype[i])
		}
		for j := 0; j < len(masks); j++ {
			//mutation with probability mutRate
			if self.MutRate > self.Random.Float32() {
				if offspring[i]&masks[j] == 0 {
					offspring[i] += masks[j]
				} else {
					offspring[i] -= masks[j]
				}
			}
		}
	}
	return &Individual{offspring, self.Fitness(offspring)}
}

//hinerits blocks of 1 byte randomly from one parent.  Mutation should be handled externally 
func (self *GA) Crossover() *Individual {
	father := self.TournamentSelection()
	mother := self.TournamentSelection()
	offspring := []uint8{}
	for i := 0; i < len(father.Genotype); i++ {
		if self.Random.Intn(2) == 1 {
			offspring = append(offspring, father.Genotype[i])
		} else {
			offspring = append(offspring, mother.Genotype[i])
		}
	}
	return &Individual{offspring, self.Fitness(offspring)}
}

//mutates an individual. Each gene may mutate independently from the others
func (self *GA) Mutation() *Individual {
	masks := [8]uint8{0b00000001, 0b00000010, 0b00000100, 0b00001000, 0b00010000, 0b00100000, 0b01000000, 0b10000000}
	candidate := self.TournamentSelection()
	offspring := []uint8{}
	offspring = append(offspring, candidate.Genotype...)
	pos1 := self.Random.Intn(len(offspring))
	pos2 := self.Random.Intn(len(masks))
	if offspring[pos1]&masks[pos2] == 0 {
		offspring[pos1] += masks[pos2]
	} else {
		offspring[pos1] -= masks[pos2]
	}
	return &Individual{offspring, self.Fitness(offspring)}
}


//A simple fitness function. The objective is to maximize the number of 1 in the output of a function f.
//The output can only be manipulated through the nonce (genotype).
//f depends also on a challenge. The challenge is generated from previous block hash and the wallet address of the miner (this means that miners get sligthly different challenges)
//f (challenge, nonce) = xor (challenge, nonce+nonce+nonce+nonce), where the + indicates a concatenation.
//Since the challenge is longer than the nonce, a given nonce may be good for the "first" xor and bad for the three others.
//There is a push to both get an high number of 1 in at least one byte and never get a low number of 1 in any bytes (bonus)   
//The fitness could be anything, but the problem should be NP (not sure if this problem is NP)
func (self *GA) Fitness(genotype []uint8) uint64 {
	var total uint64 = 0
	challenge := []uint8{}
	var bonus uint64 = 0
	max := []uint64{0, 0, 0, 0, 0, 0, 0, 0}
	min := []uint64{8, 8, 8, 8, 8, 8, 8, 8}
	//xor of minerAddr and first part of parentHash
	for i := 0; i < len(self.MinerAddr); i++ {
		challenge = append(challenge, self.MinerAddr[i]^self.ParentHash[i])
	}
	//append remaining part of parent hash
	challenge = append(challenge, self.ParentHash[len(self.MinerAddr):]...)
	for i := 0; i < len(challenge); i++ {
		v := challenge[i] ^ genotype[i%8]
		n := uint64(bits.OnesCount8(v))
		if n > max[i%8] {
			max[i%8] = n
		}
		if n < min[i%8] {
			min[i%8] = n
		}
		total += n
	}
	for i := 0; i < len(max); i++ {
		bonus = bonus + 7*max[i] + 11*min[i]
	}
	return 2*total + bonus
}
