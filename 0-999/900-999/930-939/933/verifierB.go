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

func convertNegBase(p, k int64) []int64 {
	if p == 0 {
		return []int64{0}
	}
	digits := make([]int64, 0)
	for p != 0 {
		r := ((p % int64(k)) + int64(k)) % int64(k)
		digits = append(digits, r)
		p = (p - r) / -int64(k)
	}
	return digits
}

func generateCaseB(rng *rand.Rand) (int64, int64) {
	p := rng.Int63n(1_000_000_000_000_000_000) + 1 // 1..1e18
	k := int64(rng.Intn(1999) + 2)
	return p, k
}

func runCaseB(bin string, p, k int64) error {
	input := fmt.Sprintf("%d %d\n", p, k)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outFields := strings.Fields(out.String())
	if len(outFields) == 1 && outFields[0] == "-1" {
		expected := convertNegBase(p, k)
		if len(expected) != 1 || expected[0] != -1 {
			return fmt.Errorf("unexpected -1 output")
		}
		return nil
	}
	if len(outFields) < 1 {
		return fmt.Errorf("no output")
	}
	var d int
	if _, err := fmt.Sscan(outFields[0], &d); err != nil {
		return fmt.Errorf("failed to read d: %v", err)
	}
	if d != len(outFields)-1 {
		return fmt.Errorf("length mismatch")
	}
	digits := make([]int64, d)
	for i := 0; i < d; i++ {
		if _, err := fmt.Sscan(outFields[i+1], &digits[i]); err != nil {
			return fmt.Errorf("failed to read digit %d", i)
		}
	}
	expected := convertNegBase(p, k)
	if len(expected) != len(digits) {
		return fmt.Errorf("expected length %d got %d", len(expected), len(digits))
	}
	for i := range expected {
		if expected[i] != digits[i] {
			return fmt.Errorf("digit %d expected %d got %d", i, expected[i], digits[i])
		}
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
	edge := [][2]int64{{0, 2}, {1, 2}, {12345, 2}, {12345, 3}, {999999999999999999, 2000}}
	for idx, e := range edge {
		if err := runCaseB(bin, e[0], e[1]); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		p, k := generateCaseB(rng)
		if err := runCaseB(bin, p, k); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
