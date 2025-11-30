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

// embeddedTestcasesA matches the previous contents of testcasesA.txt.
const embeddedTestcasesA = `100
7 7 1
5 9 8
7 5 8
6 10 4
9 3 5
3 2 10
5 9 10
3 5 2
2 6 8
9 2 6
7 6 10
4 9 8
8 9 5
1 9 1
2 7 1
10 8 6
4 6 2
4 10 4
4 3 9
8 2 2
6 9 8
2 5 9
5 2 9
6 9 4
10 9 10
5 8 2
10 7 6
10 4 5
3 4 3
1 10 5
8 2 2
3 3 1
2 9 7
9 5 9
4 4 10
7 10 5
8 8 6
2 6 10
2 8 10
6 4 4
1 5 2
4 6 3
6 7 1
2 3 4
1 10 9
10 2 1
2 4 10
10 2 7
2 6 2
1 10 1
4 3 2
8 4 1
1 9 7
10 2 5
2 4 2
5 6 7
3 1 9
8 1 10
2 7 4
5 6 8
10 3 4
1 3 3
6 9 5
2 10 8
3 1 8
7 10 9
5 6 7
5 3 9
1 8 2
6 1 9
5 3 4
8 6 10
5 6 10
10 3 5
7 7 2
1 10 4
6 3 4
4 8 7
10 7 1
7 10 7
1 3 8
2 5 3
8 9 8
9 10 1
1 8 6
5 8 1
7 4 9
2 3 1
7 7 6
1 4 1
1 9 10
2 4 2
10 4 5
5 3 2
8 7 2
1 5 8
2 5 3
9 6 2
3 5 1
1 1 4
`

// maxProduct mirrors the logic from 1992A.go to compute the answer for one test.
func maxProduct(a, b, c int) int {
	best := 0
	for i := 0; i <= 5; i++ {
		for j := 0; j <= 5-i; j++ {
			k := 5 - i - j
			prod := (a + i) * (b + j) * (c + k)
			if prod > best {
				best = prod
			}
		}
	}
	return best
}

func expectedFromInput(input string) ([]string, error) {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("invalid test count: %w", err)
	}
	results := make([]string, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("not enough values in test data for case %d", i+1)
		}
		a, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid a for case %d: %w", i+1, err)
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("not enough values in test data for case %d", i+1)
		}
		b, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid b for case %d: %w", i+1, err)
		}
		if !scanner.Scan() {
			return nil, fmt.Errorf("not enough values in test data for case %d", i+1)
		}
		c, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, fmt.Errorf("invalid c for case %d: %w", i+1, err)
		}
		results[i] = fmt.Sprintf("%d", maxProduct(a, b, c))
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	expected, err := expectedFromInput(embeddedTestcasesA)
	if err != nil {
		fmt.Println("could not compute expected outputs:", err)
		os.Exit(1)
	}

	cmd := exec.Command(os.Args[1])
	cmd.Stdin = strings.NewReader(embeddedTestcasesA)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out.String())
		os.Exit(1)
	}

	outScan := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < len(expected); i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
