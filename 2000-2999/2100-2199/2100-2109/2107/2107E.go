package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const smallLimit = 200000

type plan struct {
	ok     bool
	target int64
	chains []int
}

func choose3(x int64) int64 {
	return x * (x - 1) * (x - 2) / 6
}

func solveCase(n int, k int64) plan {
	maxWeight := choose3(int64(n))
	candidates := []int64{k - 1, k, k + 1}

	for _, target := range candidates {
		if target < 0 || target > maxWeight {
			continue
		}

		remNodes := n - 1 // nodes we can still place (root is node 1)
		remWeight := target
		chains := make([]int, 0)

		// Greedily carve out large chains to shrink the remaining weight.
		for remWeight > smallLimit && remNodes >= 2 {
			// Approximate the needed length by cube root of 6*remWeight.
			approx := int(math.Cbrt(float64(remWeight*6))) + 2
			if approx > remNodes+1 {
				approx = remNodes + 1
			}
			if approx < 3 {
				break
			}

			t := approx
			val := choose3(int64(t))
			for t > 3 && (int64(t-1) > int64(remNodes) || val > remWeight) {
				t--
				val = choose3(int64(t))
			}
			if t < 3 || int64(t-1) > int64(remNodes) || val > remWeight {
				break
			}

			chains = append(chains, t)
			remWeight -= val
			remNodes -= t - 1
		}

		// Quick impossibility check: the maximum achievable weight with remaining nodes
		// is by putting them all into one chain.
		if remWeight > choose3(int64(remNodes+1)) {
			continue
		}

		if remWeight == 0 {
			return plan{ok: true, target: target, chains: chains}
		}

		// Unbounded knapsack on the now-small remainder.
		maxLen := remNodes + 1
		if maxLen > 400 {
			maxLen = 400
		}

		// dp[w] = min nodes used to reach weight w.
		dp := make([]int, int(remWeight)+1)
		prev := make([][2]int, int(remWeight)+1) // previous weight, length used
		const inf = int(1e9)
		for i := range dp {
			dp[i] = inf
			prev[i] = [2]int{-1, -1}
		}
		dp[0] = 0

		for l := 3; l <= maxLen; l++ {
			val := choose3(int64(l))
			if val > remWeight {
				break
			}
			cost := l - 1
			valInt := int(val)
			for w := valInt; w <= int(remWeight); w++ {
				if dp[w-valInt]+cost < dp[w] {
					dp[w] = dp[w-valInt] + cost
					prev[w] = [2]int{w - valInt, l}
				}
			}
		}

		if dp[int(remWeight)] <= remNodes {
			// Reconstruct chains from dp.
			w := int(remWeight)
			for w > 0 {
				if prev[w][0] == -1 {
					break
				}
				chains = append(chains, prev[w][1])
				w = prev[w][0]
			}
			return plan{ok: true, target: target, chains: chains}
		}
	}

	return plan{ok: false}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}

	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)

		res := solveCase(n, k)
		if !res.ok {
			fmt.Fprintln(out, "No")
			continue
		}

		// Build edges.
		used := 1 // root already used
		edges := make([][2]int, 0, n-1)
		cur := 2
		for _, len := range res.chains {
			prev := 1
			for i := 0; i < len-1; i++ {
				edges = append(edges, [2]int{prev, cur})
				prev = cur
				cur++
				used++
			}
		}

		// Attach remaining nodes as leaves to the root.
		for used < n {
			edges = append(edges, [2]int{1, cur})
			cur++
			used++
		}

		fmt.Fprintln(out, "Yes")
		for _, e := range edges {
			fmt.Fprintf(out, "%d %d\n", e[0], e[1])
		}
	}
}
