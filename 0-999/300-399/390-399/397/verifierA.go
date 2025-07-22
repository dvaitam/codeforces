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

type interval struct {
	l int
	r int
}

func compute(n int, segs []interval) int {
	l1, r1 := segs[0].l, segs[0].r
	var inter []interval
	for i := 1; i < n; i++ {
		s := segs[i].l
		if s < l1 {
			s = l1
		}
		e := segs[i].r
		if e > r1 {
			e = r1
		}
		if s < e {
			inter = append(inter, interval{s, e})
		}
	}
	if len(inter) == 0 {
		return r1 - l1
	}
	sort.Slice(inter, func(i, j int) bool {
		if inter[i].l == inter[j].l {
			return inter[i].r < inter[j].r
		}
		return inter[i].l < inter[j].l
	})
	covered := 0
	curL, curR := inter[0].l, inter[0].r
	for _, iv := range inter[1:] {
		if iv.l <= curR {
			if iv.r > curR {
				curR = iv.r
			}
		} else {
			covered += curR - curL
			curL, curR = iv.l, iv.r
		}
	}
	covered += curR - curL
	res := (r1 - l1) - covered
	if res < 0 {
		res = 0
	}
	return res
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(100) + 1
	segs := make([]interval, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(100)
		r := l + 1 + rng.Intn(100-l)
		segs[i] = interval{l, r}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", segs[i].l, segs[i].r))
	}
	ans := compute(n, segs)
	return sb.String(), ans
}

func runCase(exe, input string, expected int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	var got int
	fmt.Sscan(outStr, &got)
	if got != expected {
		return fmt.Errorf("expected %d got %s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
