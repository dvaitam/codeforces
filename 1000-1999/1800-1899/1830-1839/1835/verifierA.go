package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func pow10(n int) int {
	res := 1
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}

func numDigits(x int) int {
	if x == 0 {
		return 1
	}
	d := 0
	for x > 0 {
		d++
		x /= 10
	}
	return d
}

func kthEquality(A, B, C int, k int64) string {
	aStart := pow10(A - 1)
	aEnd := pow10(A) - 1
	bStart := pow10(B - 1)
	bEnd := pow10(B) - 1
	cStart := pow10(C - 1)
	cEnd := pow10(C) - 1
	var list []string
	for a := aStart; a <= aEnd; a++ {
		for b := bStart; b <= bEnd; b++ {
			c := a + b
			if c >= cStart && c <= cEnd && numDigits(c) == C {
				list = append(list, fmt.Sprintf("%d + %d = %d", a, b, c))
			}
		}
	}
	sort.Strings(list)
	if int64(len(list)) < k {
		return "-1"
	}
	return list[k-1]
}

func generateCase(rng *rand.Rand) (string, string) {
	A := rng.Intn(3) + 1
	B := rng.Intn(3) + 1
	C := rng.Intn(3) + 1
	var eqs []string
	aStart := pow10(A - 1)
	aEnd := pow10(A) - 1
	bStart := pow10(B - 1)
	bEnd := pow10(B) - 1
	cStart := pow10(C - 1)
	cEnd := pow10(C) - 1
	for a := aStart; a <= aEnd; a++ {
		for b := bStart; b <= bEnd; b++ {
			c := a + b
			if c >= cStart && c <= cEnd && numDigits(c) == C {
				eqs = append(eqs, fmt.Sprintf("%d + %d = %d", a, b, c))
			}
		}
	}
	var k int64
	if len(eqs) == 0 {
		k = rng.Int63n(3) + 1
	} else {
		k = rng.Int63n(int64(len(eqs))+3) + 1
	}
	input := fmt.Sprintf("1\n%d %d %d %d\n", A, B, C, k)
	expected := kthEquality(A, B, C, k)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
