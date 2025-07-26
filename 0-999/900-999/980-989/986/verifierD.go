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

func maxProductForCost(S int64) *big.Int {
	var a, b int64
	switch S % 3 {
	case 0:
		b = S / 3
	case 1:
		if S < 4 {
			a = S / 2
		} else {
			b = (S - 4) / 3
			a = 2
		}
	case 2:
		if S < 2 {
			return big.NewInt(1)
		}
		b = (S - 2) / 3
		a = 1
	}
	res := new(big.Int).Exp(big.NewInt(3), big.NewInt(b), nil)
	if a > 0 {
		var t big.Int
		t.Exp(big.NewInt(2), big.NewInt(a), nil)
		res.Mul(res, &t)
	}
	return res
}

func solveCase(n *big.Int) int64 {
	if n.Cmp(big.NewInt(1)) == 0 {
		return 1
	}
	lo := int64(1)
	hi := int64(n.BitLen()*2 + 10)
	for lo < hi {
		mid := (lo + hi) / 2
		if maxProductForCost(mid).Cmp(n) >= 0 {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func genCase(rng *rand.Rand) (string, string) {
	n := big.NewInt(rng.Int63n(1000000) + 1)
	input := n.String() + "\n"
	expected := fmt.Sprintf("%d\n", solveCase(n))
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
