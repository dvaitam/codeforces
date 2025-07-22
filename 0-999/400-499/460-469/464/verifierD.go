package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveD(n, k int) float64 {
	d := make([]float64, 711)
	for i := 0; i < n; i++ {
		for t := 1; t < 710; t++ {
			term1 := float64(k-1) / float64(k) * d[t]
			coef := 1.0 / (float64(k) * float64(t+1))
			term2 := float64(t)*d[t] + d[t+1] + float64(t*(t+3))/2.0
			d[t] = term1 + coef*term2
		}
	}
	return d[1] * float64(k)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCases(path string) ([][2]int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	var cases [][2]int
	for {
		var n, k int
		if !scan.Scan() {
			break
		}
		n, _ = strconv.Atoi(scan.Text())
		if !scan.Scan() {
			return nil, fmt.Errorf("bad file")
		}
		k, _ = strconv.Atoi(scan.Text())
		cases = append(cases, [2]int{n, k})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesD.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		expected := solveD(tc[0], tc[1])
		in := fmt.Sprintf("%d %d\n", tc[0], tc[1])
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		scan := bufio.NewScanner(strings.NewReader(out))
		scan.Split(bufio.ScanWords)
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: no output\n", idx+1)
			os.Exit(1)
		}
		got, err := strconv.ParseFloat(scan.Text(), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx+1)
			os.Exit(1)
		}
		if scan.Scan() {
			fmt.Fprintf(os.Stderr, "case %d: extra output\n", idx+1)
			os.Exit(1)
		}
		if diff := got - expected; diff < -1e-6 || diff > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
