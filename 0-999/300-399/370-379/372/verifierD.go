package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type edge struct{ u, v int }

func connected(mask int, n int, edges []edge) bool {
	var start int = -1
	for i := 0; i < n; i++ {
		if mask&(1<<i) != 0 {
			start = i
			break
		}
	}
	if start == -1 {
		return false
	}
	visited := 0
	queue := []int{start}
	visited |= 1 << start
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		for _, e := range edges {
			a, b := e.u, e.v
			if a == u && mask&(1<<b) != 0 && visited&(1<<b) == 0 {
				visited |= 1 << b
				queue = append(queue, b)
			} else if b == u && mask&(1<<a) != 0 && visited&(1<<a) == 0 {
				visited |= 1 << a
				queue = append(queue, a)
			}
		}
	}
	return visited == mask
}

func consecutiveLen(mask int, n int) int {
	best := 0
	for l := 0; l < n; l++ {
		for r := l; r < n; r++ {
			if r-l+1 <= best {
				continue
			}
			ok := true
			for x := l; x <= r; x++ {
				if mask&(1<<x) == 0 {
					ok = false
					break
				}
			}
			if ok {
				best = r - l + 1
			}
		}
	}
	return best
}

func brute(n, k int, edges []edge) int {
	best := 0
	total := 1 << n
	for mask := 1; mask < total; mask++ {
		if bits.OnesCount(uint(mask)) > k {
			continue
		}
		if !connected(mask, n, edges) {
			continue
		}
		clen := consecutiveLen(mask, n)
		if clen > best {
			best = clen
		}
	}
	return best
}

func runCase(bin string, n, k int, edges []edge) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u+1, e.v+1))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := brute(n, k, edges)
	var got int
	fmt.Sscan(strings.TrimSpace(out.String()), &got)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	for i := 0; i < tests; i++ {
		n := rng.Intn(6) + 2
		k := rng.Intn(n) + 1
		edges := make([]edge, n-1)
		for v := 1; v < n; v++ {
			p := rng.Intn(v)
			edges[v-1] = edge{u: p, v: v}
		}
		if err := runCase(bin, n, k, edges); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
