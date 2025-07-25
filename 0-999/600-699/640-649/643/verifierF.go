package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n      uint64
	p      int
	q      int
	input  string
	expect uint32
}

func solve(n uint64, p, q int) uint32 {
	var m int
	if n == 0 {
		m = 0
	} else if p < int(n-1) {
		m = p
	} else {
		m = int(n - 1)
	}
	S := big.NewInt(0)
	cur := big.NewInt(1)
	nn := new(big.Int).SetUint64(n)
	for k := 1; k <= m; k++ {
		term := new(big.Int).Sub(nn, big.NewInt(int64(k-1)))
		cur.Mul(cur, term)
		cur.Div(cur, big.NewInt(int64(k)))
		S.Add(S, cur)
	}
	mask := new(big.Int).SetUint64(0xffffffff)
	Smod := new(big.Int).And(S, mask)
	mod := uint64(1) << 32
	var ans uint32
	for i := 1; i <= q; i++ {
		x := uint64(i)
		t := (x * x) & (mod - 1)
		t = (t * Smod.Uint64()) & (mod - 1)
		val := (x + t) & (mod - 1)
		ans ^= uint32(val)
	}
	return ans
}

func genCase(rng *rand.Rand) testCase {
	n := uint64(rng.Intn(100) + 1)
	p := rng.Intn(20) + 1
	q := rng.Intn(50) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, p, q)
	expect := solve(n, p, q)
	return testCase{n, p, q, sb.String(), expect}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		val, err := strconv.ParseUint(strings.TrimSpace(out), 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid integer output\n", i+1)
			os.Exit(1)
		}
		if uint32(val) != c.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, c.expect, val, c.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
