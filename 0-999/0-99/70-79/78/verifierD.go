package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func isqrt(n int64) int64 {
	if n <= 0 {
		return 0
	}
	x := int64(math.Sqrt(float64(n)))
	for (x+1)*(x+1) <= n {
		x++
	}
	for x*x > n {
		x--
	}
	return x
}

func solve(k int64) int64 {
	R := k - 1
	T := 4 * R * R
	Ymax := isqrt(T / 3)
	var ans int64
	for y := int64(0); y <= Ymax; y++ {
		t := T - 3*y*y
		if t < 0 {
			continue
		}
		xmax := isqrt(t)
		var cnt int64
		if y%2 == 0 {
			cnt = 2*(xmax/2) + 1
		} else {
			cnt = 2 * ((xmax + 1) / 2)
		}
		if y == 0 {
			ans += cnt
		} else {
			ans += 2 * cnt
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, int64) {
	k := int64(rng.Intn(50) + 1)
	input := fmt.Sprintf("%d\n", k)
	return input, solve(k)
}

func runCase(bin, input string, expected int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
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
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
