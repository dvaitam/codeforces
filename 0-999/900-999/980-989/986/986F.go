package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"
)

// Item for priority queue used in Dijkstra
type Item struct {
	r int
	d int64
}

// Priority queue implementation
type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].d < pq[j].d }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

const inf int64 = 1 << 62

// mulMod computes (a*b) % mod using big integers
func mulMod(a, b, mod int64) int64 {
	return new(big.Int).Mod(new(big.Int).Mul(big.NewInt(a), big.NewInt(b)), big.NewInt(mod)).Int64()
}

// modPow computes (a^e) % mod
func modPow(a, e, mod int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = mulMod(res, a, mod)
		}
		a = mulMod(a, a, mod)
		e >>= 1
	}
	return res
}

// Miller-Rabin primality test for 64-bit integers
func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	small := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}
	for _, p := range small {
		if n%p == 0 {
			return n == p
		}
	}
	d := n - 1
	s := 0
	for d&1 == 0 {
		d >>= 1
		s++
	}
	bases := []int64{2, 325, 9375, 28178, 450775, 9780504, 1795265022}
	for _, a := range bases {
		if a%n == 0 {
			continue
		}
		x := modPow(a, d, n)
		if x == 1 || x == n-1 {
			continue
		}
		composite := true
		for r := 1; r < s; r++ {
			x = mulMod(x, x, n)
			if x == n-1 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

// Pollard's Rho factorization
func pollardsRho(n int64) int64 {
	if n%2 == 0 {
		return 2
	}
	for {
		x := rand.Int63n(n-2) + 2
		y := x
		c := rand.Int63n(n-1) + 1
		d := int64(1)
		for d == 1 {
			x = (mulMod(x, x, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			y = (mulMod(y, y, n) + c) % n
			d = gcd(abs(x-y), n)
			if d == n {
				break
			}
		}
		if d > 1 && d < n {
			return d
		}
	}
}

// factor fills map m with prime factors of n
func factor(n int64, m map[int64]int) {
	if n == 1 {
		return
	}
	if isPrime(n) {
		m[n]++
		return
	}
	d := pollardsRho(n)
	factor(d, m)
	factor(n/d, m)
}

type kInfo struct {
	single   bool
	p        int64
	minPrime int64
	dist     []int64
}

var cache = make(map[int64]*kInfo)

// precompute builds information for given k
func precompute(k int64) *kInfo {
	if v, ok := cache[k]; ok {
		return v
	}
	mp := make(map[int64]int)
	factor(k, mp)
	primes := make([]int64, 0, len(mp))
	for p := range mp {
		primes = append(primes, p)
	}
	sort.Slice(primes, func(i, j int) bool { return primes[i] < primes[j] })
	if len(primes) == 1 {
		info := &kInfo{single: true, p: primes[0]}
		cache[k] = info
		return info
	}
	pmin := primes[0]
	dist := make([]int64, pmin)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = 0
	pq := &PQ{{r: 0, d: 0}}
	heap.Init(pq)
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		if it.d != dist[it.r] {
			continue
		}
		for _, p := range primes {
			nr := int((int64(it.r) + p) % pmin)
			nd := it.d + p
			if nd < dist[nr] {
				dist[nr] = nd
				heap.Push(pq, Item{r: nr, d: nd})
			}
		}
	}
	info := &kInfo{single: false, minPrime: pmin, dist: dist}
	cache[k] = info
	return info
}

func main() {
	rand.Seed(time.Now().UnixNano())
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int64
		fmt.Fscan(in, &n, &k)
		if k == 1 {
			fmt.Fprintln(out, "NO")
			continue
		}
		info := precompute(k)
		if info.single {
			if n >= info.p && n%info.p == 0 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		} else {
			r := n % info.minPrime
			if n >= info.dist[r] {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}
