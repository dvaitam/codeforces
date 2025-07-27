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

type Interval struct{ l, r int }

func solveCase(n, x, m int, segs []Interval) int {
	l, r := x, x
	for _, s := range segs {
		if s.l <= r && s.r >= l {
			if s.l < l {
				l = s.l
			}
			if s.r > r {
				r = s.r
			}
		}
	}
	return r - l + 1
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 3 {
			fmt.Printf("invalid test case %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		x, _ := strconv.Atoi(fields[1])
		m, _ := strconv.Atoi(fields[2])
		if len(fields) != 3+2*m {
			fmt.Printf("test %d wrong field count\n", idx)
			os.Exit(1)
		}
		segs := make([]Interval, m)
		for i := 0; i < m; i++ {
			li, _ := strconv.Atoi(fields[3+2*i])
			ri, _ := strconv.Atoi(fields[4+2*i])
			segs[i] = Interval{li, ri}
		}
		expected := solveCase(n, x, m, segs)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, x, m))
		for i, s := range segs {
			input.WriteString(fmt.Sprintf("%d %d", s.l, s.r))
			if i < m-1 {
				input.WriteByte('\n')
			}
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
