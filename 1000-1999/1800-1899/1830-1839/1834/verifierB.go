package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	L string
	R string
}

func expectedB(L, R string) string {
	// pad with leading zeros
	if len(L) < len(R) {
		L = strings.Repeat("0", len(R)-len(L)) + L
	} else if len(R) < len(L) {
		R = strings.Repeat("0", len(L)-len(R)) + R
	}
	n := len(L)
	ans := 0
	for i := 0; i < n; i++ {
		if L[i] != R[i] {
			diff := int(L[i]) - int(R[i])
			if diff < 0 {
				diff = -diff
			}
			ans = diff + 9*(n-i-1)
			break
		}
	}
	return fmt.Sprint(ans)
}

func genRandBig(maxDigits int) *big.Int {
	digits := rand.Intn(maxDigits) + 1
	var sb strings.Builder
	sb.WriteByte(byte('1' + rand.Intn(9)))
	for i := 1; i < digits; i++ {
		sb.WriteByte(byte('0' + rand.Intn(10)))
	}
	n := new(big.Int)
	n.SetString(sb.String(), 10)
	return n
}

func genTestsB() []testCaseB {
	rand.Seed(2)
	tests := make([]testCaseB, 0, 100)
	for len(tests) < 100 {
		a := genRandBig(18)
		b := new(big.Int).Add(a, big.NewInt(rand.Int63n(1000)))
		if rand.Intn(2) == 0 {
			a, b = b, a
		}
		if a.Cmp(b) > 0 {
			a, b = b, a
		}
		tests = append(tests, testCaseB{L: a.String(), R: b.String()})
	}
	return tests
}

func runCase(bin string, tc testCaseB) error {
	input := fmt.Sprintf("1\n%s %s\n", tc.L, tc.R)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := expectedB(tc.L, tc.R)
	if got != expect {
		return fmt.Errorf("L=%s R=%s expected %s got %s", tc.L, tc.R, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
