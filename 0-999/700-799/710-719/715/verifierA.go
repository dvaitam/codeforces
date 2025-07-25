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

func expected(n int) []string {
	xprev := big.NewInt(2)
	one := big.NewInt(1)
	res := make([]string, n)
	for k := 1; k <= n; k++ {
		sqrtFloor := new(big.Int).Sqrt(xprev)
		sqrtCeil := new(big.Int).Set(sqrtFloor)
		tmp := new(big.Int).Mul(sqrtFloor, sqrtFloor)
		if tmp.Cmp(xprev) < 0 {
			sqrtCeil.Add(sqrtCeil, one)
		}
		kp1 := big.NewInt(int64(k + 1))
		m0 := new(big.Int).Add(sqrtCeil, new(big.Int).Sub(kp1, one))
		m0.Div(m0, kp1)
		bk := big.NewInt(int64(k))
		m := new(big.Int).Add(m0, new(big.Int).Sub(bk, one))
		m.Div(m, bk)
		m.Mul(m, bk)
		t := new(big.Int).Mul(m, kp1)
		t2 := new(big.Int).Mul(t, t)
		diff := new(big.Int).Sub(t2, xprev)
		a := new(big.Int).Div(diff, bk)
		res[k-1] = a.String()
		xprev.Set(t)
	}
	return res
}

func runCase(bin string, n int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.Fields(strings.TrimSpace(out.String()))
	exp := expected(n)
	if len(got) != len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", len(exp), len(got))
	}
	for i := range exp {
		if got[i] != exp[i] {
			return fmt.Errorf("mismatch at line %d: expected %s got %s", i+1, exp[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
