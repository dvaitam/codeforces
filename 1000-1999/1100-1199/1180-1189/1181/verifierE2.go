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
	a, b, c, d int
}

func solve(rects []Rect) bool {
	ids0 := make([]int, len(rects))
	for i := range rects {
		ids0[i] = i
	}
	stack := [][]int{ids0}
	for len(stack) > 0 {
		ids := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		m := len(ids)
		if m <= 1 {
			continue
		}
		idsA := make([]int, m)
		copy(idsA, ids)
		sort.Slice(idsA, func(i, j int) bool { return rects[idsA[i]].a < rects[idsA[j]].a })
		maxC := rects[idsA[0]].c
		split := -1
		for i := 0; i < m-1; i++ {
			if rects[idsA[i]].c > maxC {
				maxC = rects[idsA[i]].c
			}
			if maxC <= rects[idsA[i+1]].a {
				split = i
				break
			}
		}
		if split >= 0 {
			left := make([]int, split+1)
			right := make([]int, m-split-1)
			copy(left, idsA[:split+1])
			copy(right, idsA[split+1:])
			stack = append(stack, left, right)
			continue
		}
		idsB := make([]int, m)
		copy(idsB, ids)
		sort.Slice(idsB, func(i, j int) bool { return rects[idsB[i]].b < rects[idsB[j]].b })
		maxD := rects[idsB[0]].d
		split = -1
		for i := 0; i < m-1; i++ {
			if rects[idsB[i]].d > maxD {
				maxD = rects[idsB[i]].d
			}
			if maxD <= rects[idsB[i+1]].b {
				split = i
				break
			}
		}
		if split >= 0 {
			left := make([]int, split+1)
			right := make([]int, m-split-1)
			copy(left, idsB[:split+1])
			copy(right, idsB[split+1:])
			stack = append(stack, left, right)
			continue
		}
		return false
	}
	return true
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
		fmt.Println("usage: go run verifierE2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 1
		rects := make([]Rect, n)
		for j := 0; j < n; j++ {
			x1 := rng.Intn(10)
			y1 := rng.Intn(10)
			w := rng.Intn(3) + 1
			h := rng.Intn(3) + 1
			rects[j] = Rect{x1, y1, x1 + w, y1 + h}
		}
		if err := runCase(bin, rects); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
