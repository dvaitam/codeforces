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

const modE = 1000000007

func modPowE(a, e int64) int64 {
	res := int64(1)
	a %= modE
	for e > 0 {
		if e&1 == 1 {
			res = res * a % modE
		}
		a = a * a % modE
		e >>= 1
	}
	return res
}

func expectedAnswerE(n, k int64) int64 {
	if n == 0 {
		if k == 0 {
			return 1
		}
		return 0
	}
	if k > 4 {
		return 0
	}
	perm4 := int64(1)
	for i := int64(0); i < k; i++ {
		perm4 = perm4 * (4 - i) % modE
	}
	pw := modPowE(n+1, k)
	return perm4 * pw % modE
}

func generateCaseE(rng *rand.Rand) (int64, int64) {
	n := int64(rng.Intn(10))
	k := int64(rng.Intn(6))
	return n, k
}

func runCaseE(bin string, n, k int64) error {
	input := fmt.Sprintf("%d %d\n", n, k)
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
		return fmt.Errorf("bad output: %v", err)
	}
	expected := expectedAnswerE(n, k)
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k := generateCaseE(rng)
		if err := runCaseE(bin, n, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%d %d\n", i+1, err, n, k)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
