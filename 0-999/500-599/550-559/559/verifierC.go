package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type Point struct{ x, y int }

type testCase struct {
	h, w  int
	cells []Point
}

const mod int64 = 1000000007

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func comb(n, k int, fact, invFact []int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	return fact[n] * invFact[k] % mod * invFact[n-k] % mod
}

func solveCase(h, w int, cells []Point) int64 {
	n := len(cells)
	pts := append([]Point(nil), cells...)
	pts = append(pts, Point{h, w})
	sort.Slice(pts, func(i, j int) bool {
		if pts[i].x == pts[j].x {
			return pts[i].y < pts[j].y
		}
		return pts[i].x < pts[j].x
	})
	maxN := h + w
	fact := make([]int64, maxN+1)
	invFact := make([]int64, maxN+1)
	fact[0] = 1
	for i := 1; i <= maxN; i++ {
		fact[i] = fact[i-1] * int64(i) % mod
	}
	invFact[maxN] = modPow(fact[maxN], mod-2)
	for i := maxN; i > 0; i-- {
		invFact[i-1] = invFact[i] * int64(i) % mod
	}
	dp := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		p := pts[i]
		dp[i] = comb(p.x+p.y-2, p.x-1, fact, invFact)
		for j := 0; j < i; j++ {
			q := pts[j]
			if q.x <= p.x && q.y <= p.y {
				ways := dp[j] * comb(p.x-q.x+p.y-q.y, p.x-q.x, fact, invFact) % mod
				dp[i] = (dp[i] - ways + mod) % mod
			}
		}
	}
	return dp[n]
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", tc.h, tc.w, len(tc.cells))
	for _, c := range tc.cells {
		fmt.Fprintf(&sb, "%d %d\n", c.x, c.y)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := solveCase(tc.h, tc.w, tc.cells)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	cases := make([]testCase, 100)
	for i := range cases {
		h := rand.Intn(10) + 2
		w := rand.Intn(10) + 2
		n := rand.Intn(5)
		mp := make(map[Point]bool)
		for len(mp) < n {
			p := Point{rand.Intn(h) + 1, rand.Intn(w) + 1}
			if p.x == h && p.y == w {
				continue
			}
			if p.x == 1 && p.y == 1 {
				continue
			}
			mp[p] = true
		}
		cells := make([]Point, 0, len(mp))
		for p := range mp {
			cells = append(cells, p)
		}
		cases[i] = testCase{h, w, cells}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
