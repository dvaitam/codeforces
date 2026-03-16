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
const testcasesBRaw = `
3 9 3 9 4
8 7 16 13 7 4 16 1 13 14
10 0 15 9 8 19 4 11 1 1 1 18
1 6 7
7 0 17 8 15 16 18 8 12
4 10 8 15 10 1
7 8 4 6 10 4 11 17 14
9 10 7 10 10 19 16 17 13 19 2
8 3 13 14 6 12 18 12 3 15
9 1 6 17 13 12 16 1 16 2 10
10 9 19 13 6 6 17 8 1 7 18 18
4 6 17 12 19 12
8 4 18 20 1 13 17 5 17 18
4 6 2 16 12 19
9 3 17 14 16 12 14 12 1 18 18
10 9 11 15 20 1 8 6 18 19 6 3
9 4 2 3 3 1 15 1 9 8 9
2 9 6 12
5 1 6 6 9 17 6
5 10 10 15 11 16 16
2 0 10 13
6 6 7 9 4 9 17 7
10 6 1 8 1 13 5 2 6 15 17 14
9 3 17 15 8 17 1 13 19 11 14
1 4 5
4 0 10 3 3 10
5 2 14 19 9 5 1
9 0 19 7 19 15 6 20 17 2 13
4 5 4 7 19 14
10 3 16 4 13 10 17 16 1 11 20 13
5 0 6 7 11 19 5
6 6 7 9 4 13 18 12
9 7 18 8 3 2 3 5 6 6 18
4 4 11 20 17 9
6 5 11 4 10 8 20 16
3 9 18 4 11
1 6 3
7 2 5 11 4 20 19 13 3
10 8 8 19 3 9 12 10 19 18 4 15
5 1 2 10 1 20 1
2 6 4 2
4 3 19 14 6 4
8 2 8 6 4 14 13 18 10 18
5 7 11 4 7 11 2
1 0 10
10 5 15 13 11 13 3 3 11 20 15 4
5 3 20 18 16 12 9
3 8 7 10 7
4 5 3 9 3 15
2 10 19 11
4 6 10 2 11 6
6 9 10 8 11 4 18 20
10 9 3 8 8 1 8 13 3 9 18 3
2 0 1 10
6 7 16 5 4 17 11 3
9 10 6 6 5 5 11 10 4 17 20
5 2 7 5 18 2 11
10 10 18 7 6 10 14 18 6 2 8 9
2 10 15 14
9 4 18 15 18 15 1 13 11 6 9
8 0 14 19 1 2 12 19 5 19
3 2 9 9 13
10 6 6 20 3 8 16 1 6 17 11 17
8 10 8 8 11 16 16 8 14 11
9 9 9 8 2 3 17 12 6 17 7
5 4 10 18 12 6 15
10 1 4 20 17 19 13 6 5 9 14 7
10 0 16 13 12 13 17 6 18 2 17 3
5 10 4 9 3 5 20
2 7 8 13
7 6 6 11 15 5 20 16 7
2 6 20 18
7 1 10 9 8 13 18 1 7
9 7 19 1 1 20 8 9 7 6 10
3 8 7 9 10
10 4 15 6 18 12 16 14 4 7 19 13
4 4 4 1 4 19
1 8 10
3 1 17 12 19
5 6 17 12 17 11 1
2 7 15 12
5 8 13 11 19 16 4
7 6 7 18 1 9 20 17 7
8 9 17 14 10 6 15 20 17 7
6 8 1 13 19 14 13 11
10 9 3 16 8 10 1 14 5 13 9 6
2 9 1 12
5 6 18 10 5 15 9
8 2 15 17 2 9 17 4 19 14
2 5 3 15
1 2 17
3 1 13 9 20
5 3 17 7 8 11 9
2 1 17 12
8 8 18 2 6 10 18 9 12 20
4 6 18 13 6 16
5 9 11 8 9 20 8
1 9 13
6 6 8 9 7 3 6 19
8 9 5 20 9 15 17 6 5 5
`


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

	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
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
