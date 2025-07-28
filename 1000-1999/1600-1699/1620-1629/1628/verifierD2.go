package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(n, m int, k int64) int64 {
	if m == 0 || m > n {
		return 0
	}
	mod := big.NewInt(1_000_000_007)
	inv := func(x *big.Int) *big.Int {
		exp := new(big.Int).Sub(mod, big.NewInt(2))
		return new(big.Int).Exp(x, exp, mod)
	}
	half := big.NewRat(1, 2)
	prev := make([]*big.Rat, m+2)
	curr := make([]*big.Rat, m+2)
	for i := range prev {
		prev[i] = new(big.Rat)
		curr[i] = new(big.Rat)
	}
	for i := 1; i <= n; i++ {
		upto := m
		if i < m {
			upto = i
		}
		for j := 1; j <= upto; j++ {
			if j == i {
				curr[j].SetInt64(int64(i))
				curr[j].Mul(curr[j], big.NewRat(k, 1))
				continue
			}
			delta := new(big.Rat).Sub(prev[j], prev[j-1])
			cmp0 := delta.Cmp(new(big.Rat))
			twoK := new(big.Rat).SetInt64(2 * k)
			if cmp0 <= 0 {
				curr[j].Set(prev[j])
			} else if delta.Cmp(twoK) >= 0 {
				curr[j].Set(prev[j-1])
				curr[j].Add(curr[j], big.NewRat(k, 1))
			} else {
				curr[j].Set(prev[j-1])
				delta.Mul(delta, half)
				curr[j].Add(curr[j], delta)
			}
		}
		if i <= m {
			curr[i+1] = new(big.Rat)
		}
		prev, curr = curr, prev
	}
	ans := prev[m]
	p := new(big.Int).Mod(ans.Num(), mod)
	q := new(big.Int).Mod(ans.Denom(), mod)
	qInv := inv(q)
	p.Mul(p, qInv)
	p.Mod(p, mod)
	return p.Int64()
}

func buildInput(n, m int, k int64) string {
	return fmt.Sprintf("1\n%d %d %d\n", n, m, k)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		m := rng.Intn(n) + 1
		k := int64(rng.Intn(5) + 1)
		input := buildInput(n, m, k)
		exp := fmt.Sprintf("%d", solve(n, m, k))
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("case %d wrong answer\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
