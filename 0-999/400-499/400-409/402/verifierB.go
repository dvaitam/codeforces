package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierB.go /path/to/binary\n")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refB")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "402B.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n%s", err, out)
		os.Exit(1)
	}
	defer os.Remove(ref)

	f, err := os.Open(filepath.Join(dir, "testcasesB.txt"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	idx := 0
	for {
		if !scan.Scan() {
			break
		}
		line1 := strings.TrimSpace(scan.Text())
		if line1 == "" {
			continue
		}
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "unexpected EOF in testcases\n")
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scan.Text())
		idx++
		input := line1 + "\n" + line2 + "\n"
		// compute expected minimal operations count
		parts := strings.Fields(line1)
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "bad test line: %s\n", line1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		arrStr := strings.Fields(line2)
		if len(arrStr) != n {
			fmt.Fprintf(os.Stderr, "bad array line for test %d\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := range arr {
			v, _ := strconv.Atoi(arrStr[i])
			arr[i] = v
		}
		freq := map[int]int{}
		for i := 0; i < n; i++ {
			h := arr[i] - i*k
			if h >= 1 {
				freq[h]++
			}
		}
		bestCnt := 0
		for _, c := range freq {
			if c > bestCnt {
				bestCnt = c
			}
		}
		expectedOps := n - bestCnt

		candOut, cErr := runBinary(candidate, input)
		if cErr != nil {
			fmt.Fprintf(os.Stderr, "test %d: candidate error: %v\n", idx, cErr)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(candOut), "\n")
		if len(lines) == 0 {
			fmt.Fprintf(os.Stderr, "test %d: empty output\n", idx)
			os.Exit(1)
		}
		p, err := strconv.Atoi(strings.Fields(lines[0])[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid count line\n", idx)
			os.Exit(1)
		}
		if p != expectedOps {
			fmt.Fprintf(os.Stderr, "test %d: expected %d operations got %d\n", idx, expectedOps, p)
			os.Exit(1)
		}
		if len(lines)-1 != p {
			fmt.Fprintf(os.Stderr, "test %d: expected %d operation lines, got %d\n", idx, p, len(lines)-1)
			os.Exit(1)
		}
		// apply operations
		for i := 1; i <= p; i++ {
			parts := strings.Fields(lines[i])
			if len(parts) != 3 {
				fmt.Fprintf(os.Stderr, "test %d: bad operation line %d\n", idx, i)
				os.Exit(1)
			}
			op := parts[0]
			j, _ := strconv.Atoi(parts[1])
			x, _ := strconv.Atoi(parts[2])
			if j < 1 || j > n {
				fmt.Fprintf(os.Stderr, "test %d: invalid index %d\n", idx, j)
				os.Exit(1)
			}
			if op == "+" {
				arr[j-1] += x
			} else if op == "-" {
				arr[j-1] -= x
			} else {
				fmt.Fprintf(os.Stderr, "test %d: invalid operation %s\n", idx, op)
				os.Exit(1)
			}
		}
		// check arithmetic progression
		for i := 1; i < n; i++ {
			if arr[i]-arr[i-1] != k {
				fmt.Fprintf(os.Stderr, "test %d: result not an AP\n", idx)
				os.Exit(1)
			}
		}
	}
	if err := scan.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
