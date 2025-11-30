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

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type segment struct {
	start, end int
}

// solution logic from 1927D.go
func solveCase(fields []string) (string, error) {
	if len(fields) < 2 {
		return "", fmt.Errorf("invalid test line")
	}
	pos := 0
	n, err := strconv.Atoi(fields[pos])
	if err != nil {
		return "", err
	}
	pos++
	if len(fields) < pos+n+1 {
		return "", fmt.Errorf("not enough data for array")
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i], _ = strconv.Atoi(fields[pos+i])
	}
	pos += n
	q, err := strconv.Atoi(fields[pos])
	if err != nil {
		return "", err
	}
	pos++
	limit := pos + 2*q
	if len(fields) < limit {
		return "", fmt.Errorf("expected %d query numbers got %d", 2*q, len(fields)-pos)
	}
	if len(fields) > limit {
		fields = fields[:limit]
	}
	if len(fields) != pos+2*q {
		return "", fmt.Errorf("expected %d query numbers got %d", 2*q, len(fields)-pos)
	}

	segs := make([]segment, 0, n)
	for i := 0; i < n; i++ {
		j := i
		for j < n && a[j] == a[i] {
			j++
		}
		segs = append(segs, segment{i + 1, j})
		i = j - 1
	}

	var outLines []string
	for i := 0; i < q; i++ {
		l, _ := strconv.Atoi(fields[pos+2*i])
		r, _ := strconv.Atoi(fields[pos+2*i+1])
		idx := sort.Search(len(segs), func(i int) bool {
			return segs[i].start > l
		})
		idx--
		if idx >= 0 && segs[idx].end < r {
			outLines = append(outLines, fmt.Sprintf("%d %d", l, segs[idx].end+1))
		} else {
			outLines = append(outLines, "-1 -1")
		}
	}
	return strings.Join(outLines, "\n"), nil
}

