package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(n int) string {
	if n == 1 {
		return "0/1"
	}

	a := 0
	m := n
	for m%2 == 0 {
		a++
		m /= 2
	}

	if m == 1 {
		return fmt.Sprintf("%d/1", a)
	}

	var b big.Int
	var tmp big.Int

	r := 1
	l := 0
	for {
		b.Lsh(&b, 1)
		tmp.SetInt64(int64(r))
		b.Add(&b, &tmp)
		r = (r * 2) % m
		l++
		if r == 1 {
			break
		}
	}

	one := big.NewInt(1)

	var den big.Int
	den.SetInt64(1)
	den.Lsh(&den, uint(l))
	den.Sub(&den, one)

	var num big.Int
	num.Lsh(&b, 1)

	if a > 0 {
		var add big.Int
		add.SetInt64(int64(a))
		add.Mul(&add, &den)
		num.Add(&num, &add)
	}

	var g big.Int
	g.GCD(nil, nil, &num, &den)
	num.Quo(&num, &g)
	den.Quo(&den, &g)

	return fmt.Sprintf("%s/%s", num.String(), den.String())
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10000) + 1
	return fmt.Sprintf("%d\n", n)
}

func runCase(bin string, input, expected string) error {
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
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Also test small fixed cases
	for n := 1; n <= 20; n++ {
		input := fmt.Sprintf("%d\n", n)
		expected := solve(n)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case n=%d failed: %v\ninput:%s", n, err, input)
			os.Exit(1)
		}
	}

	w := bufio.NewWriter(os.Stderr)
	_ = w
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		var n int
		fmt.Sscan(strings.TrimSpace(input), &n)
		expected := solve(n)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
