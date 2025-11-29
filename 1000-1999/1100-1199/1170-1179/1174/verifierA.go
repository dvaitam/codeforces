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

const testcasesRaw = `4 14 2 9 17 16 13 10 16
3 19 7 17 5 10 5
1 20 9
5 20 5 10 4 3 11 16 18 4 12
4 11 20 7 18 16 15 17 9
5 1 20 11 1 6 20 9 10 3 20
3 13 17 10 17 20 15
2 5 20 18 15
5 7 2 2 19 5 16 1 16 19 8
5 11 14 3 15 10 7 10 7 8 1
2 9 4 15 9
2 15 11 3 16
2 17 13 19 11
4 11 10 2 6 3 7 14 3
3 1 17 1 15 3 2
1 17 10
1 12 14
5 19 8 6 19 6 19 4 3 19 14
4 2 13 1 1 6 9 5 9
3 1 11 1 10 3 6
1 20 18
2 10 4 17 11
2 14 2 1 14
2 10 12 3 2
1 18 18
5 16 11 4 10 1 11 5 16 1 18
5 14 14 9 7 18 19 14 19 19 10
1 4 10
1 4 3
4 5 19 4 9 18 5 20 4
3 9 11 12 13 3 12
5 6 18 7 2 13 14 18 8 17 14
1 5 2
4 11 13 8 20 6 17 4 13
3 2 8 17 17 17 13
4 20 14 4 1 16 13 11 18
1 18 5
2 15 20 2 8
5 16 16 10 14 8 20 11 5 6 5
3 20 3 11 2 7 2
3 17 16 15 2 1 2
2 5 10 10 4
5 8 5 11 5 10 19 2 3 19 12
4 4 1 11 8 13 19 7 6
4 9 19 12 7 8 9 10 3
5 17 6 1 12 12 3 4 12 19 3
5 9 7 17 12 15 6 13 1 7 4
5 1 18 15 7 17 7 8 16 4 11
5 5 8 2 14 8 14 17 3 11 4
1 7 12
1 2 8
4 13 10 11 20 11 14 16 8
4 16 5 16 9 4 13 9 10
5 3 3 3 11 8 12 14 7 4 8
2 15 9 7 13
4 18 19 17 9 10 12 3 6
4 16 1 20 10 5 5 2 3
5 20 12 7 16 10 6 16 10 3 18
1 3 19
1 3 18
4 19 2 11 8 5 16 10 1
4 17 8 3 12 3 9 16 16
2 8 10 3 13
3 20 12 12 19 13 18
3 13 10 1 18 1 20
1 7 2
1 12 18
1 15 18
2 1 3 18 17
4 2 16 14 15 1 12 8 8
3 16 17 1 18 10 19
1 4 15
2 14 10 14 12
2 15 19 20 7
3 9 14 15 14 14 19
2 9 19 6 14
4 12 5 15 7 20 5 18 7
2 18 4 20 3
2 1 14 12 16
1 13 18
3 5 12 17 8 9 13
4 9 7 10 17 11 1 11 6
1 12 3
2 14 2 2 9
2 5 6 14 10
5 7 8 20 16 12 18 10 20 19 20
5 18 9 10 3 16 4 7 7 8 17
1 11 7
4 1 19 12 5 11 3 16 7
5 11 13 1 5 9 11 7 8 19 12
1 16 9
1 2 15
2 14 17 1 16
1 19 12
5 2 9 17 5 3 2 19 17 13 7
3 2 1 7 6 9 11
2 4 15 17 5
5 6 15 6 11 20 16 4 18 16 4
1 5 16
1 20 10
2 20 8 19 13
1 11 9`

func referenceSolve(vals []int64) string {
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	allEqual := true
	for i := 0; i < len(vals)-1; i++ {
		if vals[i] != vals[i+1] {
			allEqual = false
			break
		}
	}
	if allEqual {
		return "-1"
	}
	var out strings.Builder
	for i, v := range vals {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(strconv.FormatInt(v, 10))
	}
	return out.String()
}

func runCase(bin string, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return fmt.Errorf("invalid test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid n: %v", err)
	}
	if len(fields) != 1+2*n {
		return fmt.Errorf("expected %d numbers got %d", 1+2*n, len(fields))
	}
	vals := make([]int64, 2*n)
	for i := 0; i < 2*n; i++ {
		v, err := strconv.ParseInt(fields[1+i], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid value: %v", err)
		}
		vals[i] = v
	}
	expect := referenceSolve(vals)
	input := line + "\n"

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
