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

type Rect struct {
	a, b, c, d int64
}

func solve(rs []Rect) bool {
	n := len(rs)
	if n <= 1 {
		return true
	}
	vs := make([]Rect, n)
	copy(vs, rs)
	sort.Slice(vs, func(i, j int) bool { return vs[i].a < vs[j].a })
	prefixMax := make([]int64, n)
	prefixMax[0] = vs[0].c
	for i := 1; i < n; i++ {
		if vs[i].c > prefixMax[i-1] {
			prefixMax[i] = vs[i].c
		} else {
			prefixMax[i] = prefixMax[i-1]
		}
	}
	suffixMin := make([]int64, n)
	suffixMin[n-1] = vs[n-1].a
	for i := n - 2; i >= 0; i-- {
		if vs[i].a < suffixMin[i+1] {
			suffixMin[i] = vs[i].a
		} else {
			suffixMin[i] = suffixMin[i+1]
		}
	}
	for i := 0; i < n-1; i++ {
		if prefixMax[i] <= suffixMin[i+1] {
			if solve(vs[:i+1]) && solve(vs[i+1:]) {
				return true
			}
		}
	}
	hs := make([]Rect, n)
	copy(hs, rs)
	sort.Slice(hs, func(i, j int) bool { return hs[i].b < hs[j].b })
	prefixMax = make([]int64, n)
	prefixMax[0] = hs[0].d
	for i := 1; i < n; i++ {
		if hs[i].d > prefixMax[i-1] {
			prefixMax[i] = hs[i].d
		} else {
			prefixMax[i] = prefixMax[i-1]
		}
	}
	suffixMin = make([]int64, n)
	suffixMin[n-1] = hs[n-1].b
	for i := n - 2; i >= 0; i-- {
		if hs[i].b < suffixMin[i+1] {
			suffixMin[i] = hs[i].b
		} else {
			suffixMin[i] = suffixMin[i+1]
		}
	}
	for i := 0; i < n-1; i++ {
		if prefixMax[i] <= suffixMin[i+1] {
			if solve(hs[:i+1]) && solve(hs[i+1:]) {
				return true
			}
		}
	}
	return false
}

func runCase(bin string, rects []Rect) error {
	var sb strings.Builder
	n := len(rects)
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, r := range rects {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r.a, r.b, r.c, r.d))
	}
	input := sb.String()
	expect := "NO"
	if solve(rects) {
		expect = "YES"
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.ToUpper(got) != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		rects := make([]Rect, n)
		for j := 0; j < n; j++ {
			x1 := int64(rng.Intn(10))
			y1 := int64(rng.Intn(10))
			w := int64(rng.Intn(3) + 1)
			h := int64(rng.Intn(3) + 1)
			rects[j] = Rect{x1, y1, x1 + w, y1 + h}
		}
		if err := runCase(bin, rects); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
