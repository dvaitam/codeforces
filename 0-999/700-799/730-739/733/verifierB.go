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

const testcasesBRaw = `100
6
2 4
4 1
3 1
2 4
2 5
5 2
3
2 3
5 2
2 1
8
4 4
3 4
4 2
4 2
5 3
5 4
3 5
2 1
8
1 3
4 2
5 5
5 4
5 2
2 3
4 1
3 3
5
4 4
2 2
4 1
5 5
2 1
4
2 2
4 4
5 3
5 2
8
2 2
2 2
3 1
3 4
2 3
1 5
2 3
2 4
8
2 2
3 3
5 4
2 5
3 4
4 4
5 1
3 4
2
3 3
4 4
4
3 2
5 3
3 5
1 5
8
5 2
5 3
3 4
4 2
5 1
4 1
4 2
2 2
7
1 4
1 4
1 3
5 5
4 4
3 2
3 1
8
1 4
4 5
2 4
5 4
5 5
2 5
5 4
2 2
3
4 1
2 1
2 3
7
3 2
1 2
5 2
2 2
2 5
4 4
5 5
8
1 2
3 5
2 4
4 4
1 2
5 4
3 4
1 5
7
5 2
4 3
4 3
1 4
3 3
2 2
4 1
4
3 3
3 1
3 1
5 4
1
3 1
1
5 4
5
1 1
3 4
1 5
1 5
2 4
4
3 3
3 4
2 4
1 3
7
1 3
4 2
5 3
3 4
1 1
5 4
4 1
2
1 3
4 3
6
1 5
4 5
1 2
3 5
1 1
2 4
7
4 1
4 3
1 3
2 3
5 1
5 2
5 4
8
5 3
1 3
3 1
3 2
4 3
4 2
3 5
2 3
2
1 1
1 5
5
1 2
2 2
3 2
3 5
4 1
5
5 3
4 2
5 3
4 1
1 4
1
5 3
3
1 3
5 2
4 1
3
3 1
1 5
4 2
5
3 3
3 3
5 5
1 5
4 1
4
2 3
3 4
1 4
3 4
2
2 5
4 5
1
1 5
6
5 5
4 5
5 4
2 3
5 4
3 1
3
5 5
4 1
3 5
6
5 3
3 4
4 5
3 4
5 5
5 1
1
1 3
4
5 2
5 4
4 5
3 5
4
5 2
2 4
2 3
3 1
2
3 4
3 1
8
1 5
5 4
2 3
3 3
2 3
2 2
5 2
4 5
3
4 2
3 5
1 2
5
2 1
4 5
3 2
2 3
3 4
4
2 5
1 4
3 5
4 3
6
1 5
1 1
2 1
4 5
1 3
1 3
1
3 5
7
3 2
3 4
4 4
5 5
5 2
5 4
3 2
3
1 1
2 5
3 3
5
3 1
1 2
2 2
1 1
5 5
7
3 2
1 5
2 2
2 1
2 4
4 2
2 2
2
3 2
4 1
6
3 4
3 4
3 5
4 5
1 1
1 4
5
2 5
2 2
2 3
2 1
4 3
6
5 1
1 4
2 2
1 4
2 5
3 5
6
2 3
5 4
4 1
4 3
3 1
3 3
2
3 3
3 1
4
1 3
5 5
2 3
1 3
5
5 4
2 4
1 5
4 1
2 1
1
5 3
2
1 1
1 1
8
1 1
3 5
5 3
5 5
4 2
3 3
3 4
4 5
1
2 5
7
2 2
2 5
3 2
3 2
5 2
4 3
1 1
7
5 5
4 3
3 4
5 2
3 2
4 2
5 1
2
1 3
3 5
6
3 1
1 2
5 3
5 2
1 4
4 2
4
5 4
2 5
1 1
1 2
3
4 3
5 3
1 4
3
1 5
2 1
4 1
8
3 2
5 2
1 4
2 2
1 5
5 1
4 3
5 2
7
4 4
4 1
5 1
1 2
1 1
1 1
5 1
7
3 1
2 1
2 4
4 5
4 4
4 5
5 4
3
1 3
4 4
3 5
3
5 5
3 5
2 2
5
2 4
5 2
1 1
2 1
4 2
1
1 1
5
4 2
2 5
5 3
5 3
1 1
6
5 3
5 2
2 1
4 3
2 3
2 4
3
4 1
5 5
2 5
2
3 4
2 4
5
1 3
2 4
3 5
4 5
4 5
5
5 1
3 2
3 3
3 1
4 3
2
5 1
2 2
2
4 3
2 3
4
5 4
3 1
1 5
3 3
6
5 1
3 5
3 3
3 3
4 5
2 4
6
5 4
2 3
5 3
2 5
2 2
4 1
6
3 2
1 1
2 4
2 5
5 1
3 5
4
4 5
3 1
3 3
3 5
4
3 5
4 1
4 4
3 5
6
3 2
1 4
4 5
1 5
5 2
4 2
3
3 3
3 5
3 2
6
4 3
2 4
5 4
2 4
2 2
2 4
7
4 4
2 5
3 3
1 4
2 1
5 1
5 3
8
2 4
3 2
3 3
2 3
3 3
3 5
1 3
1 4
8
4 5
4 4
1 1
1 3
3 5
1 3
3 3
3 2
`

func solveCase(n int, l, r []int) string {
	sumL := 0
	sumR := 0
	for i := 0; i < n; i++ {
		sumL += l[i]
		sumR += r[i]
	}
	diff := abs(sumL - sumR)
	best := diff
	bestIdx := 0
	for i := 0; i < n; i++ {
		newL := sumL - l[i] + r[i]
		newR := sumR - r[i] + l[i]
		d := abs(newL - newR)
		if d > best {
			best = d
			bestIdx = i + 1
		}
	}
	return fmt.Sprintf("%d\n", bestIdx)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		l := make([]int, n)
		r := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			l[j], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			r[j], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", l[j], r[j]))
		}
		exp := solveCase(n, l, r)
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
