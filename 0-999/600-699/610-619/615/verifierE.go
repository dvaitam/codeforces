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

func spiralPos(n int64) (int64, int64) {
	if n == 0 {
		return 0, 0
	}
	low, high := int64(0), int64(1000000000)
	for low < high {
		mid := (low + high) / 2
		if 3*mid*(mid+1) >= n {
			high = mid
		} else {
			low = mid + 1
		}
	}
	r := low
	prev := 3 * (r - 1) * r
	k := n - prev
	x, y := r, int64(0)
	dirs := [6][2]int64{{-1, 1}, {-1, 0}, {0, -1}, {1, -1}, {1, 0}, {0, 1}}
	for i := 0; i < 6 && k > 0; i++ {
		step := r
		if k < step {
			step = k
		}
		x += dirs[i][0] * step
		y += dirs[i][1] * step
		k -= step
	}
	return x, y
}

func genCase(rng *rand.Rand) (string, [2]int64) {
	n := rng.Int63n(1_000_000_000)
	x, y := spiralPos(n)
	input := fmt.Sprintf("%d\n", n)
	return input, [2]int64{x, y}
}

func runCase(bin, input string, exp [2]int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var x, y int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &x, &y); err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if x != exp[0] || y != exp[1] {
		return fmt.Errorf("expected %d %d got %d %d", exp[0], exp[1], x, y)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases := []string{
		"0\n",
	}
	exps := [][2]int64{{0, 0}}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		in, exp := genCase(rng)
		cases = append(cases, in)
		exps = append(exps, exp)
	}

	for i := range cases {
		if err := runCase(bin, cases[i], exps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, cases[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
