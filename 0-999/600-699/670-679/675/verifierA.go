package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func expected(a, b, c int64) string {
	if a == b {
		return "YES"
	}
	if c == 0 {
		return "NO"
	}
	if (b-a)%c == 0 && (b-a)/c >= 0 {
		return "YES"
	}
	return "NO"
}

const testcasesARaw = `100
-12 16 -4
-4 -13 2
8 10 5
4 -7 -4
11 -19 1
7 18 -5
8 -3 -2
17 -14 0
-19 -19 -5
14 -20 1
-7 7 -5
13 -6 2
11 15 -2
2 -6 5
-6 9 -1
-19 6 3
-14 -9 5
-2 -13 0
12 7 3
-8 -1 -1
17 11 3
5 17 -5
10 -5 1
6 -9 0
15 3 -4
8 12 -4
-10 13 1
3 11 -5
10 -18 -1
19 17 4
5 -10 -3
12 -6 -5
-8 14 3
-6 5 3
2 16 0
9 -3 5
15 18 -5
4 12 -3
13 15 -2
7 -17 2
3 16 3
-8 12 1
11 2 1
2 -20 3
14 19 4
1 9 4
-19 -6 5
-9 15 4
-9 -15 3
-4 -18 5
-16 -15 -5
8 -20 -1
-5 -3 -4
19 -9 0
-2 -16 -3
-10 -4 3
-10 -3 5
-2 9 0
11 10 -4
-19 -1 1
1 6 -2
-4 -14 -1
12 -7 4
7 -19 -2
-19 5 -3
-18 -10 2
12 7 3
-6 20 3
8 -6 3
-19 5 5
16 0 5
20 7 -5
-1 -12 -2
-17 -1 -4
-16 -1 -1
-10 6 4
-4 -12 -5
15 -18 4
-7 16 2
-10 19 3
-18 4 -2
2 -14 -2
16 7 4
-8 11 -4
4 -2 3
11 -19 0
19 5 -1
-19 -10 -2
0 16 -3
1 7 -2
-3 -14 1
15 2 5
14 11 3
-5 -16 -5
-15 -12 -3
-10 14 -2
-3 1 4
12 -4 0
1 1 -4
-2 -5 4`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data := []byte(testcasesARaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expectedOut := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		aVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		bVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		cVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		expectedOut[i] = expected(aVal, bVal, cVal)
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("execution failed: %v\n%s", err, out)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expectedOut[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expectedOut[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
