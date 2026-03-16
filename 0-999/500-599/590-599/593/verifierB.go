package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type line struct {
	y1 int64
	y2 int64
}

func solve(n int, x1, x2 int64, ks, bs []int64) string {
	lines := make([]line, n)
	for i := 0; i < n; i++ {
		lines[i].y1 = ks[i]*x1 + bs[i]
		lines[i].y2 = ks[i]*x2 + bs[i]
	}
	sort.Slice(lines, func(i, j int) bool {
		if lines[i].y1 == lines[j].y1 {
			return lines[i].y2 < lines[j].y2
		}
		return lines[i].y1 < lines[j].y1
	})
	for i := 0; i < n-1; i++ {
		if lines[i].y2 > lines[i+1].y2 {
			return "Yes"
		}
	}
	return "No"
}

const testcasesBRaw = `100
4 -1 5
0 5
1 8
0 9
2 6
4 -7 9
-5 2
0 0
3 -10
-2 8
2 -3 2
0 3
4 -3
6 -10 9
3 -7
-5 8
0 -9
0 -5
0 -2
2 5
8 -4 1
5 -3
5 10
2 1
5 -5
-5 5
3 -5
1 -1
5 8
5 -7 1
2 9
4 -5
4 4
0 8
0 3
6 0 3
3 3
1 -3
0 2
-1 -3
-2 -2
-2 7
6 -5 8
-2 9
0 9
5 -6
0 -3
-3 -2
3 -5
6 -1 2
2 5
-2 -5
4 1
4 -3
-4 7
2 -4
4 -4 5
3 -4
0 3
0 -10
-5 -4
3 -1 3
3 4
-4 -2
3 -1
7 -3 7
0 2
-5 -7
-5 0
-3 -2
4 7
2 -1
3 3
8 -3 6
-2 -8
-4 10
-4 6
-5 8
3 0
4 -7
4 -3
-5 -9
8 -9 5
-4 8
1 -6
4 9
3 -4
-1 -1
-3 -5
4 -3
-3 -6
3 -6 8
-1 -6
4 -2
-2 -8
3 0 7
-3 -4
4 -10
0 7
7 -3 3
-1 -8
4 -2
4 -3
-5 -6
-2 8
5 -10
1 -6
3 -3 9
-5 6
0 -3
5 -5
5 -5 8
-2 10
5 -2
5 -4
-4 1
2 -2
4 -2 10
-1 -4
-1 7
1 -8
-3 2
2 -1 2
3 -6
2 -6
5 -6 9
-2 -5
-4 -2
2 0
1 6
-5 1
8 -4 7
3 -6
3 8
-3 10
-5 -1
4 4
3 8
1 -4
2 -4
5 -9 2
3 2
-4 6
5 4
5 -5
-5 -4
3 -7 2
4 2
1 -8
3 -4
6 -7 4
2 -8
-2 -7
3 -6
4 -4
3 8
3 10
6 -4 10
1 10
-1 5
5 10
4 -10
0 1
-2 10
2 -7 8
-2 6
5 -4
8 -10 3
-2 4
2 -8
4 -2
5 -1
1 0
5 -1
5 -8
4 -7
8 -3 6
4 9
5 -8
-2 -6
1 -3
-5 5
-1 -6
3 10
0 -6
2 -2 6
0 -8
0 -4
3 -6 9
1 10
1 2
-5 -6
5 -9 5
5 -6
3 -4
4 5
0 -6
4 -6
6 -4 5
-2 -7
-3 -6
2 7
1 -5
-5 -1
2 8
2 0 6
5 -7
-4 5
3 -2 1
-5 8
-4 4
-1 8
4 -2 7
-3 -3
5 -7
5 8
3 -5
5 0 4
-5 -3
-3 9
5 -6
5 -6
4 9
3 -9 2
-1 -7
-1 9
2 8
8 -6 10
1 -9
-4 -6
2 -9
0 0
-1 10
1 2
-5 -1
5 -3
3 -9 10
-3 -3
3 -4
2 0
8 -10 5
1 8
3 7
-3 0
4 10
3 2
-2 -1
-4 4
5 10
4 -7 5
-1 -8
0 -6
2 -10
1 10
5 -10 5
3 9
1 4
4 8
3 -4
-3 -10
4 -9 7
-1 -4
2 5
-1 6
1 -6
6 -3 5
-3 -10
-2 0
-1 -7
5 -10
1 -1
0 9
7 -7 6
1 -6
5 -7
-2 1
2 4
-3 1
1 3
1 1
6 -9 2
3 4
-3 -7
-5 6
1 10
2 -3
1 10
7 -9 7
-5 -1
4 -5
-1 -2
3 -7
3 0
-5 -6
-5 -8
4 -5 8
-2 -3
0 8
2 -2
3 -5
5 -4 1
-4 -10
-4 0
-2 -1
3 -7
3 -1
2 -9 9
4 -1
-5 -3
8 -9 3
2 -5
-3 4
-1 -5
2 8
-5 7
1 -10
-4 -7
-1 -8
7 -4 4
-1 2
-2 4
5 9
3 7
0 -5
1 10
-5 8
5 -1 4
-1 1
3 0
2 10
2 -2
0 -2
8 -8 3
-2 4
0 10
-5 1
2 9
-1 -3
-5 -7
2 6
-3 -4
8 -4 10
0 7
-3 -9
2 3
2 4
0 -10
-2 8
-2 -8
2 -6
7 -9 9
5 7
1 9
-1 -4
-4 3
-5 7
2 0
-5 -4
6 -6 7
2 6
-1 -8
-5 -8
0 2
5 3
1 10
7 0 4
-1 5
1 7
-5 10
1 10
-1 5
-3 2
5 3
5 -1 10
-1 5
3 1
-1 3
4 10
-3 1
2 -4 8
2 -1
-4 -3
4 -7 7
-5 9
0 -6
-3 5
-4 -7
3 -1 1
-2 10
-3 -10
5 6
7 -4 7
-3 -1
1 4
-1 7
-5 -5
2 -5
-4 -7
2 -1
7 -6 5
1 -9
4 -2
0 4
-4 -9
5 -1
-2 4
5 -1
6 -7 7
-4 3
-5 -9
-4 1
-5 -1
-1 -9
-5 -9
6 -10 3
4 -5
2 10
-4 8
-3 -10
3 -6
-4 -8
7 -9 10
5 -7
-2 -7
2 3
-4 -7
2 4
3 -5
3 -7
3 -5 7
3 -6
-4 10
3 -9
7 -2 4
-4 2
0 6
0 7
4 7
1 -1
-2 8
-4 -6
6 -4 6
-1 1
3 1
-5 -10
-5 -6
-1 -4
5 -6
3 -9 8
2 -6
-1 -2
-3 -9
3 0 6
0 -4
5 -7
3 5
3 0 5
0 -4
5 -7
4 -3
4 -4 2
-1 2
3 1
-5 10
0 -3
5 -10 7
3 -2
-3 -2
-5 -2
-4 0
5 3
2 0 5
-2 10
-3 -2
4 -2 2
-3 5
-2 8
4 8
-4 3
4 0 1
-1 -10
-5 9
-3 0
-3 2
3 -1 9
-2 -9
5 7
-2 4
6 0 8
4 9
0 1
0 -7
-5 3
0 -5
2 -3
7 0 7
4 -10
5 -4
4 5
1 2
-3 -6
5 3
1 4
8 0 4
-1 -8
-3 8
-5 -1
-3 -7
-3 -2
-3 2
-1 6
5 3
4 -2 7
2 4
-4 -2
4 0
-1 -2
7 -6 7
4 6
-5 5
5 6
3 -4
1 -2
5 3
-5 -10
5 -2 2
-1 -10
-5 5
5 2
3 1
3 5
6 -3 9
2 8
-3 8
-5 9
0 -4
-1 3
-4 -3
2 -8 10
2 10
-3 5
4 -3 3
2 -2
2 0
3 -3
-5 1
6 -6 7
3 10
4 -7
-1 -6
3 -6
-3 -8
0 -4
2 -10 8
2 -8
4 6
2 -8 8
4 -6
4 -2
2 -6 6
-3 4
0 -9
4 -9 2
4 -7
-1 -1
4 1
-3 -10
7 -3 7
1 10
-2 -8
1 -2
2 3
5 5
1 -6
1 -4
8 -4 4
-1 -8
-3 4
0 1
-1 4
0 -3
-3 4
-3 0
-2 5
3 -9 9
4 9
2 0
-3 -8
7 -5 3
-1 1
5 -7
-4 2
0 0
-1 0
-5 -4
5 10
4 0 3
-3 7
-3 5
-2 -8
3 -2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesBRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	type input struct {
		n      int
		x1, x2 int64
		ks, bs []int64
	}
	inputs := make([]input, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		x1, _ := strconv.ParseInt(scan.Text(), 10, 64)
		scan.Scan()
		x2, _ := strconv.ParseInt(scan.Text(), 10, 64)
		ks := make([]int64, n)
		bs := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			ks[j], _ = strconv.ParseInt(scan.Text(), 10, 64)
			scan.Scan()
			bs[j], _ = strconv.ParseInt(scan.Text(), 10, 64)
		}
		inputs[i] = input{n, int64(x1), int64(x2), ks, bs}
	}
	expected := make([]string, t)
	for i, in := range inputs {
		expected[i] = solve(in.n, in.x1, in.x2, in.ks, in.bs)
	}
	for i, in := range inputs {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", in.n, in.x1, in.x2))
		for j := 0; j < in.n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", in.ks[j], in.bs[j]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
