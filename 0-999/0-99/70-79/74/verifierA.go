package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func randScore(minVal, maxVal int) int {
	if rand.Intn(2) == 0 {
		return 0
	}
	return minVal + rand.Intn(maxVal-minVal+1)
}

func generateTest() (string, string) {
	n := rand.Intn(10) + 1 // 1..10
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", n)
	bestScore := -1 << 60
	bestHandle := ""
	for i := 0; i < n; i++ {
		handle := fmt.Sprintf("h%d_%d", i, rand.Intn(1000))
		plus := rand.Intn(51)
		minus := rand.Intn(51)
		a := randScore(150, 500)
		b := randScore(300, 1000)
		c := randScore(450, 1500)
		d := randScore(600, 2000)
		e := randScore(750, 2500)
		fmt.Fprintf(&buf, "%s %d %d %d %d %d %d %d\n", handle, plus, minus, a, b, c, d, e)
		score := a + b + c + d + e + 100*plus - 50*minus
		if score > bestScore {
			bestScore = score
			bestHandle = handle
		}
	}
	return buf.String(), bestHandle
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
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
