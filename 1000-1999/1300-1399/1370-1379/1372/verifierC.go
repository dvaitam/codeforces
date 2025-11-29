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

const testcasesRaw = `100
1
1
2
2 1
6
6 1 4 5 3 2
10
10 2 8 9 5 4 6 7 3 1
9
6 9 2 4 3 8 7 1 5
7
4 3 1 6 7 2 5
1
1
3
3 1 2
3
1 2 3
9
1 2 6 8 5 7 4 3 9
6
1 5 3 6 2 4
9
1 2 6 7 5 9 3 8 4
8
3 1 2 4 7 8 5 6
6
1 4 5 3 2 6
5
3 1 2 4 5
8
2 7 4 3 1 5 8 6
4
2 3 4 1
5
4 1 3 2 5
9
6 8 9 4 1 7 2 5 3
1
1
6
4 1 6 5 2 3
2
2 1
2
2 1
1
1
1
1
6
4 1 6 5 2 3
9
4 3 6 8 2 9 5 7 1
6
6 4 2 3 1 5
1
1
5
5 3 2 1 4
10
5 2 9 7 10 6 1 4 8 3
3
1 2 3
8
1 7 4 5 8 3 2 6
1
1
9
4 7 6 5 2 8 9 1 3
8
3 8 7 2 1 6 5 4
2
1 2
2
2 1
10
2 8 7 1 3 4 9 5 6 10
7
6 3 7 5 1 2 4
2
2 1
3
2 3 1
4
4 2 3 1
8
7 1 3 2 4 6 8 5
7
3 7 2 6 1 5 4
2
1 2
6
4 2 3 6 5 1
5
4 2 1 5 3
2
2 1
5
3 4 5 1 2
7
5 1 3 2 7 4 6
1
1
5
4 5 2 3 1
4
2 3 1 4
10
10 5 9 1 4 3 2 8 7 6
2
2 1
10
3 2 5 8 9 1 10 4 7 6
5
5 2 1 4 3
5
1 2 3 5 4
8
1 8 3 7 6 2 4 5
9
3 8 4 1 9 6 5 2 7
3
2 1 3
7
5 3 4 2 6 7 1
4
4 1 2 3
9
3 6 7 1 4 9 8 2 5
9
8 1 9 6 7 4 2 3 5
8
2 7 3 4 5 1 6 8
7
6 7 2 3 4 1 5
7
1 2 5 3 4 7 6
2
2 1
9
5 3 2 1 9 6 8 7 4
9
9 2 5 4 6 1 7 3 8
4
3 4 2 1
5
3 1 4 2 5
1
1
10
6 1 3 4 2 10 9 5 8 7
2
1 2
1
1
8
1 6 8 3 4 5 7 2
9
7 6 4 2 3 5 9 1 8
3
2 3 1
8
5 8 1 2 3 6 4 7
9
6 2 7 4 5 9 3 1 8
4
4 1 2 3
8
7 4 2 8 6 3 1 5
7
3 6 4 2 7 1 5
9
7 8 3 6 4 9 1 5 2
1
1
2
2 1
7
4 1 5 6 2 3 7
6
1 5 6 2 3 4
5
1 5 3 2 4
4
4 1 2 3
2
2 1
8
5 1 8 2 3 6 7 4
3
3 2 1
9
1 2 4 5 3 6 8 9 7
7
4 5 2 1 3 6 7
7
1 7 6 5 3 2 4
2
2 1
`

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(arr []int) int {
	n := len(arr)
	sorted := true
	for i := 0; i < n; i++ {
		if arr[i] != i+1 {
			sorted = false
			break
		}
	}
	if sorted {
		return 0
	}
	l := 0
	for l < n && arr[l] == l+1 {
		l++
	}
	r := n - 1
	for r >= 0 && arr[r] == r+1 {
		r--
	}
	needTwo := false
	for i := l; i <= r; i++ {
		if arr[i] == i+1 {
			needTwo = true
			break
		}
	}
	if needTwo {
		return 2
	}
	return 1
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][]int, t)
	for idx := 0; idx < t; idx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = v
		}
		cases[idx] = arr
	}
	input := []byte(testcasesRaw + "\n")
	out, err := runCandidate(bin, input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(strings.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for idx, arr := range cases {
		expect := solveCase(arr)
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(outScan.Text())
		if got != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
