package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func computeAnswer(n int, taps []int) int {
	ans := 0
	for i := 1; i <= n; i++ {
		minDist := n + 1
		for _, x := range taps {
			d := x - i
			if d < 0 {
				d = -d
			}
			if d < minDist {
				minDist = d
			}
		}
		tm := minDist + 1
		if tm > ans {
			ans = tm
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	const tests = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, tests)
	expected := make([]int, tests)
	for t := 0; t < tests; t++ {
		n := rand.Intn(200) + 1
		k := rand.Intn(n) + 1
		taps := make([]int, k)
		used := make(map[int]bool)
		for i := 0; i < k; i++ {
			for {
				v := rand.Intn(n) + 1
				if !used[v] {
					used[v] = true
					taps[i] = v
					break
				}
			}
		}
		// sort taps
		for i := 0; i < k; i++ {
			for j := i + 1; j < k; j++ {
				if taps[j] < taps[i] {
					taps[i], taps[j] = taps[j], taps[i]
				}
			}
		}
		fmt.Fprintf(&input, "%d %d\n", n, k)
		for i, v := range taps {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
		expected[t] = computeAnswer(n, taps)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to run binary:", err)
		os.Exit(1)
	}
	parts := strings.Fields(string(out))
	if len(parts) != tests {
		fmt.Printf("Expected %d outputs, got %d\n", tests, len(parts))
		os.Exit(1)
	}
	for i, p := range parts {
		got, err := strconv.Atoi(p)
		if err != nil || got != expected[i] {
			fmt.Printf("Test %d failed: expected %d got %s\n", i+1, expected[i], p)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
