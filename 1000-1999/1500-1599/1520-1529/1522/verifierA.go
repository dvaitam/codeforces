package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveA mirrors the reference solution 1522A.go.
// It reads the number of matches and prints 0 for each match.
func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	for i := 0; i < n; i++ {
		line, _ := in.ReadString('\n')
		_ = line
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("0\n")
	}
	return sb.String()
}

// generateMatchLine creates a line describing a match with random numbers.
func generateMatchLine(rng *rand.Rand) string {
	fields := rng.Intn(10) + 5
	var sb strings.Builder
	for i := 0; i < fields; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(100)))
	}
	return sb.String()
}

// generateCaseA builds a random test case for problem A.
func generateCaseA(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(generateMatchLine(rng))
		sb.WriteByte('\n')
	}
	return sb.String()
}

// runBinary executes the provided binary with the given input.
func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseA(rng)
	}
	for i, tc := range cases {
		expect := solveA(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
