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

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func modPow(base, exp, mod int64) int64 {
	result := int64(1)
	for exp > 0 {
		if exp%2 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return result
}

func solve(n int) int64 {
	mersenneExponents := []int{
		2, 3, 5, 7, 13, 17, 19, 31, 61, 89,
		107, 127, 521, 607, 1279, 2203, 2281, 3217, 4253, 4423,
		9689, 9941, 11213, 19937, 21701, 23209, 44497, 86243, 110503, 132049,
		216091, 756839, 859433, 1257787, 1398269, 2976221, 3021377, 6972593, 13466917, 20996011,
	}
	if n > len(mersenneExponents) {
		// Fallback or error, though problem says n<=40
		return 0
	}
	q := mersenneExponents[n-1]
	mod := int64(1000000007)
	res := modPow(2, int64(q-1), mod)
	res = (res - 1 + mod) % mod
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Test all cases from 1 to 40 since they are few and deterministic
	for n := 1; n <= 40; n++ {
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test n=%d runtime error: %v\nInput:%s", n, err, input)
			return
		}
		expected := solve(n)
		got, errp := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if errp != nil || got != expected {
			fmt.Printf("Test n=%d FAILED\nInput:%sExpected:%d Got:%s\n", n, input, expected, out)
			return
		}
	}
	
	// Also run some random tests (redundant but keeping structure)
	for t := 1; t <= 20; t++ {
		n := r.Intn(40) + 1
		input := fmt.Sprintf("%d\n", n)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Random Test %d runtime error: %v\nInput:%s", t, err, input)
			return
		}
		expected := solve(n)
		got, errp := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
		if errp != nil || got != expected {
			fmt.Printf("Random Test %d FAILED\nInput:%sExpected:%d Got:%s\n", t, input, expected, out)
			return
		}
	}

	fmt.Println("All tests passed")
}