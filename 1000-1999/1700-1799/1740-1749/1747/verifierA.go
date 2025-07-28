package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

// generateTests returns the input and expected outputs for problem A
func generateTests() (string, []int64) {
	rand.Seed(1)
	t := 100
	var buf bytes.Buffer
	var answers []int64
	fmt.Fprintln(&buf, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1 // n in [1,5]
		fmt.Fprintln(&buf, n)
		sum := int64(0)
		for j := 0; j < n; j++ {
			x := rand.Intn(21) - 10 // values in [-10,10]
			fmt.Fprint(&buf, x)
			if j+1 < n {
				fmt.Fprint(&buf, " ")
			}
			sum += int64(x)
		}
		fmt.Fprintln(&buf)
		if sum < 0 {
			sum = -sum
		}
		answers = append(answers, sum)
	}
	return buf.String(), answers
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	input, expected := generateTests()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(scanner.Text(), &got)
		if got != exp {
			fmt.Printf("case %d: expected %d, got %d\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
