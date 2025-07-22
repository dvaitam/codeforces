package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type building struct {
	v int
	c int
}

type testCaseE struct {
	n  int
	s  int64
	bs []building
}

func generateCase(rng *rand.Rand) (string, testCaseE) {
	n := rng.Intn(6) + 1
	s := int64(rng.Intn(100) + 1)
	bs := make([]building, n)
	for i := 0; i < n; i++ {
		bs[i] = building{v: rng.Intn(10) + 1, c: rng.Intn(10)}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", n, s)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d %d\n", bs[i].v, bs[i].c)
	}
	return b.String(), testCaseE{n: n, s: s, bs: bs}
}

func expected(tc testCaseE) int64 {
	best := make(map[int]int)
	for _, b := range tc.bs {
		if pc, ok := best[b.v]; !ok || b.c < pc {
			best[b.v] = b.c
		}
	}
	vs := make([]int, 0, len(best))
	for v := range best {
		vs = append(vs, v)
	}
	sort.Ints(vs)
	filtered := make([]building, len(vs))
	for i, v := range vs {
		filtered[i] = building{v: v, c: best[v]}
	}
	const INF = 1e300
	dp := make([]float64, len(filtered))
	type line struct{ m, b float64 }
	type hull struct {
		lines []line
		xs    []float64
	}
	addLine := func(h *hull, m, b float64) {
		for len(h.lines) > 0 {
			last := h.lines[len(h.lines)-1]
			if m == last.m {
				if b >= last.b {
					return
				}
				h.lines = h.lines[:len(h.lines)-1]
				h.xs = h.xs[:len(h.xs)-1]
				continue
			}
			x := (last.b - b) / (m - last.m)
			if len(h.xs) > 0 && x <= h.xs[len(h.xs)-1] {
				h.lines = h.lines[:len(h.lines)-1]
				h.xs = h.xs[:len(h.xs)-1]
				continue
			}
			h.xs = append(h.xs, x)
			break
		}
		if len(h.lines) == 0 {
			h.xs = append(h.xs, -1e300)
		}
		h.lines = append(h.lines, line{m, b})
	}
	query := func(h hull, x float64) float64 {
		l, r := 0, len(h.lines)-1
		for l < r {
			m := (l + r + 1) / 2
			if h.xs[m] <= x {
				l = m
			} else {
				r = m - 1
			}
		}
		ln := h.lines[l]
		return ln.m*x + ln.b
	}
	var h hull
	for i, b := range filtered {
		if b.c == 0 {
			dp[i] = 0
		} else if len(h.lines) == 0 {
			dp[i] = INF
		} else {
			dp[i] = query(h, float64(b.c))
		}
		if dp[i] < INF {
			addLine(&h, 1/float64(b.v), dp[i])
		}
	}
	res := 1e300
	for i, b := range filtered {
		if dp[i] >= INF {
			continue
		}
		t := dp[i] + float64(tc.s)/float64(b.v)
		if t < res {
			res = t
		}
	}
	ans := int64(res)
	if float64(ans) < res-1e-12 {
		ans++
	}
	return ans
}

func runCase(bin string, input string, tc testCaseE) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expectedAns := fmt.Sprint(expected(tc))
	if outStr != expectedAns {
		return fmt.Errorf("expected %s got %s", expectedAns, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
