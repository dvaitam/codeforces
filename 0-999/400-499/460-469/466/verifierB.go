package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(n, a, b int64) (int64, int64, int64) {
	m := 6 * n
	if a*b >= m {
		return a * b, a, b
	}
	swapped := false
	if a > b {
		a, b = b, a
		swapped = true
	}
	bestArea := int64(1<<62 - 1)
	var bestA, bestB int64
	lim := int64(math.Sqrt(float64(m))) + 2
	for A := a; A <= lim; A++ {
		B := (m + A - 1) / A
		if B < b {
			B = b
		}
		area := A * B
		if area < bestArea {
			bestArea = area
			bestA, bestB = A, B
		}
	}
	A := (m + b - 1) / b
	if A < a {
		A = a
	}
	area := A * b
	if area < bestArea {
		bestArea = area
		bestA, bestB = A, b
	}
	if swapped {
		bestA, bestB = bestB, bestA
	}
	return bestArea, bestA, bestB
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) != 3 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.ParseInt(parts[0], 10, 64)
		a, _ := strconv.ParseInt(parts[1], 10, 64)
		b, _ := strconv.ParseInt(parts[2], 10, 64)
		expS, expA, expB := solve(n, a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", n, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(string(out)))
		if len(outFields) < 3 {
			fmt.Printf("Test %d invalid output\n", idx)
			os.Exit(1)
		}
		gotS, _ := strconv.ParseInt(outFields[0], 10, 64)
		gotA, _ := strconv.ParseInt(outFields[1], 10, 64)
		gotB, _ := strconv.ParseInt(outFields[2], 10, 64)
		if gotS != expS || gotA != expA || gotB != expB {
			fmt.Printf("Test %d failed: expected %d %d %d got %s\n", idx, expS, expA, expB, strings.Join(outFields, " "))
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
