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

func weakness(arr []float64, x float64) float64 {
	maxEnd, minEnd := 0.0, 0.0
	maxVal, minVal := 0.0, 0.0
	for _, a := range arr {
		v := a - x
		maxEnd = math.Max(0, maxEnd+v)
		if maxEnd > maxVal {
			maxVal = maxEnd
		}
		minEnd = math.Min(0, minEnd+v)
		if -minEnd > minVal {
			minVal = -minEnd
		}
	}
	if maxVal > minVal {
		return maxVal
	}
	return minVal
}

func solveC(nums []int64) float64 {
	arr := make([]float64, len(nums))
	for i, v := range nums {
		arr[i] = float64(v)
	}
	l, r := -20000.0, 20000.0
	for iter := 0; iter < 60; iter++ {
		m1 := (l + r) / 2
		w1 := weakness(arr, m1)
		m2 := m1 + 1e-7
		w2 := weakness(arr, m2)
		if w1 < w2 {
			r = m2
		} else {
			l = m1
		}
	}
	return weakness(arr, (l+r)/2)
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		idx++
		n, _ := strconv.Atoi(line1)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing numbers line\n", idx)
			os.Exit(1)
		}
		numsStr := strings.Fields(strings.TrimSpace(scanner.Text()))
		if len(numsStr) != n {
			fmt.Fprintf(os.Stderr, "case %d expected %d numbers got %d\n", idx, n, len(numsStr))
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i, s := range numsStr {
			v, _ := strconv.ParseInt(s, 10, 64)
			arr[i] = v
		}
		expected := solveC(arr)
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.ParseFloat(got, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-float output %s\n", idx, got)
			os.Exit(1)
		}
		if diff := math.Abs(val - expected); diff > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
