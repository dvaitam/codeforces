package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"7 14 1 3 5 4 4 7 7 3 4 3 5 2 5 2",
	"5 5 1 5 3 5 5",
	"3 10 1 3 1 3 2 2 3 1 2 2",
	"6 20 6 2 5 4 4 5 3 1 5 1 1 6 4 6 6 6 1 5 4 3",
	"4 11 1 2 2 2 2 4 1 1 3 4 1",
	"5 18 3 1 5 3 5 2 5 5 5 3 4 1 5 4 3 5 2 3",
	"3 7 1 1 3 3 2 2 1",
	"2 5 1 1 1 2 2",
	"9 8 4 7 5 8 8 6 2 6",
	"10 4 8 10 6 4",
	"4 1 3",
	"2 8 2 1 2 2 1 1 1 1",
	"1 19 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1",
	"10 4 5 2 4 2",
	"5 12 4 2 1 5 4 1 5 1 4 2 3 3",
	"8 19 3 4 1 3 3 6 5 2 8 3 1 8 7 5 6 7 5 3 1",
	"8 3 6 1 5",
	"3 8 2 2 3 2 3 2 3 3",
	"10 5 5 7 7 2 1",
	"10 7 6 3 4 4 8 7 10",
	"7 2 4 7",
	"10 14 1 3 8 2 5 3 8 9 8 9 10 1 1 8",
	"6 10 4 1 4 2 5 6 1 6 2 1",
	"7 14 3 1 2 1 6 7 1 7 6 5 5 1 2 1",
	"10 7 5 5 3 2 8 7 2",
	"1 9 1 1 1 1 1 1 1 1 1",
	"1 2 1 1",
	"9 11 6 1 8 8 7 6 9 3 4 7 5",
	"1 5 1 1 1 1 1",
	"2 11 1 1 2 1 1 2 2 2 1 2 1",
	"8 8 1 5 3 2 5 7 6 5",
	"7 4 1 5 4 4",
	"6 11 1 4 1 6 4 4 1 3 3 6 6",
	"3 6 3 3 2 3 1 1",
	"2 7 1 1 2 1 1 2 2",
	"8 16 4 7 2 6 4 5 3 7 4 6 2 2 1 8 4 2",
	"8 13 5 4 1 4 3 2 4 8 7 6 3 2 8",
	"3 19 2 3 3 2 3 2 3 2 2 2 3 3 1 3 3 1 1 2 3",
	"6 11 1 5 2 3 5 2 4 5 3 6 6",
	"8 3 2 1 2",
	"4 5 1 3 1 4 3",
	"3 5 3 2 2 3 2",
	"9 17 1 2 9 2 7 4 5 9 7 8 7 4 1 1 3 5 9",
	"10 9 6 2 8 5 5 7 7 7 1",
	"3 5 1 2 3 2 1",
	"1 16 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1",
	"1 20 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1",
	"8 12 4 8 6 2 1 1 8 5 1 4 4 2",
	"9 17 7 9 5 2 3 7 7 2 2 7 2 2 7 3 1 8 7",
	"7 1 4",
	"6 9 1 3 1 1 3 6 1 3 3",
	"3 1 1",
	"6 3 5 2 2",
	"1 7 1 1 1 1 1 1 1",
	"3 15 1 2 2 3 2 1 1 1 2 2 2 2 2 3 3",
	"6 6 5 1 1 5 5 3",
	"3 13 1 1 1 2 3 1 1 1 2 2 2 3 1",
	"3 20 1 2 1 3 1 1 2 2 3 2 3 1 3 2 2 3 2 1 1 2",
	"6 17 1 4 4 1 4 6 3 4 2 3 3 4 1 2 1 3 1",
	"9 20 3 8 7 3 7 7 3 4 8 6 9 3 6 8 2 8 4 5 1 8",
	"10 15 1 4 5 2 5 9 10 3 7 8 2 8 4 9 7",
	"5 1 1",
	"5 2 1 3",
	"7 17 5 6 4 4 1 6 3 3 3 7 7 6 2 5 1 1 1",
	"5 10 5 3 1 5 2 2 1 4 3 3",
	"9 5 9 4 9 2 7",
	"9 13 5 5 8 6 3 3 2 2 7 7 8 3 9",
	"5 12 4 4 2 4 4 5 3 4 1 4 3 2",
	"8 2 4 1",
	"6 16 4 1 5 1 6 1 6 6 6 4 1 3 1 1 5 1",
	"5 10 2 2 5 3 2 1 4 4 3 4",
	"3 11 2 3 3 2 1 2 3 1 3 2 1",
	"4 6 4 3 4 4 4 4",
	"4 7 4 2 1 4 1 2 1",
	"3 12 1 3 3 3 1 1 3 2 3 1 3 3",
	"5 12 4 4 1 5 5 4 5 4 4 3 4 2",
	"6 9 1 1 1 2 3 1 3 6 1",
	"3 3 2 3 1",
	"10 13 9 4 8 4 6 10 2 10 2 6 6 9 8",
	"6 9 1 5 1 2 3 1 2 5 3",
	"4 7 3 3 3 4 3 4 3",
	"4 2 3 1",
	"1 15 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1",
	"8 20 2 7 7 1 3 4 8 4 3 5 6 6 7 2 5 4 5 8 8 5",
	"5 8 1 1 5 1 2 4 2 2",
	"5 1 5",
	"9 14 1 2 7 5 2 6 4 9 5 4 4 2 9 5",
	"6 8 3 6 4 3 5 2 2 1",
	"9 17 6 6 1 3 7 3 3 9 2 3 4 8 4 4 3 4 7",
	"6 20 5 2 6 4 1 5 1 5 5 3 4 4 3 1 2 5 6 2 6 4",
	"8 18 6 2 5 3 7 4 6 5 7 1 4 1 6 4 6 8 4 5",
	"6 6 3 1 3 5 5 1",
	"3 12 1 2 3 1 1 1 1 1 1 3 2 1",
	"1 12 1 1 1 1 1 1 1 1 1 1 1 1",
	"1 14 1 1 1 1 1 1 1 1 1 1 1 1 1 1",
	"4 1 2",
	"9 6 2 7 2 8 3 1",
	"5 11 4 1 1 4 1 3 3 2 4 2 5",
	"6 6 6 4 3 3 4 4",
	"1 10 1 1 1 1 1 1 1 1 1 1",
}