const testcasesData = `
4 5 5 2 3 10 4 4 1 4 3 3 2 4 4 4 4 4 2 4 2 4 4 4 1 2
13 5 1 3 1 3 4 5 4 4 4 5 4 2 6 2 2 3 10 4 8 11 12 13 13 7 11
14 4 5 3 5 5 4 5 2 3 1 3 5 2 3 9 10 14 2 13 11 12 11 13 5 6 2 9 14 14 2 7 13 13
7 2 1 3 4 4 1 1 10 5 5 4 6 5 6 5 5 1 3 1 1 1 5 5 5 2 5 3 7
5 2 1 3 3 3 3 4 5 4 5 5 5
10 5 3 4 2 3 4 3 5 3 5 6 1 7 10 10 1 7 10 10 1 6 8 9
11 3 5 3 4 1 5 1 1 3 3 4 5 10 11 3 8 3 8 6 10 5 7
13 4 1 1 5 2 3 5 2 3 2 3 2 4 2 2 11 6 11
11 2 4 2 1 3 2 5 4 3 2 1 1 9 9
6 5 2 3 3 1 5 6 5 5 4 5 5 6 4 5 6 6 3 6
10 4 1 4 2 2 1 4 5 5 4 9 4 4 8 10 9 10 9 10 4 10 2 6 2 5 1 1 9 9
15 4 5 1 1 4 1 2 5 3 2 1 5 5 4 1 10 2 7 3 7 14 15 13 13 6 9 4 5 9 15 2 4 4 8 15 15
14 1 4 5 4 1 3 2 3 5 5 5 4 1 4 6 13 13 14 14 13 13 1 2 1 2 8 8
14 1 5 5 4 3 2 3 1 3 4 4 5 3 3 5 4 9 7 8 3 11 1 12 12 13
13 1 5 2 1 3 4 5 5 4 1 5 4 1 6 11 12 13 13 7 12 7 10 1 4 4 12
5 5 1 4 2 4 3 1 3 3 5 3 3
8 1 5 4 1 3 5 5 1 10 1 8 3 4 7 7 2 6 2 7 7 7 1 6 2 2 2 7 8 8
10 3 1 1 5 5 5 2 1 5 1 9 1 9 6 10 3 4 4 5 4 7 10 10 5 7 10 10 6 10
7 1 4 5 2 4 2 4 10 7 7 2 7 4 5 2 2 4 7 6 7 5 7 7 7 2 4 7 7
3 5 5 3 4 3 3 3 3 3 3 1 2
1 3 8 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
13 1 4 2 2 1 2 3 2 2 3 2 2 5 1 5 7
14 1 3 2 4 1 1 1 1 3 3 1 3 4 5 6 1 1 6 11 7 13 8 8 4 14 10 13
7 2 5 3 1 3 1 4 2 4 6 1 5
12 3 3 4 3 3 1 3 5 5 5 1 4 9 6 6 12 12 11 11 11 11 3 8 11 12 2 3 9 10 6 11
12 5 4 5 3 2 4 4 3 2 1 1 2 9 9 12 6 6 5 9 7 7 3 3 8 12 5 8 12 12 6 9
8 5 1 3 4 1 2 3 5 2 2 6 2 3
12 2 5 4 4 2 5 5 2 4 2 5 5 3 10 10 4 8 6 12
5 1 4 4 4 3 9 5 5 4 5 4 4 5 5 1 1 2 5 2 5 2 3 5 5
1 5 8 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
14 3 2 1 5 1 1 3 2 1 2 4 2 4 3 2 13 13 7 8
14 2 2 4 4 4 1 5 4 2 4 3 1 4 4 7 8 9 8 8 12 13 6 11 8 12 6 12 4 4
13 2 3 3 2 4 5 2 2 2 1 2 5 4 9 3 13 1 3 2 11 3 10 8 9 1 1 7 10 6 12 1 12
12 1 2 4 1 4 4 1 2 2 1 4 4 4 3 8 10 10 6 6 10 10
13 3 3 4 3 4 2 5 3 1 3 3 3 1 1 11 12
2 5 5 1 1 1
11 1 1 2 5 1 4 1 2 5 3 2 8 6 9 6 11 1 7 5 11 10 11 2 6 3 9 2 10
7 5 3 5 4 2 4 5 6 2 4 7 7 4 5 4 7 3 5 7 7
9 5 4 4 1 2 2 1 4 1 2 6 7 1 1
8 4 3 3 1 1 2 4 1 6 5 5 8 8 4 5 1 3 3 7 1 2
4 3 2 2 5 9 4 4 2 3 2 4 1 1 1 1 1 2 3 3 1 4 4 4
14 5 1 4 3 4 5 4 1 3 5 4 4 2 1 9 1 8 5 6 6 11 4 11 13 13 9 14 12 13 7 8 10 14
4 2 3 1 3 4 4 4 1 3 2 3 2 4
15 2 2 5 4 5 4 3 2 5 4 4 4 4 5 1 5 15 15 3 4 1 12 3 7 9 13
1 4 1 1 1
5 4 4 1 3 4 4 3 3 1 4 3 4 4 4
4 2 4 3 4 6 2 3 3 4 4 4 3 4 3 3 4 4
15 3 3 2 5 1 1 2 5 2 5 1 2 1 1 2 8 6 11 9 9 8 10 4 4 5 11 15 15 14 14 10 14
13 1 4 3 3 2 4 4 3 3 1 1 4 1 3 13 13 3 9 9 11
9 2 3 1 2 1 1 4 1 4 8 9 9 6 7 5 6 2 4 9 9 7 8 8 9 5 7
8 3 5 2 5 3 2 5 1 7 8 8 8 8 1 6 2 5 3 5 5 5 7 7
6 1 5 4 4 2 3 6 2 5 4 5 1 3 2 2 3 5 5 6
9 3 1 2 4 3 4 1 4 2 4 2 3 2 2 5 5 7 8
8 4 3 3 5 5 1 2 2 8 3 4 3 8 3 8 8 8 1 3 7 7 3 8 3 5
4 2 2 1 5 3 2 3 1 4 4 4
1 3 2 1 1 1 1
3 4 3 5 2 1 1 2 2
8 2 1 3 5 4 1 3 4 1 6 8
12 4 5 4 4 5 5 3 4 2 5 4 4 4 10 11 12 12 5 7 6 8
14 2 2 1 5 1 4 5 2 1 5 1 2 1 4 10 11 13 9 14 13 14 2 9 9 13 8 10 3 12 10 11 6 13 9 13
6 3 5 3 5 3 3 5 3 5 2 2 5 5 6 6 6 6
5 4 3 4 2 4 6 5 5 5 5 5 5 1 4 4 5 1 2
14 4 4 2 1 3 5 2 1 5 4 1 4 4 5 6 1 5 6 14 6 12 8 10 11 11 10 13
3 3 2 1 6 1 3 1 1 1 1 2 3 2 2 1 1
9 3 5 4 3 2 4 3 5 1 9 7 7 8 8 5 7 9 9 2 4 2 8 1 6 7 7 8 8
8 3 2 4 2 5 4 5 3 7 5 5 2 2 8 8 2 8 1 6 8 8 2 8
11 3 4 5 2 3 2 4 1 1 5 1 9 9 9 10 10 3 6 8 11 5 7 8 9 4 8 6 10 3 11
6 5 5 3 3 5 5 4 4 4 6 6 3 4 4 4
3 1 4 1 9 2 2 2 2 1 2 2 2 2 2 1 2 2 3 1 1 2 3
6 3 1 3 1 2 2 5 4 4 5 5 2 5 6 6 2 4
10 2 2 5 4 5 5 4 5 2 3 8 6 6 2 10 8 9 1 10 3 7 9 10 8 8 7 10
6 4 5 4 2 1 4 6 6 6 5 6 1 3 5 5 3 6 6 6
14 4 3 4 1 5 3 2 5 1 5 5 3 5 3 9 6 12 14 14 1 7 11 13 7 14 6 8 4 8 12 12 10 13
12 1 3 5 2 3 4 2 5 4 4 2 4 9 10 11 2 2 10 11 12 12 3 11 11 11 2 3 9 9 7 11
9 3 4 3 1 4 3 3 2 3 2 4 4 8 9
7 2 3 3 2 3 5 5 4 6 6 2 4 7 7 6 7
15 1 1 2 1 2 4 3 1 3 2 2 1 4 3 5 3 6 13 12 15 9 13
11 5 4 4 3 3 5 3 5 4 2 1 1 3 9
11 1 2 2 5 2 3 3 1 5 5 3 4 2 4 4 5 8 11 5 6
2 5 3 10 1 2 1 1 1 1 1 1 1 2 2 2 2 2 1 1 2 2 1 1
1 3 9 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
7 4 1 4 3 5 1 5 10 2 4 1 1 4 4 5 6 2 3 4 6 6 7 6 6 5 7 4 6
13 1 2 4 5 5 1 3 5 5 2 4 5 2 5 9 10 4 9 9 11 2 3 8 8
8 2 4 1 1 1 4 5 1 1 8 8
1 3 5 1 1 1 1 1 1 1 1 1 1
2 2 2 2 1 2 1 1
2 3 5 9 1 2 1 2 2 2 2 2 2 2 2 2 1 2 2 2 2 2
12 3 2 3 1 3 4 4 4 2 3 5 1 4 10 10 10 11 4 10 3 8
2 2 5 8 2 2 2 2 2 2 1 2 2 2 2 2 1 1 1 1
13 5 1 4 1 5 4 4 4 4 5 5 5 4 10 4 11 7 7 11 13 12 12 1 4 10 10 3 9 13 13 2 3 10 11
15 3 5 2 4 4 2 2 2 4 4 3 4 5 1 3 9 1 10 6 11 9 12 4 5 5 8 2 5 15 15 4 7 14 15
4 1 3 5 3 4 2 4 4 4 3 3 2 3
14 5 1 3 3 2 3 2 3 2 2 5 1 4 4 5 1 9 7 11 10 12 14 14 8 9
5 3 1 1 4 4 6 5 5 3 5 1 5 5 5 1 2 3 4
10 2 2 3 1 1 3 5 5 1 4 6 9 10 1 8 8 10 8 8 8 8 2 4
6 4 3 4 5 2 4 7 4 5 6 6 6 6 5 6 4 6 4 6 4 4
8 1 3 2 2 5 1 2 2 5 6 6 6 8 1 4 6 6 1 7
3 1 5 4 6 3 3 1 2 3 3 3 3 1 1 3 3
4 1 3 1 3 1 2 2
12 4 1 1 5 2 3 4 5 3 3 1 2 9 11 11 4 10 3 8 8 12 9 12 11 12 3 12 12 12 11 12
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		pos := 0
		n, _ := strconv.Atoi(fields[pos])
		pos++
		if len(fields) < pos+n {
			fmt.Printf("test %d short array\n", idx)
			os.Exit(1)
		}
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			v, _ := strconv.Atoi(fields[pos+i])
			a[i] = v
			sb.WriteString(fields[pos+i])
		}
		sb.WriteByte('\n')
		pos += n
		if len(fields) <= pos {
			fmt.Printf("test %d missing q\n", idx)
			os.Exit(1)
		}
		q, _ := strconv.Atoi(fields[pos])
		pos++
		sb.WriteString(fmt.Sprintf("%d\n", q))
		limit := pos + 2*q
		if len(fields) < limit {
			fmt.Printf("test %d expected %d query numbers got %d\n", idx, 2*q, len(fields)-pos)
			os.Exit(1)
		}
		if len(fields) > limit {
			fields = fields[:limit]
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, _ := strconv.Atoi(fields[pos+2*i])
			r, _ := strconv.Atoi(fields[pos+2*i+1])
			queries[i] = [2]int{l, r}
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		input := sb.String()
		want, err := solveCase(fields)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx, err)
			os.Exit(1)
		}
		outScanner := bufio.NewScanner(strings.NewReader(got))
		answers := make([]string, 0, q)
		for outScanner.Scan() {
			line := strings.TrimSpace(outScanner.Text())
			if line != "" {
				answers = append(answers, line)
			}
		}
		if err := outScanner.Err(); err != nil {
			fmt.Printf("test %d output scan error: %v\n", idx, err)
			os.Exit(1)
		}
		if len(answers) != q {
			fmt.Printf("test %d expected %d lines, got %d\n", idx, q, len(answers))
			os.Exit(1)
		}
		if strings.TrimSpace(strings.Join(answers, "\n")) != want {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx, input, want, strings.Join(answers, "\n"))
			os.Exit(1)
		}
		for i := 0; i < q; i++ {
			parts := strings.Fields(answers[i])
			if len(parts) != 2 {
				fmt.Printf("test %d line %d invalid output\n", idx, i+1)
				os.Exit(1)
			}
			x, err1 := strconv.Atoi(parts[0])
			y, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				fmt.Printf("test %d line %d non-integer output\n", idx, i+1)
				os.Exit(1)
			}
			l, r := queries[i][0], queries[i][1]
			if x == -1 && y == -1 {
				same := true
				for k := l; k < r; k++ {
					if a[k-1] != a[k] {
						same = false
						break
					}
				}
				if !same {
					fmt.Printf("test %d line %d incorrect -1 -1\n", idx, i+1)
					os.Exit(1)
				}
			} else {
				if x < l || x > r || y < l || y > r {
					fmt.Printf("test %d line %d indices out of range\n", idx, i+1)
					os.Exit(1)
				}
				if a[x-1] == a[y-1] {
					fmt.Printf("test %d line %d values not distinct\n", idx, i+1)
					os.Exit(1)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
