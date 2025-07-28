package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func reference(n, k int, x, h []int, queries [][2]int) []int {
	type tower struct{ x, h int }
	towers := make([]tower, k)
	for i := 0; i < k; i++ {
		towers[i] = tower{x[i], h[i]}
	}
	sort.Slice(towers, func(i, j int) bool { return towers[i].x < towers[j].x })
	for i := 0; i < k; i++ {
		x[i] = towers[i].x
		h[i] = towers[i].h
	}
	const INF = int(1e18)
	prefix := make([]int, n+2)
	maxR := -INF
	t := 0
	for i := 1; i <= n; i++ {
		for t < k && x[t] <= i {
			if v := x[t] + h[t]; v > maxR {
				maxR = v
			}
			t++
		}
		prefix[i] = maxR
	}
	suffix := make([]int, n+2)
	minL := INF
	t = k - 1
	for i := n; i >= 1; i-- {
		for t >= 0 && x[t] >= i {
			if v := x[t] - h[t]; v < minL {
				minL = v
			}
			t--
		}
		suffix[i] = minL
	}
	costPoint := func(p int) int {
		best := INF
		if prefix[p] != -INF {
			v := p - prefix[p]
			if v < 0 {
				v = 0
			}
			if v < best {
				best = v
			}
		}
		if suffix[p] != INF {
			v := suffix[p] - p
			if v < 0 {
				v = 0
			}
			if v < best {
				best = v
			}
		}
		return best
	}
	ans := make([]int, len(queries))
	for qi, q := range queries {
		l, r := q[0], q[1]
		if l > r {
			l, r = r, l
		}
		cpL := costPoint(l)
		cpR := costPoint(r)
		separate := cpL + cpR
		mid := (l + r) / 2
		candidates := []int{INF, INF, INF, INF}
		if prefix[l] != -INF {
			v := r - prefix[l]
			if v < 0 {
				v = 0
			}
			candidates[0] = v
		}
		if suffix[r] != INF {
			v := suffix[r] - l
			if v < 0 {
				v = 0
			}
			candidates[1] = v
		}
		if prefix[mid] != -INF {
			v := r - prefix[mid]
			if v < 0 {
				v = 0
			}
			candidates[2] = v
		}
		if mid+1 <= n && suffix[mid+1] != INF {
			v := suffix[mid+1] - l
			if v < 0 {
				v = 0
			}
			candidates[3] = v
		}
		bestOne := INF
		for _, v := range candidates {
			if v < bestOne {
				bestOne = v
			}
		}
		if separate < bestOne {
			bestOne = separate
		}
		if bestOne < separate {
			ans[qi] = bestOne
		} else {
			ans[qi] = separate
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rand.Seed(48)
	n := rand.Intn(20) + 2
	k := rand.Intn(n) + 1
	x := make([]int, k)
	h := make([]int, k)
	used := map[int]bool{}
	for i := 0; i < k; i++ {
		for {
			v := rand.Intn(n) + 1
			if !used[v] {
				used[v] = true
				x[i] = v
				break
			}
		}
		h[i] = rand.Intn(n)
	}
	qn := 100
	queries := make([][2]int, qn)
	for i := 0; i < qn; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n) + 1
		for l == r {
			r = rand.Intn(n) + 1
		}
		queries[i] = [2]int{l, r}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, k)
	for i := 0; i < k; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, x[i])
	}
	input.WriteByte('\n')
	for i := 0; i < k; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprint(&input, h[i])
	}
	input.WriteByte('\n')
	fmt.Fprintln(&input, qn)
	for i := 0; i < qn; i++ {
		fmt.Fprintf(&input, "%d %d\n", queries[i][0], queries[i][1])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input.String())
	outBytes, err := cmd.Output()
	if err != nil {
		fmt.Println("error running binary:", err)
		os.Exit(1)
	}

	outLines := strings.Fields(strings.TrimSpace(string(outBytes)))
	if len(outLines) != qn {
		fmt.Printf("expected %d lines, got %d\n", qn, len(outLines))
		os.Exit(1)
	}
	want := reference(n, k, append([]int(nil), x...), append([]int(nil), h...), queries)
	for i, s := range outLines {
		var got int
		fmt.Sscan(s, &got)
		if got != want[i] {
			fmt.Printf("mismatch on query %d expected %d got %d\n", i+1, want[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
