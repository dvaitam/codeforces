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
const testcasesDRaw = `
7 5 5 7 1 5 9 7 1 7 5 7 4 6
4 2 2 10 6 3 1 7
8 5 8 2 2 7 9 9 3 3 3 4 3 4 1
6 5 8 6 10 5 6 2 7 5 3 6 9
5 3 7 9 5 4 7 6 7 8
7 5 7 7 2 3 3 1 10 10 9 2 4 10
9 5 2 5 10 3 7 2 1 1 2 6 8 6 2 8
7 3 8 4 3 9 9 2 9 3 1 3
8 5 10 4 8 7 5 1 10 3 7 3 8 10 1
10 4 2 10 7 6 4 9 8 1 8 10 2 5 9 8
10 5 7 8 5 3 4 9 6 3 5 10 3 8 2 2 8
8 4 7 9 2 5 8 4 2 5 3 6 2 3
3 1 10 9 4 1
7 1 9 10 8 2 8 9 6 6
1 1 4 4
10 4 5 5 4 1 8 6 9 6 2 2 5 10 7 4
9 3 3 4 5 4 8 6 5 7 10 3 2 7
10 3 4 6 6 7 5 6 7 5 2 8 5 2 8
7 2 4 3 6 8 4 2 7 7 8
8 5 10 10 4 7 9 5 8 4 6 9 1 2 8
9 3 4 7 1 10 1 7 2 5 4 6 3 2
7 2 1 4 1 1 7 9 1 3 2
9 5 4 2 8 4 1 9 10 7 2 9 3 4 4 7
7 4 1 7 4 7 1 10 5 1 10 6 6
10 3 3 10 9 2 5 2 2 5 1 3 10 3 7
7 2 4 7 9 9 2 9 2 8 2
8 5 8 3 8 9 6 3 7 5 7 2 10 9 6
9 2 4 6 8 7 1 8 1 9 7 8 4
5 4 5 8 8 3 4 8 5 6 8
6 5 9 2 4 5 9 2 1 3 6 1 6
6 5 8 7 4 7 8 4 3 7 1 6 9
7 2 3 9 9 10 10 7 3 10 8
6 5 10 8 1 4 8 10 6 9 4 1 1
3 1 4 8 10 4
9 1 7 5 6 8 4 3 6 9 6 8
7 4 6 8 2 3 4 2 3 7 8 9 4
3 3 1 7 4 2 7 10
8 5 9 4 1 3 9 8 5 7 9 9 2 2 10
10 4 1 4 7 10 6 5 2 1 4 10 1 10 9 6
9 2 5 2 10 8 5 4 7 7 9 6 4
7 5 2 1 3 6 10 10 8 2 9 8 7 6
9 4 1 2 7 10 3 3 1 2 7 4 1 5 5
4 4 6 5 6 10 4 9 9 1
3 3 3 8 1 5 6 3
4 3 10 9 7 2 5 10 2
5 5 10 3 8 1 9 4 6 2 9 4
5 2 7 4 10 8 6 3 5
6 3 4 8 6 3 5 9 1 2 6
8 3 1 2 10 1 5 2 3 5 4 9 9
10 3 2 7 9 9 7 2 2 10 8 6 3 5 4
9 4 5 1 1 1 5 2 6 6 9 5 1 4 9
6 2 10 1 7 7 5 10 4 5
10 3 8 2 7 6 3 9 1 4 8 7 1 4 2
4 3 3 3 8 1 6 9 7
9 1 4 8 5 5 7 7 8 1 3 9
6 3 4 3 4 5 3 3 9 5 9
3 1 4 7 3 4
6 5 10 3 4 9 5 1 8 4 1 7 3
8 4 7 4 7 9 9 8 2 3 2 9 1 2
4 1 7 8 5 7 6
6 4 2 3 10 4 4 1 1 9 2 8
9 2 3 10 9 5 1 10 3 6 6 6 5
10 5 5 3 1 8 4 10 7 6 8 5 8 3 8 5 6
9 3 7 10 7 9 5 10 4 4 7 4 5 9
10 2 3 10 9 9 9 8 6 5 7 5 8 9
4 1 5 8 2 1 4
5 5 5 1 4 3 10 3 6 10 5 1
1 1 10 3
6 3 7 8 6 6 2 10 7 2 1
2 2 10 9 5 6
8 2 4 6 3 3 8 3 10 10 9 7
10 5 4 1 1 7 10 1 3 9 2 7 5 6 10 4 1
5 3 6 5 8 8 10 6 5 2
10 4 3 2 9 5 9 7 10 8 7 3 10 7 3 2
5 5 3 8 5 2 7 5 5 2 6 9
5 5 1 1 10 1 3 2 7 7 8 1
5 3 2 7 10 4 1 4 7 1
10 1 7 3 9 5 7 7 9 1 8 9 3
6 5 4 10 3 2 10 8 5 5 5 4 7
1 1 9 2
5 1 3 3 9 3 2 6
6 1 3 6 10 2 2 6 4
10 3 6 8 8 9 9 2 2 6 8 7 7 9 8
8 5 9 2 4 3 7 9 6 3 2 3 3 1 7
5 5 2 9 5 10 5 6 8 7 10 6
8 5 9 6 4 9 1 7 6 6 9 1 4 10 2
4 4 6 7 10 4 9 1 9 9
1 1 10 6
6 1 3 6 5 8 8 7 1
9 2 6 3 1 8 9 3 7 5 10 8 4
8 1 10 10 3 10 8 6 10 6 7
10 5 3 9 7 7 4 6 7 5 4 4 5 4 9 10 6
5 4 10 6 10 3 6 4 3 5 9
10 2 10 2 8 8 4 3 10 6 10 4 2 9
5 4 5 7 5 9 4 6 7 6 1
4 1 4 10 3 9 8
5 1 9 5 8 7 6 7
2 1 1 6 3
10 3 5 9 7 8 6 4 2 2 2 9 1 10 9
4 3 10 4 1 1 2 5 2
8 5 3 10 10 9 4 5 5 3 3 10 4 3 3
6 2 8 5 3 10 5 4 2 3
5 4 4 10 1 6 9 8 3 5 2
10 5 5 10 4 1 3 4 10 1 7 3 3 8 9 1 9
6 3 3 8 6 8 4 5 4 2 6
6 4 6 9 10 5 4 5 9 7 2 3
4 4 10 2 9 3 8 3 7 5
6 5 3 2 5 8 6 7 9 1 6 7 7
6 5 9 10 3 1 5 10 8 9 10 7 4
5 4 7 2 2 9 7 2 9 8 3
4 1 4 4 2 10 7
4 4 9 4 9 4 4 8 2 4
9 1 3 10 4 9 10 9 10 7 10 7
10 3 10 1 6 6 4 1 4 8 4 9 4 1 3
9 2 8 2 10 5 5 7 6 5 2 10 2
6 4 10 10 2 4 3 7 4 1 8 5
10 2 6 3 2 9 6 6 2 6 4 4 5 9
9 4 7 7 1 3 8 2 5 3 9 4 10 8 6
4 2 10 7 5 4 2 3
9 2 9 7 3 1 5 8 10 2 6 2 2
`


