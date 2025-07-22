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

func computeExpected(p, q, l, r int, a, b, c, d []int) int {
	count := 0
	for t := l; t <= r; t++ {
		ok := false
		for i := 0; i < p && !ok; i++ {
			for j := 0; j < q; j++ {
				if c[j]+t <= b[i] && d[j]+t >= a[i] {
					ok = true
					break
				}
			}
		}
		if ok {
			count++
		}
	}
	return count
}

func generateCase(rng *rand.Rand) (string, string) {
	p := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	l := rng.Intn(20)
	r := l + rng.Intn(20)

	a := make([]int, p)
	b := make([]int, p)
	for i := 0; i < p; i++ {
		a[i] = rng.Intn(30)
		b[i] = a[i] + rng.Intn(30)
	}
	c := make([]int, q)
	d := make([]int, q)
	for i := 0; i < q; i++ {
		c[i] = rng.Intn(30)
		d[i] = c[i] + rng.Intn(30)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", p, q, l, r)
	for i := 0; i < p; i++ {
		fmt.Fprintf(&sb, "%d %d\n", a[i], b[i])
	}
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d %d\n", c[i], d[i])
	}

	expCount := computeExpected(p, q, l, r, a, b, c, d)
	return sb.String(), fmt.Sprintf("%d", expCount)
}

func runCase(bin string, input string, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
