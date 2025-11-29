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

const testcases = `6 R 7 L 13 ? 13 R 12 ? 12 ? 7
7 R 10 L 8 ? 10 R 4 L 6 L 9 L 11
2 L 6 ? 6
4 R 7 R 10 ? 10 ? 7
10 R 17 R 3 R 9 ? 17 R 4 R 19 ? 3 R 20 R 1 L 12
9 L 13 L 15 ? 13 R 8 ? 13 ? 15 L 14 L 18 L 5
6 R 13 R 11 ? 11 ? 13 L 9 R 4
10 L 8 ? 8 ? 8 ? 8 ? 8 L 4 R 17 ? 4 ? 17 L 14
9 R 16 R 10 R 13 L 3 R 17 R 19 ? 16 R 12 R 8
7 R 11 R 7 L 12 ? 12 ? 7 ? 11 R 1
1 R 3
2 L 5 ? 5
8 R 2 ? 2 R 17 ? 17 ? 2 R 8 R 6 L 18
5 L 2 R 3 ? 2 R 5 ? 3
1 L 1
6 L 6 R 9 ? 9 ? 6 L 11 ? 6
8 R 13 R 1 ? 13 R 8 L 2 R 5 ? 1 ? 8
9 R 1 R 6 R 8 L 14 ? 1 L 19 ? 14 ? 14 ? 14
5 L 12 L 2 ? 2 L 3 R 4
6 R 14 R 5 ? 5 R 9 R 11 ? 9
10 L 8 ? 8 ? 8 L 7 ? 8 ? 8 ? 8 ? 7 ? 7 L 14
1 R 2
8 L 12 ? 12 ? 12 ? 12 ? 12 L 18 ? 12 ? 18
5 R 6 L 11 L 10 ? 10 R 4
10 L 15 ? 15 ? 15 R 14 ? 14 ? 15 R 6 ? 6 ? 15 L 21
1 R 4
7 R 2 L 12 ? 2 ? 2 ? 12 L 14 ? 14
3 L 2 R 3 R 1
8 R 18 L 17 L 3 L 2 ? 3 ? 3 ? 2 R 6
10 L 15 ? 15 R 17 R 12 R 8 L 21 ? 8 ? 12 ? 17 ? 8
9 R 6 R 13 L 10 L 16 ? 10 ? 13 ? 13 L 7 ? 10
9 L 6 R 18 L 4 L 15 ? 18 ? 15 R 16 ? 15 ? 15
1 L 1
5 L 10 ? 10 ? 10 ? 10 ? 10
10 L 8 R 1 R 2 L 4 ? 4 ? 1 L 10 R 3 ? 3 R 13
10 L 8 R 1 R 5 R 7 R 21 L 14 ? 1 ? 7 R 20 ? 21
4 R 3 ? 3 ? 3 ? 3
4 L 8 L 4 R 7 L 1
9 R 7 L 20 ? 7 ? 20 L 16 ? 7 R 3 R 19 ? 16
7 R 11 ? 11 ? 11 R 12 ? 11 R 15 L 16
6 L 13 ? 13 L 4 ? 4 L 14 L 5
3 R 8 ? 8 ? 8
2 R 4 ? 4
9 R 13 R 6 ? 6 ? 13 ? 13 ? 6 R 8 R 4 L 19
5 R 11 L 4 ? 4 ? 4 R 12
4 L 5 ? 5 ? 5 L 10
3 R 1 R 3 R 8
5 R 3 ? 3 L 1 ? 3 R 4
8 L 6 L 12 R 14 ? 12 L 16 ? 14 L 7 R 15
6 L 8 L 4 L 13 ? 4 ? 8 ? 4
1 R 2
3 R 3 L 4 L 8
7 R 7 R 5 R 12 ? 12 ? 5 L 16 L 4
8 R 11 ? 11 R 1 ? 1 L 12 ? 11 L 8 ? 8
10 R 9 ? 9 R 7 ? 7 ? 7 R 3 ? 7 ? 3 ? 7 R 14
5 R 7 L 10 ? 7 R 3 ? 7
3 R 1 R 2 L 6
3 R 8 L 2 ? 2
6 R 13 L 8 R 2 L 1 L 5 L 4
8 R 3 ? 3 ? 3 ? 3 ? 3 ? 3 ? 3 ? 3
2 R 6 ? 6
10 L 3 R 5 ? 3 R 20 ? 5 L 13 ? 13 R 12 L 10 ? 10
4 L 6 R 7 ? 6 ? 6
5 R 12 L 2 L 6 ? 6 ? 2
5 L 12 ? 12 ? 12 ? 12 R 11
10 L 20 ? 20 ? 20 L 11 ? 20 L 12 L 22 R 4 R 1 R 14
8 R 15 ? 15 L 4 ? 4 ? 15 R 9 ? 4 ? 9
8 L 18 ? 18 R 16 ? 18 R 6 L 13 ? 6 ? 6
7 L 4 ? 4 ? 4 L 5 ? 5 L 3 ? 3
2 R 3 L 6
2 R 4 L 3
7 L 7 ? 7 L 5 ? 5 ? 7 R 2 R 13
2 R 1 ? 1
4 R 5 ? 5 L 3 R 9
3 R 3 L 1 R 7
7 R 6 ? 6 R 9 ? 9 L 8 L 12 L 5
2 L 2 L 4
4 R 4 L 9 L 6 L 8
2 R 1 R 4
8 L 5 ? 5 ? 5 ? 5 ? 5 ? 5 L 17 ? 5
10 R 1 ? 1 R 15 ? 15 ? 1 R 22 ? 15 L 14 ? 15 L 5
10 L 7 ? 7 L 3 ? 3 R 9 R 17 ? 3 ? 17 ? 7 L 19
4 L 3 ? 3 L 10 L 5
6 L 7 ? 7 ? 7 ? 7 ? 7 ? 7
2 R 4 ? 4
2 R 2 ? 2
10 L 11 L 5 L 14 R 1 R 10 ? 1 ? 11 R 6 L 17 ? 6
8 R 1 ? 1 R 17 ? 1 ? 1 ? 1 ? 17 L 13
10 L 8 ? 8 L 18 ? 8 L 3 ? 3 R 7 ? 18 ? 8 ? 3
2 L 3 L 6
10 R 13 ? 13 ? 13 L 1 R 14 R 8 R 20 ? 20 R 15 ? 15
1 L 4
6 L 8 R 12 R 3 ? 8 ? 12 R 10
4 L 7 ? 7 L 10 ? 7
3 R 5 R 1 L 3
8 R 9 R 14 L 2 R 13 R 4 ? 2 ? 9 ? 9
9 R 14 L 4 L 9 R 7 ? 14 ? 4 R 10 ? 14 ? 7
7 L 16 ? 16 L 6 R 14 ? 16 ? 6 ? 6
9 L 12 L 8 ? 12 L 19 ? 8 ? 19 ? 8 L 6 R 4
5 L 9 R 6 ? 6 L 4 R 11`