func expectedD(n, w int, a, b []int64) int {
	if w == 1 {
		return n
	}
	m := w - 1
	pattern := make([]int64, m)
	for i := 0; i < m; i++ {
		pattern[i] = b[i+1] - b[i]
	}
	text := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		text[i] = a[i+1] - a[i]
	}
	pi := make([]int, m)
	for i, j := 1, 0; i < m; i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = pi[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		pi[i] = j
	}
	count := 0
	for i, j := 0, 0; i < n-1; i++ {
		for j > 0 && text[i] != pattern[j] {
			j = pi[j-1]
		}
		if text[i] == pattern[j] {
			j++
		}
		if j == m {
			count++
			j = pi[j-1]
		}
	}
	return count
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesDRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		wVal, _ := strconv.Atoi(parts[1])
		need := 2 + nVal + wVal
		if len(parts) != need {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		a := make([]int64, nVal)
		b := make([]int64, wVal)
		for i := 0; i < nVal; i++ {
			v, _ := strconv.ParseInt(parts[2+i], 10, 64)
			a[i] = v
		}
		for i := 0; i < wVal; i++ {
			v, _ := strconv.ParseInt(parts[2+nVal+i], 10, 64)
			b[i] = v
		}
		expect := expectedD(nVal, wVal, a, b)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", nVal, wVal))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, v := range b {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(gotStr)
		if err != nil {
			fmt.Printf("test %d: invalid output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
