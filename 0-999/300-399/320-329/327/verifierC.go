package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const modC = 1000000007

func modPow(x, e int64) int64 {
	res := int64(1)
	x %= modC
	for e > 0 {
		if e&1 == 1 {
			res = (res * x) % modC
		}
		x = (x * x) % modC
		e >>= 1
	}
	return res
}

func modInv(x int64) int64 {
	return modPow(x, modC-2)
}

func expectedAnswer(a string, k int64) int64 {
	m := int64(len(a))
	var sumP int64
	var pow2 int64 = 1
	for i := int64(0); i < m; i++ {
		c := a[i]
		if c == '0' || c == '5' {
			sumP = (sumP + pow2) % modC
		}
		pow2 = (pow2 * 2) % modC
	}
	twoPowM := pow2
	var GS int64
	if twoPowM == 1 {
		GS = k % modC
	} else {
		numerator := (modPow(twoPowM, k) - 1 + modC) % modC
		denom := (twoPowM - 1 + modC) % modC
		GS = numerator * modInv(denom) % modC
	}
	return GS * sumP % modC
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(6) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('0' + rng.Intn(10))
	}
	k := int64(rng.Intn(20) + 1)
	return string(b), k
}

func runCase(bin string, a string, k int64) error {
	input := fmt.Sprintf("%s\n%d\n", a, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	expect := expectedAnswer(a, k)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		a, k := generateCase(rng)
		if err := runCase(bin, a, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
