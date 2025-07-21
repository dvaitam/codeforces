package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func generateTest() (string, string) {
	n := int64(rand.Intn(100) + 2)
	m := int64(rand.Intn(100) + 2)
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d\n", n, m)
	result := gcd(n-1, m-1) + 1
	return buf.String(), fmt.Sprintf("%d", result)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest()
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(inp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