type op struct {
	cmd string
	id  int
}

func solve(ops []op) []int {
	pos := make(map[int]int)
	l, r := 0, 0
	size := 0
	res := []int{}
	for _, o := range ops {
		switch o.cmd {
		case "L":
			if size == 0 {
				pos[o.id] = 0
				l = 0
				r = 0
				size = 1
			} else {
				l--
				pos[o.id] = l
				size++
			}
		case "R":
			if size == 0 {
				pos[o.id] = 0
				l = 0
				r = 0
				size = 1
			} else {
				r++
				pos[o.id] = r
				size++
			}
		case "?":
			p := pos[o.id]
			left := p - l
			right := r - p
			if left < right {
				res = append(res, left)
			} else {
				res = append(res, right)
			}
		}
	}
	return res
}

func runCase(bin string, ops []op) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, o := range ops {
		input.WriteString(fmt.Sprintf("%s %d\n", o.cmd, o.id))
	}
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotLines := strings.Fields(strings.TrimSpace(out.String()))
	expectVals := solve(ops)
	if len(gotLines) != len(expectVals) {
		return fmt.Errorf("expected %d lines got %d", len(expectVals), len(gotLines))
	}
	for i, v := range expectVals {
		if gotLines[i] != strconv.Itoa(v) {
			return fmt.Errorf("mismatch at %d: expected %d got %s", i+1, v, gotLines[i])
		}
	}
	return nil
}

func parseOps(parts []string) ([]op, error) {
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty line")
	}
	q, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	if len(parts) != 1+2*q {
		return nil, fmt.Errorf("need %d tokens", 1+2*q)
	}
	ops := make([]op, q)
	idx := 1
	for i := 0; i < q; i++ {
		cmd := parts[idx]
		id, _ := strconv.Atoi(parts[idx+1])
		idx += 2
		ops[i] = op{cmd, id}
	}
	return ops, nil
}

func loadTestcases() ([][]op, error) {
	scanner := bufio.NewScanner(strings.NewReader(testcases))
	var cases [][]op
	lineNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		lineNum++
		parts := strings.Fields(line)
		ops, err := parseOps(parts)
		if err != nil {
			return nil, fmt.Errorf("invalid test line %d: %w", lineNum, err)
		}
		cases = append(cases, ops)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := loadTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for idx, ops := range cases {
		if err := runCase(bin, ops); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
