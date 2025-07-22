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

func charValue(r rune) int {
	if r >= '0' && r <= '9' {
		return int(r - '0')
	}
	return int(r-'A') + 10
}

func evaluate(digits []int, base int) int {
	v := 0
	for _, d := range digits {
		v = v*base + d
	}
	return v
}

func expectedRadixes(t string) []int {
	parts := strings.Split(t, ":")
	aStr, bStr := parts[0], parts[1]
	aDigits := make([]int, len(aStr))
	bDigits := make([]int, len(bStr))
	maxDigit := 0
	for i, r := range aStr {
		d := charValue(r)
		aDigits[i] = d
		if d > maxDigit {
			maxDigit = d
		}
	}
	for i, r := range bStr {
		d := charValue(r)
		bDigits[i] = d
		if d > maxDigit {
			maxDigit = d
		}
	}
	minBase := maxDigit + 1
	if minBase < 2 {
		minBase = 2
	}
	aConst := true
	for i := 0; i < len(aDigits)-1; i++ {
		if aDigits[i] != 0 {
			aConst = false
			break
		}
	}
	bConst := true
	for i := 0; i < len(bDigits)-1; i++ {
		if bDigits[i] != 0 {
			bConst = false
			break
		}
	}
	aVal := evaluate(aDigits, minBase)
	bVal := evaluate(bDigits, minBase)
	if aConst && bConst && aVal <= 23 && bVal <= 59 {
		return []int{-1}
	}
	res := []int{}
	for base := minBase; base <= 60; base++ {
		av := evaluate(aDigits, base)
		bv := evaluate(bDigits, base)
		if av <= 23 && bv <= 59 {
			res = append(res, base)
		}
	}
	if len(res) == 0 {
		return []int{0}
	}
	return res
}

func generateCase(rng *rand.Rand) string {
	l1 := rng.Intn(5) + 1
	l2 := rng.Intn(5) + 1
	sb := strings.Builder{}
	for i := 0; i < l1; i++ {
		v := rng.Intn(36)
		if v < 10 {
			sb.WriteByte(byte('0' + v))
		} else {
			sb.WriteByte(byte('A' + v - 10))
		}
	}
	sb.WriteByte(':')
	for i := 0; i < l2; i++ {
		v := rng.Intn(36)
		if v < 10 {
			sb.WriteByte(byte('0' + v))
		} else {
			sb.WriteByte(byte('A' + v - 10))
		}
	}
	return sb.String()
}

func runCase(bin string, input string, expect []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	got := []int{}
	for _, f := range fields {
		var v int
		fmt.Sscanf(f, "%d", &v)
		got = append(got, v)
	}
	if len(got) != len(expect) {
		return fmt.Errorf("expected %v got %v", expect, got)
	}
	for i := range got {
		if got[i] != expect[i] {
			return fmt.Errorf("expected %v got %v", expect, got)
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
	// deterministic edge cases
	for _, t := range []string{"0:0", "1:1", "A:1", "1:A"} {
		exp := expectedRadixes(t)
		if err := runCase(bin, t, exp); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %s failed: %v\n", t, err)
			os.Exit(1)
		}
	}
	for i := 0; i < 100; i++ {
		t := generateCase(rng)
		exp := expectedRadixes(t)
		if err := runCase(bin, t, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %s\n", i+1, err, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