const testcasesCount = 100

func referenceSolve(n int, moves []int) []string {
	pos := make([]int, n+1)
	posts := make([]int, n+1)
	minPos := make([]int, n+1)
	maxPos := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pos[i] = i
		posts[i] = i
		minPos[i] = i
		maxPos[i] = i
	}
	for _, x := range moves {
		if pos[x] > 1 {
			y := posts[pos[x]-1]
			posts[pos[x]-1], posts[pos[x]] = posts[pos[x]], posts[pos[x]-1]
			pos[x]--
			pos[y]++
			if minPos[x] > pos[x] {
				minPos[x] = pos[x]
			}
			if maxPos[y] < pos[y] {
				maxPos[y] = pos[y]
			}
		}
		if maxPos[x] < pos[x] {
			maxPos[x] = pos[x]
		}
	}
	res := make([]string, n)
	for i := 1; i <= n; i++ {
		res[i-1] = fmt.Sprintf("%d %d", minPos[i], maxPos[i])
	}
	return res
}

func runCase(bin string, line string, idx int) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("case %d: invalid line", idx)
	}
	n, _ := strconv.Atoi(fields[0])
	m, _ := strconv.Atoi(fields[1])
	if len(fields) != 2+m {
		return fmt.Errorf("case %d: expected %d moves got %d", idx, m, len(fields)-2)
	}
	moves := make([]int, m)
	for i := 0; i < m; i++ {
		moves[i], _ = strconv.Atoi(fields[2+i])
	}
	expectLines := referenceSolve(n, moves)
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i, v := range moves {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(strconv.Itoa(v))
	}
	input.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(gotLines) != len(expectLines) {
		return fmt.Errorf("expected %d lines got %d", len(expectLines), len(gotLines))
	}
	for i := range expectLines {
		if strings.TrimSpace(gotLines[i]) != expectLines[i] {
			return fmt.Errorf("line %d: expected %s got %s", i+1, expectLines[i], strings.TrimSpace(gotLines[i]))
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	if len(testcases) != testcasesCount {
		fmt.Fprintf(os.Stderr, "unexpected testcase count: got %d want %d\n", len(testcases), testcasesCount)
		os.Exit(1)
	}

	for i, tc := range testcases {
		if err := runCase(bin, tc, i+1); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
