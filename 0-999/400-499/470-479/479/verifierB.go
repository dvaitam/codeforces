package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleB")
	cmd := exec.Command("go", "build", "-o", oracle, "479B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		if len(parts) != n+2 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		arrStr := make([]string, n)
		for i := 0; i < n; i++ {
			arrStr[i] = parts[2+i]
			val, _ := strconv.Atoi(parts[2+i])
			arr[i] = val
		}
		input := fmt.Sprintf("%d %d\n%s\n", n, k, strings.Join(arrStr, " "))

		// run oracle to get minimal achievable difference
		cmdO := exec.Command(oracle)
		cmdO.Stdin = strings.NewReader(input)
		var outO bytes.Buffer
		cmdO.Stdout = &outO
		if err := cmdO.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "oracle run error: %v\n", err)
			os.Exit(1)
		}
		var expectedDiff, dummy int
		if _, err := fmt.Fscan(bytes.NewReader(outO.Bytes()), &expectedDiff, &dummy); err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\n", err)
			os.Exit(1)
		}

		// run candidate solution
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

		reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
		var s, t int
		if _, err := fmt.Fscan(reader, &s, &t); err != nil {
			fmt.Printf("test %d: failed to read s and t: %v\n", idx, err)
			os.Exit(1)
		}
		if t < 0 || t > k {
			fmt.Printf("test %d: invalid number of operations %d\n", idx, t)
			os.Exit(1)
		}

		// apply operations
		h := make([]int, n)
		copy(h, arr)
		for op := 0; op < t; op++ {
			var i, j int
			if _, err := fmt.Fscan(reader, &i, &j); err != nil {
				fmt.Printf("test %d: failed to read operation %d: %v\n", idx, op+1, err)
				os.Exit(1)
			}
			if i < 1 || i > n || j < 1 || j > n {
				fmt.Printf("test %d: operation %d has invalid indices %d %d\n", idx, op+1, i, j)
				os.Exit(1)
			}
			if h[i-1] <= h[j-1] {
				fmt.Printf("test %d: operation %d moves from non-taller tower %d to %d\n", idx, op+1, i, j)
				os.Exit(1)
			}
			h[i-1]--
			h[j-1]++
		}

		// compute resulting difference
		minH, maxH := h[0], h[0]
		for i := 1; i < n; i++ {
			if h[i] < minH {
				minH = h[i]
			}
			if h[i] > maxH {
				maxH = h[i]
			}
		}
		diff := maxH - minH
		if diff != s {
			fmt.Printf("test %d: reported diff %d but got %d\n", idx, s, diff)
			os.Exit(1)
		}
		if s != expectedDiff {
			fmt.Printf("test %d: expected minimal diff %d, but got %d\n", idx, expectedDiff, s)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
