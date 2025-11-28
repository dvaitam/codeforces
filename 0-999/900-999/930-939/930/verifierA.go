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

func expectedA(n int, parents []int) int {
	depth := make([]int, n+1)
	count := map[int]int{}
	count[0] = 1
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		depth[i] = depth[p] + 1
		count[depth[i]]++
	}
	ans := 0
	for _, c := range count {
		if c%2 == 1 {
			ans++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	numTests := 100 // Number of test cases to generate
	for idx := 1; idx <= numTests; idx++ {
		input, expect := generateCase(rng)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: input:\n%s\nexpected:%s\ngot:%s\n", idx, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", numTests)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(99999) + 2 // N between 2 and 100,000
	
	parents := make([]int, n-1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	
	// Generate parents for p2, p3, ..., pn
	for i := 2; i <= n; i++ {
		// pi is a random inflorescence from 1 to i-1
		parents[i-2] = rng.Intn(i-1) + 1 
	}

	// Write parents to string builder
	for i := 0; i < n-1; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", parents[i])
	}
	sb.WriteByte('\n')

	// Calculate expected output
	expectedOutput := expectedA(n, parents)
	
	return sb.String(), strconv.Itoa(expectedOutput)
}
