package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type interval struct{ l, r int }

func solve(n, m int, intervals []interval, queries []interval) []int {
	res := make([]int, m)
	for qi, q := range queries {
		curr := q.l
		cnt := 0
		for curr < q.r {
			bestR := curr
			for _, in := range intervals {
				if in.l <= curr && in.r > bestR {
					bestR = in.r
				}
			}
			if bestR == curr {
				cnt = -1
				break
			}
			curr = bestR
			cnt++
		}
		res[qi] = cnt
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	ints := make([]interval, n)
	for i := range ints {
		l := rng.Intn(20)
		r := l + rng.Intn(10) + 1
		ints[i] = interval{l, r}
	}
	qs := make([]interval, m)
	for i := range qs {
		x := rng.Intn(20)
		y := x + rng.Intn(10) + 1
		qs[i] = interval{x, y}
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for _, in := range ints {
		fmt.Fprintf(&input, "%d %d\n", in.l, in.r)
	}
	for _, q := range qs {
		fmt.Fprintf(&input, "%d %d\n", q.l, q.r)
	}
	ans := solve(n, m, ints, qs)
	var output strings.Builder
	for i, v := range ans {
		if i > 0 {
			output.WriteByte(' ')
		}
		fmt.Fprintf(&output, "%d", v)
	}
	output.WriteByte('\n')
	return input.String(), output.String()
}

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
