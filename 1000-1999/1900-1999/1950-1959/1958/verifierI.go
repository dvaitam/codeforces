package main

import (
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(n int, p1, p2 []int) int {
	anc1 := make([][]bool, n+1)
	anc2 := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		anc1[i] = make([]bool, n+1)
		anc2[i] = make([]bool, n+1)
	}
	for v := 1; v <= n; v++ {
		x := v
		for x != 0 {
			anc1[x][v] = true
			x = p1[x]
		}
		x = v
		for x != 0 {
			anc2[x][v] = true
			x = p2[x]
		}
	}
	m := n - 1
	adj := make([]uint64, m)
	for i := 2; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			var o1, o2 int
			if anc1[i][j] {
				o1 = 1
			} else if anc1[j][i] {
				o1 = -1
			}
			if anc2[i][j] {
				o2 = 1
			} else if anc2[j][i] {
				o2 = -1
			}
			if o1 != o2 {
				idxi := i - 2
				idxj := j - 2
				adj[idxi] |= 1 << idxj
				adj[idxj] |= 1 << idxi
			}
		}
	}
	if m == 0 {
		return 0
	}
	m1 := m / 2
	m2 := m - m1
	mask2 := uint64(1<<m2) - 1
	adjFirst := make([]uint64, m1)
	cross := make([]uint64, m1)
	adjSecond := make([]uint64, m2)
	for i := 0; i < m; i++ {
		if i < m1 {
			adjFirst[i] = adj[i] & ((1 << m1) - 1)
			cross[i] = (adj[i] >> m1) & mask2
		} else {
			adjSecond[i-m1] = (adj[i] >> m1) & mask2
		}
	}
	dp := make([]int, 1<<m2)
	for mask := 1; mask < 1<<m2; mask++ {
		v := bits.TrailingZeros(uint(mask))
		mWithout := mask &^ (1 << v)
		mWith := mWithout &^ int(adjSecond[v])
		choose := 1 + dp[mWith]
		skip := dp[mWithout]
		if choose > skip {
			dp[mask] = choose
		} else {
			dp[mask] = skip
		}
	}
	best := 0
	for mask := 0; mask < 1<<m1; mask++ {
		independent := true
		unionCross := uint64(0)
		for i := 0; i < m1 && independent; i++ {
			if mask&(1<<i) != 0 {
				if adjFirst[i]&uint64(mask) != 0 {
					independent = false
					break
				}
				unionCross |= cross[i]
			}
		}
		if !independent {
			continue
		}
		allowed := int(mask2 &^ unionCross)
		candidate := bits.OnesCount(uint(mask)) + dp[allowed]
		if candidate > best {
			best = candidate
		}
	}
	return 2 * (n - (1 + best))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(50)
	n := rand.Intn(8) + 2
	p1 := make([]int, n+1)
	p2 := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p1[i] = rand.Intn(i-1) + 1
	}
	for i := 2; i <= n; i++ {
		p2[i] = rand.Intn(i-1) + 1
	}

	var input strings.Builder
	fmt.Fprintln(&input, n)
	for i := 2; i <= n; i++ {
		if i > 2 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, p1[i])
	}
	input.WriteByte('\n')
	for i := 2; i <= n; i++ {
		if i > 2 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, p2[i])
	}
	input.WriteByte('\n')

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}
	var got int
	fmt.Sscan(strings.TrimSpace(string(outBytes)), &got)
	want := solve(n, p1, p2)
	if got != want {
		fmt.Printf("expected %d got %d\n", want, got)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
