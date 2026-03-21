package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var rngState uint64 = 1337

func nextRand() uint64 {
	rngState ^= rngState << 13
	rngState ^= rngState >> 7
	rngState ^= rngState << 17
	return rngState
}

type Hash struct {
	h1, h2 uint64
}

func (a Hash) Xor(b Hash) Hash {
	return Hash{a.h1 ^ b.h1, a.h2 ^ b.h2}
}

type Element struct {
	h  Hash
	id int
}

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	minPrime := make([]int, n+1)
	primes := make([]int, 0, n/10)
	for i := 2; i <= n; i++ {
		if minPrime[i] == 0 {
			minPrime[i] = i
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p*i > n {
				break
			}
			minPrime[p*i] = p
			if i%p == 0 {
				break
			}
		}
	}

	primeHash := make([]Hash, n+1)
	for _, p := range primes {
		primeHash[p] = Hash{nextRand(), nextRand()}
	}

	valHash := make([]Hash, n+1)
	for i := 2; i <= n; i++ {
		p := minPrime[i]
		valHash[i] = valHash[i/p].Xor(primeHash[p])
	}

	factHash := make([]Hash, n+1)
	var P Hash
	for i := 1; i <= n; i++ {
		factHash[i] = factHash[i-1].Xor(valHash[i])
		P = P.Xor(factHash[i])
	}

	if P.h1 == 0 && P.h2 == 0 {
		printRes(n, nil)
		return
	}

	elements := make([]Element, n)
	for i := 1; i <= n; i++ {
		elements[i-1] = Element{factHash[i], i}
	}

	sort.Slice(elements, func(i, j int) bool {
		if elements[i].h.h1 != elements[j].h.h1 {
			return elements[i].h.h1 < elements[j].h.h1
		}
		return elements[i].h.h2 < elements[j].h.h2
	})

	find := func(h Hash) int {
		l, r := 0, n-1
		for l <= r {
			m := (l + r) / 2
			if elements[m].h == h {
				return elements[m].id
			}
			if elements[m].h.h1 < h.h1 || (elements[m].h.h1 == h.h1 && elements[m].h.h2 < h.h2) {
				l = m + 1
			} else {
				r = m - 1
			}
		}
		return -1
	}

	id := find(P)
	if id != -1 {
		printRes(n, []int{id})
		return
	}

	for i := 1; i <= n; i++ {
		target := P.Xor(factHash[i])
		id2 := find(target)
		if id2 != -1 && id2 != i {
			printRes(n, []int{i, id2})
			return
		}
	}

	k := n / 2
	printRes(n, []int{2, k, n})
}

func writeInt(out *bufio.Writer, n int) {
	if n == 0 {
		out.WriteByte('0')
		out.WriteByte(' ')
		return
	}
	var buf [20]byte
	i := 19
	for n > 0 {
		buf[i] = byte(n%10 + '0')
		i--
		n /= 10
	}
	out.Write(buf[i+1:])
	out.WriteByte(' ')
}

func printRes(n int, removed []int) {
	remMap := make(map[int]bool)
	for _, x := range removed {
		remMap[x] = true
	}
	fmt.Println(n - len(removed))
	out := bufio.NewWriter(os.Stdout)
	for i := 1; i <= n; i++ {
		if !remMap[i] {
			writeInt(out, i)
		}
	}
	out.WriteByte('\n')
	out.Flush()
}
