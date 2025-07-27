package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &nums[i])
	}

	const limit = 1000
	primes := make([]int, 0)
	isPrime := make([]bool, limit+1)
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	for i := 2; i <= limit; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}

	type pair struct{ a, b int }
	pairCnt := make(map[pair]int)
	primeID := map[int]int{}
	adj := [][]int{}
	getID := func(x int) int {
		if id, ok := primeID[x]; ok {
			return id
		}
		id := len(adj)
		primeID[x] = id
		adj = append(adj, []int{})
		return id
	}

	for _, val := range nums {
		x := val
		factors := make([]int, 0, 2)
		tmp := x
		for _, p := range primes {
			if p*p > tmp {
				break
			}
			cnt := 0
			for tmp%p == 0 {
				tmp /= p
				cnt++
			}
			if cnt%2 == 1 {
				factors = append(factors, p)
			}
		}
		if tmp > 1 {
			factors = append(factors, tmp)
		}
		if len(factors) == 0 {
			fmt.Fprintln(writer, 1)
			return
		}
		if len(factors) == 1 {
			factors = append(factors, 1)
		}
		if factors[0] > factors[1] {
			factors[0], factors[1] = factors[1], factors[0]
		}
		pr := pair{factors[0], factors[1]}
		pairCnt[pr]++
		if pairCnt[pr] >= 2 {
			fmt.Fprintln(writer, 2)
			return
		}
		u := getID(factors[0])
		v := getID(factors[1])
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	startIDs := make([]int, 0)
	for p, id := range primeID {
		if p == 1 || p <= limit {
			startIDs = append(startIDs, id)
		}
	}

	const INF = int(1e9)
	minCycle := INF
	dist := make([]int, len(adj))
	parent := make([]int, len(adj))

	bfs := func(s int) int {
		for i := range dist {
			dist[i] = -1
			parent[i] = -1
		}
		q := make([]int, 1)
		q[0] = s
		dist[s] = 0
		head := 0
		best := INF
		for head < len(q) {
			u := q[head]
			head++
			if dist[u]*2 >= best {
				continue
			}
			for _, v := range adj[u] {
				if dist[v] == -1 {
					dist[v] = dist[u] + 1
					parent[v] = u
					q = append(q, v)
				} else if parent[u] != v {
					c := dist[u] + dist[v] + 1
					if c < best {
						best = c
					}
				}
			}
		}
		return best
	}

	for _, id := range startIDs {
		res := bfs(id)
		if res < minCycle {
			minCycle = res
		}
	}

	if minCycle == INF {
		fmt.Fprintln(writer, -1)
	} else {
		fmt.Fprintln(writer, minCycle)
	}
}
