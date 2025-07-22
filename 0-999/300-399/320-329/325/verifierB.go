package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func expected(n *big.Int) []uint64 {
	two := big.NewInt(2)
	eight := big.NewInt(8)
	results := make(map[uint64]struct{})
	for k := 0; k <= 60; k++ {
		p := new(big.Int).Lsh(big.NewInt(1), uint(k))
		B := new(big.Int).Lsh(big.NewInt(1), uint(k+1))
		B.Sub(B, big.NewInt(3))
		D := new(big.Int).Mul(B, B)
		tmp := new(big.Int).Mul(eight, n)
		D.Add(D, tmp)
		sqrtD := new(big.Int).Sqrt(D)
		if new(big.Int).Mul(sqrtD, sqrtD).Cmp(D) != 0 {
			continue
		}
		r := new(big.Int).Sub(sqrtD, B)
		if r.Sign() <= 0 || r.Bit(0) == 0 {
			continue
		}
		r.Div(r, two)
		if r.Sign() <= 0 {
			continue
		}
		T := new(big.Int).Mul(r, p)
		if T.BitLen() > 64 {
			continue
		}
		t64 := T.Uint64()
		if t64 > 0 {
			results[t64] = struct{}{}
		}
	}
	out := make([]uint64, 0, len(results))
	for t := range results {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func generateCase(rng *rand.Rand) (string, []uint64) {
	// generate n up to 1e12
	n := rng.Int63n(1_000_000_000_000) + 1
	nBig := big.NewInt(n)
	out := expected(nBig)
	return fmt.Sprintf("%d\n", n), out
}

func runCase(bin string, input string, exp []uint64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(exp) == 0 {
		if len(fields) != 1 || fields[0] != "-1" {
			return fmt.Errorf("expected -1 got %v", fields)
		}
		return nil
	}
	if len(fields) != len(exp) {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(fields))
	}
	for i, f := range fields {
		var val uint64
		if _, err := fmt.Sscan(f, &val); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if val != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, fields)
		}
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
