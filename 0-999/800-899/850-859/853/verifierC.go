package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveCase(line string) string {
	fields := strings.Fields(line)
	idx := 0
	if len(fields) < 2 {
		return ""
	}
	n, _ := strconv.Atoi(fields[idx])
	idx++
	q, _ := strconv.Atoi(fields[idx])
	idx++
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i], _ = strconv.Atoi(fields[idx])
		idx++
	}
	type Query struct{ l, d, r, u int }
	qs := make([]Query, q)
	for i := 0; i < q; i++ {
		l, _ := strconv.Atoi(fields[idx])
		idx++
		d, _ := strconv.Atoi(fields[idx])
		idx++
		r, _ := strconv.Atoi(fields[idx])
		idx++
		u, _ := strconv.Atoi(fields[idx])
		idx++
		qs[i] = Query{l, d, r, u}
	}
	var sb strings.Builder
	for _, qu := range qs {
		cnt := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				r1 := i + 1
				r2 := j + 1
				c1 := p[i]
				c2 := p[j]
				if r1 > r2 {
					r1, r2 = r2, r1
				}
				if c1 > c2 {
					c1, c2 = c2, c1
				}
				if r2 >= qu.l && r1 <= qu.r && c2 >= qu.d && c1 <= qu.u {
					cnt++
				}
			}
		}
		sb.WriteString(fmt.Sprintf("%d\n", cnt))
	}
	return strings.TrimSpace(sb.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `3 2 1 3 2 2 3 2 3 1 1 3 1
2 3 1 2 2 2 2 2 2 1 2 2 1 1 2 1
2 2 2 1 1 2 1 2 2 1 2 2
3 2 1 2 3 2 3 3 3 2 2 2 2
2 3 1 2 2 2 2 2 2 1 2 1 1 2 2 2
5 2 4 2 5 1 3 5 1 5 1 5 2 5 2
5 2 3 5 2 1 4 1 5 4 5 1 4 2 4
2 3 1 2 1 2 2 2 1 1 2 1 1 1 1 2
4 2 4 2 1 3 1 1 1 1 3 2 3 4
3 1 3 1 2 3 1 3 2
6 1 6 5 4 2 3 1 3 4 6 6
6 1 6 5 2 4 3 1 5 5 6 6
6 1 1 5 3 4 6 2 2 2 2 3
2 2 2 1 1 2 2 2 2 2 2 2
6 2 4 6 3 1 5 2 3 5 5 5 2 1 6 4
2 2 1 2 1 1 1 1 2 2 2 2
4 2 1 3 2 4 1 4 2 4 2 4 2 4
5 3 2 5 1 4 3 4 5 4 5 3 3 5 4 5 5 5 5
5 1 5 1 2 4 3 2 2 4 4
5 1 2 1 4 3 5 3 2 3 4
6 1 4 3 2 1 6 5 2 5 3 6
6 1 3 2 5 4 1 6 4 5 5 5
4 2 2 1 4 3 1 2 2 4 2 4 4 4
4 3 1 2 4 3 4 4 4 4 1 3 1 3 3 3 4 4
5 2 5 4 3 1 2 1 1 5 5 5 5 5 5
6 1 4 2 3 5 6 1 1 4 5 6
3 3 2 1 3 2 1 3 1 3 1 3 1 1 2 2 3
3 3 1 2 3 1 1 1 3 2 1 2 3 1 3 2 3
3 3 3 1 2 3 3 3 3 1 1 1 3 3 2 3 2
3 2 2 1 3 3 1 3 3 1 2 2 3
3 3 1 2 3 3 1 3 2 2 2 3 3 1 3 2 3
2 1 1 2 1 1 1 2
3 2 3 1 2 1 1 3 2 2 2 3 2
4 2 2 1 3 4 2 1 4 3 4 2 4 3
3 2 2 3 1 3 1 3 3 3 3 3 3
2 1 2 1 1 2 1 2
2 3 2 1 1 2 1 2 1 1 1 1 2 2 2 2
6 1 4 1 5 2 6 3 2 2 5 6
3 1 2 3 1 2 1 3 3
4 3 1 4 2 3 3 1 3 4 3 2 3 2 4 1 4 4
4 1 2 1 4 3 1 3 3 3
4 1 1 3 4 2 4 1 4 2
5 1 4 1 5 2 3 5 3 5 3
6 1 2 1 6 3 4 5 1 2 2 5
6 2 4 5 6 3 1 2 6 6 6 6 1 2 4 6
2 1 2 1 2 2 2 2
3 1 2 1 3 3 2 3 2
2 1 2 1 1 1 1 1
2 2 2 1 2 2 2 2 1 1 1 2
2 2 1 2 2 2 2 2 2 1 2 1
6 3 1 3 4 5 6 2 1 4 3 6 5 6 5 6 6 2 6 4
2 2 1 2 1 2 2 2 2 2 2 2
5 3 5 2 3 1 4 5 4 5 4 2 3 2 3 4 2 5 5
5 2 3 2 1 4 5 3 1 3 2 3 1 4 1
5 2 1 5 4 2 3 2 2 2 4 4 2 5 4
5 2 2 1 5 3 4 1 4 5 4 1 5 1 5
5 1 3 5 2 1 4 3 2 5 4
4 3 4 2 3 1 1 1 2 1 3 4 4 4 2 4 3 4
6 1 2 3 5 1 6 4 1 6 3 6
4 3 4 2 3 1 1 3 1 3 2 2 2 3 2 4 2 4
3 2 2 1 3 3 1 3 1 1 1 2 3
5 1 5 2 1 4 3 5 2 5 5
2 3 1 2 1 2 2 2 2 2 2 2 2 1 2 1
5 2 5 4 2 3 1 2 1 2 4 5 2 5 5
6 2 4 3 6 1 2 5 4 3 5 3 2 2 6 2
5 1 5 2 3 4 1 2 4 2 4
6 2 5 2 4 1 3 6 4 6 5 6 1 1 1 3
5 3 2 1 3 5 4 3 4 4 4 4 4 5 4 4 4 5 5
5 3 4 1 2 5 3 3 4 3 5 4 1 5 2 3 4 3 4
6 2 5 3 1 4 6 2 2 1 6 2 6 1 6 6
6 3 1 3 5 2 6 4 4 5 5 5 2 2 5 2 6 4 6 6
4 3 1 4 2 3 1 3 3 4 1 3 1 3 4 2 4 3
2 2 2 1 1 2 1 2 2 1 2 1
2 3 2 1 1 2 1 2 1 1 1 1 1 1 1 1
2 2 1 2 1 2 1 2 1 1 1 2
5 3 5 4 3 1 2 5 1 5 4 3 3 5 5 5 3 5 3
6 3 6 3 5 4 2 1 3 1 5 1 1 3 2 4 6 2 6 2
5 3 1 3 5 2 4 2 2 2 4 3 1 3 5 4 3 4 4
3 3 2 3 1 2 1 3 2 2 1 3 3 2 2 2 2
5 3 5 1 4 2 3 3 1 5 3 4 3 4 3 2 3 4 5
2 1 1 2 1 2 1 2
2 2 2 1 1 2 1 2 2 1 2 1
6 1 1 2 3 6 5 4 3 1 5 1
5 2 3 5 1 4 2 5 2 5 5 4 2 5 3
5 2 3 1 5 2 4 2 3 2 5 2 2 2 4
3 2 1 3 2 2 2 3 3 1 3 2 3
4 2 4 1 2 3 4 4 4 4 4 1 4 4
6 1 4 6 2 5 1 3 5 5 6 6
6 2 6 1 3 2 4 5 5 2 6 4 4 2 4 2
3 3 3 1 2 2 2 3 2 2 2 3 3 3 1 3 2
6 2 3 6 2 1 5 4 5 2 5 2 3 1 4 1
5 2 3 2 4 5 1 5 5 5 5 1 3 3 3
5 2 2 4 5 3 1 5 3 5 5 3 3 5 3
2 3 1 2 2 1 2 1 2 2 2 2 2 1 2 2
3 2 2 1 3 2 1 2 2 2 2 3 3
4 1 3 1 4 2 3 4 3 4
2 3 1 2 2 1 2 2 1 2 2 2 1 2 2 2
4 2 3 2 4 1 2 2 2 3 2 3 2 3
2 2 2 1 2 2 2 2 1 1 1 1
2 2 2 1 2 1 2 1 2 1 2 2`

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		expected := solveCase(line)
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
