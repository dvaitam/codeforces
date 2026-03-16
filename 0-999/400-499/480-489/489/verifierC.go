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
const testcasesCRaw = `
2 2
3 11
6 47
10 32
20 54
20 9
19 40
14 81
13 102
17 95
18 113
17 68
2 0
12 59
11 48
14 114
17 42
18 45
8 29
1 2
11 22
5 32
17 92
17 143
6 28
14 94
17 93
19 90
12 57
6 48
13 91
15 135
8 62
9 63
17 131
12 84
15 118
12 72
18 116
16 56
11 89
6 39
9 61
10 38
17 143
17 129
20 150
14 39
7 62
17 93
20 19
11 92
1 3
4 3
19 167
2 8
19 58
4 33
5 17
8 26
2 13
2 1
12 46
6 15
1 1
4 4
1 0
1 5
9 16
6 47
6 33
1 6
19 11
8 19
2 0
12 78
4 18
11 62
1 4
15 11
9 51
20 180
5 30
8 11
11 13
1 7
5 33
19 100
16 131
11 18
11 33
9 77
14 83
1 8
5 42
2 8
2 4
6 10
4 29
8 65
`


func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(m, s int) (string, string) {
	if (s == 0 && m > 1) || s > 9*m {
		return "-1", "-1"
	}
	if s == 0 && m == 1 {
		return "0", "0"
	}
	rem := s
	max := make([]byte, m)
	for i := 0; i < m; i++ {
		d := 9
		if rem < 9 {
			d = rem
		}
		max[i] = byte('0' + d)
		rem -= d
	}
	rem = s
	min := make([]byte, m)
	for i := 0; i < m; i++ {
		low := 0
		if i == 0 {
			low = 1
		}
		for d := low; d <= 9; d++ {
			if rem-d < 0 {
				break
			}
			if rem-d <= 9*(m-i-1) {
				min[i] = byte('0' + d)
				rem -= d
				break
			}
		}
	}
	return string(min), string(max)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesCRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Fprintf(os.Stderr, "bad test line %d\n", idx)
			os.Exit(1)
		}
		mVal, _ := strconv.Atoi(fields[0])
		sVal, _ := strconv.Atoi(fields[1])
		in := line + "\n"
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		wantMin, wantMax := solve(mVal, sVal)
		tokens := strings.Fields(out)
		if len(tokens) != 2 {
			fmt.Fprintf(os.Stderr, "case %d: invalid output\n", idx)
			os.Exit(1)
		}
		if tokens[0] != wantMin || tokens[1] != wantMax {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s %s got %s %s\n", idx, wantMin, wantMax, tokens[0], tokens[1])
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
