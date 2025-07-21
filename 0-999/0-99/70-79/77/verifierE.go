package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveE(R, r, K float64) float64 {
	delta := R - r
	x1 := R + r
	y1 := 2 * delta * K
	d1 := math.Sqrt(x1*x1 + y1*y1)
	ans := 2 * R * r / (d1 - delta)
	ans -= 2 * R * r / (d1 + delta)
	return ans
}

func run(bin, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	t := 100
	cases := make([][3]float64, t)
	for i := 0; i < t; i++ {
		R := float64(rand.Intn(10000-2) + 2)
		r := float64(rand.Intn(int(R-1)) + 1)
		K := float64(rand.Intn(10000) + 1)
		cases[i] = [3]float64{R, r, K}
	}
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d %d %d\n", int(cases[i][0]), int(cases[i][1]), int(cases[i][2]))
	}
	expected := make([]float64, t)
	for i := 0; i < t; i++ {
		expected[i] = solveE(cases[i][0], cases[i][1], cases[i][2])
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
			fmt.Fprintf(os.Stderr, "missing output on test %d\n", i+1)
			os.Exit(1)
		}
		valStr := scanner.Text()
		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad float on test %d: %s\n", i+1, valStr)
			os.Exit(1)
		}
		if diff := val - expected[i]; diff > 1e-6 || diff < -1e-6 {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %.6f got %.6f\n", i+1, expected[i], val)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", t)
}
