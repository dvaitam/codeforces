package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveB(a, b int64) float64 {
	if b == 0 {
		return 1.0
	} else if a <= 4*b {
		return 0.5 + float64(a)/(16.0*float64(b))
	}
	return 1.0 - float64(b)/float64(a)
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	t := 100
	cases := make([][2]int64, t)
	for i := 0; i < t; i++ {
		cases[i][0] = rand.Int63n(1_000_001)
		cases[i][1] = rand.Int63n(1_000_001)
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d %d\n", cases[i][0], cases[i][1])
	}
	expected := make([]float64, t)
	for i := 0; i < t; i++ {
		expected[i] = solveB(cases[i][0], cases[i][1])
	}
	outStr, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(strings.NewReader(outStr))
	scanner.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		valStr := scanner.Text()
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad float on test %d: %s\n", i+1, valStr)
			os.Exit(1)
		}
		if diff := val - expected[i]; diff > 1e-6 || diff < -1e-6 {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %.10f got %.10f\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
