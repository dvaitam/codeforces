package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n    int
	logs []string
}

// solve is the embedded logic from 16D.go.
func solve(tc testCase) int {
	days := 1
	prevTime := -1
	for i := 0; i < tc.n; i++ {
		line := tc.logs[i]
		idx := strings.Index(line, ": ")
		if idx < 0 {
			continue
		}
		timeStr := line[:idx]
		parts := strings.Split(timeStr, " ")
		if len(parts) < 2 {
			continue
		}
		hm := parts[0]
		ap := parts[1]
		hmParts := strings.Split(hm, ":")
		if len(hmParts) != 2 {
			continue
		}
		h, _ := strconv.Atoi(hmParts[0])
		m, _ := strconv.Atoi(hmParts[1])
		var hour24 int
		if len(ap) > 0 && ap[0] == 'a' {
			if h == 12 {
				hour24 = 0
			} else {
				hour24 = h
			}
		} else {
			if h == 12 {
				hour24 = 12
			} else {
				hour24 = h + 12
			}
		}
		t := hour24*60 + m
		if i > 0 && t < prevTime {
			days++
		}
		prevTime = t
	}
	return days
}

// Embedded copy of testcasesD.txt (each line: n|log1|log2|...).
const testcaseData = `
1|02:42 p.m. msg
3|12:50 p.m. msg|01:05 a.m. msg|05:51 a.m. msg
2|12:17 a.m. msg|09:35 p.m. msg
4|02:26 p.m. msg|02:05 p.m. msg|06:11 p.m. msg|04:40 p.m. msg
3|01:01 a.m. msg|11:54 p.m. msg|02:55 p.m. msg
2|09:43 a.m. msg|04:39 a.m. msg
2|05:37 p.m. msg|04:39 p.m. msg
1|04:17 p.m. msg
2|10:16 p.m. msg|10:27 a.m. msg
4|11:01 a.m. msg|01:22 p.m. msg|10:59 a.m. msg|02:44 p.m. msg
4|07:35 p.m. msg|12:50 a.m. msg|03:23 a.m. msg|09:13 a.m. msg
2|01:19 p.m. msg|05:21 p.m. msg
2|12:36 p.m. msg|09:48 a.m. msg
1|07:57 a.m. msg
2|03:37 a.m. msg|11:53 a.m. msg
4|06:18 a.m. msg|07:43 a.m. msg|08:07 a.m. msg|09:52 p.m. msg
1|08:19 a.m. msg
2|10:15 p.m. msg|09:37 p.m. msg
2|12:30 a.m. msg|05:28 a.m. msg
4|11:34 a.m. msg|09:43 a.m. msg|04:43 a.m. msg|11:52 a.m. msg
2|11:43 p.m. msg|01:51 p.m. msg
5|07:59 a.m. msg|10:00 p.m. msg|02:52 a.m. msg|07:13 p.m. msg|12:39 a.m. msg
3|12:58 a.m. msg|02:21 a.m. msg|11:02 p.m. msg
2|08:07 a.m. msg|11:35 a.m. msg
2|11:11 a.m. msg|02:09 a.m. msg
2|10:43 p.m. msg|09:23 p.m. msg
3|05:15 p.m. msg|05:32 a.m. msg|04:39 p.m. msg
4|10:09 p.m. msg|10:53 p.m. msg|06:48 a.m. msg|04:05 a.m. msg
3|01:53 a.m. msg|09:13 p.m. msg|08:32 p.m. msg
5|11:25 a.m. msg|11:27 p.m. msg|09:43 p.m. msg|08:02 p.m. msg|11:48 p.m. msg
2|08:49 a.m. msg|04:55 p.m. msg
3|09:36 p.m. msg|03:39 p.m. msg|12:01 p.m. msg
4|07:34 p.m. msg|12:53 a.m. msg|09:27 a.m. msg|05:22 a.m. msg
1|06:44 a.m. msg
2|04:54 p.m. msg|02:30 p.m. msg
1|01:07 a.m. msg
2|09:21 p.m. msg|09:44 p.m. msg
3|03:41 p.m. msg|06:27 a.m. msg|04:34 a.m. msg
2|04:26 p.m. msg|10:54 a.m. msg
3|09:16 a.m. msg|08:40 a.m. msg|12:30 p.m. msg
1|12:59 p.m. msg
4|01:07 p.m. msg|06:05 a.m. msg|01:14 p.m. msg|12:28 p.m. msg
3|02:46 a.m. msg|01:23 a.m. msg|09:04 a.m. msg
4|06:14 p.m. msg|04:13 p.m. msg|08:57 a.m. msg|03:18 p.m. msg
2|01:04 p.m. msg|06:30 a.m. msg
3|04:39 a.m. msg|04:44 a.m. msg|04:31 a.m. msg
1|11:43 p.m. msg
3|11:37 a.m. msg|04:42 p.m. msg|02:03 p.m. msg
3|08:29 p.m. msg|01:46 a.m. msg|03:12 p.m. msg
2|09:06 a.m. msg|08:21 a.m. msg
4|12:00 p.m. msg|02:42 p.m. msg|01:45 p.m. msg|08:52 p.m. msg
2|11:23 p.m. msg|01:06 p.m. msg
3|11:39 a.m. msg|10:10 a.m. msg|12:12 p.m. msg
1|08:07 a.m. msg
5|04:33 p.m. msg|07:28 p.m. msg|02:55 p.m. msg|09:35 a.m. msg|10:13 p.m. msg
1|01:54 p.m. msg
2|07:37 a.m. msg|08:21 p.m. msg
2|09:16 a.m. msg|08:02 p.m. msg
4|05:58 p.m. msg|11:29 a.m. msg|04:42 p.m. msg|05:36 p.m. msg
1|04:41 a.m. msg
2|11:52 p.m. msg|01:33 p.m. msg
4|12:59 p.m. msg|09:33 a.m. msg|05:23 a.m. msg|08:24 p.m. msg
2|05:43 a.m. msg|06:29 a.m. msg
4|09:22 p.m. msg|01:00 a.m. msg|12:13 p.m. msg|04:05 a.m. msg
3|12:57 p.m. msg|07:22 a.m. msg|05:01 a.m. msg
3|11:49 a.m. msg|06:02 a.m. msg|01:04 p.m. msg
3|06:45 a.m. msg|03:07 p.m. msg|12:11 a.m. msg
1|02:23 p.m. msg
1|11:08 p.m. msg
2|05:59 a.m. msg|03:14 a.m. msg
4|07:38 p.m. msg|07:18 a.m. msg|06:34 p.m. msg|09:43 a.m. msg
4|11:32 p.m. msg|06:15 a.m. msg|06:38 p.m. msg|12:47 p.m. msg
3|10:50 p.m. msg|09:26 p.m. msg|06:25 a.m. msg
3|08:31 a.m. msg|02:18 p.m. msg|03:37 a.m. msg
2|12:57 p.m. msg|07:21 p.m. msg
2|05:25 p.m. msg|03:10 a.m. msg
2|07:00 a.m. msg|03:26 p.m. msg
4|10:28 p.m. msg|04:33 p.m. msg|01:52 p.m. msg|01:59 p.m. msg
3|12:54 a.m. msg|09:56 a.m. msg|09:02 a.m. msg
4|09:15 p.m. msg|10:18 p.m. msg|09:19 p.m. msg|12:29 a.m. msg
5|02:59 p.m. msg|07:48 p.m. msg|03:11 a.m. msg|02:45 p.m. msg|10:52 p.m. msg
4|05:01 p.m. msg|06:45 a.m. msg|09:34 p.m. msg|07:54 p.m. msg
2|03:19 a.m. msg|11:53 a.m. msg
4|02:38 p.m. msg|05:19 p.m. msg|06:02 p.m. msg|12:57 a.m. msg
1|08:06 a.m. msg
3|11:52 p.m. msg|12:59 a.m. msg|11:24 a.m. msg
4|02:33 p.m. msg|06:38 a.m. msg|10:20 p.m. msg|08:11 p.m. msg
4|06:31 p.m. msg|10:04 p.m. msg|05:25 p.m. msg|02:17 p.m. msg
4|09:46 p.m. msg|01:40 p.m. msg|12:53 p.m. msg|11:16 p.m. msg
2|12:10 a.m. msg|09:28 a.m. msg
4|10:37 a.m. msg|11:12 a.m. msg|03:35 p.m. msg|08:36 a.m. msg
2|05:09 p.m. msg|02:03 a.m. msg
2|09:44 a.m. msg|08:10 a.m. msg
4|08:19 a.m. msg|07:23 p.m. msg|05:25 p.m. msg|09:46 a.m. msg
2|04:15 p.m. msg|11:05 a.m. msg
2|12:12 p.m. msg|09:27 p.m. msg
3|04:02 p.m. msg|06:55 p.m. msg|07:54 p.m. msg
3|10:42 p.m. msg|09:06 p.m. msg|02:58 p.m. msg
4|01:24 a.m. msg|12:57 a.m. msg|04:39 p.m. msg|06:26 p.m. msg
2|10:20 a.m. msg|10:30 a.m. msg
1|04:19 p.m. msg
3|06:12 a.m. msg|12:38 p.m. msg|10:52 a.m. msg
1|10:54 p.m. msg
4|11:44 a.m. msg|07:06 p.m. msg|04:08 p.m. msg|12:16 a.m. msg
1|04:44 a.m. msg
3|07:55 p.m. msg|04:04 p.m. msg|09:06 p.m. msg
3|01:38 p.m. msg|08:01 p.m. msg|12:00 p.m. msg
4|11:33 a.m. msg|05:17 p.m. msg|04:12 p.m. msg|05:19 p.m. msg
2|12:12 p.m. msg|09:06 a.m. msg
2|04:39 p.m. msg|01:01 a.m. msg
3|12:01 a.m. msg|04:05 a.m. msg|06:43 a.m. msg
2|11:57 p.m. msg|04:13 p.m. msg
3|08:04 p.m. msg|08:50 p.m. msg|09:49 p.m. msg
`

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 2 {
			return nil, fmt.Errorf("line %d: invalid format", i+1)
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		if len(parts)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d logs got %d", i+1, n, len(parts)-1)
		}
		logs := parts[1:]
		tests = append(tests, testCase{n: n, logs: logs})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected int) error {
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d\n", tc.n))
	for _, log := range tc.logs {
		input.WriteString(log)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	val, err := strconv.Atoi(got)
	if err != nil {
		return fmt.Errorf("non-integer output %q", got)
	}
	if val != expected {
		return fmt.Errorf("expected %d got %d", expected, val)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := solve(tc)
		if err := runCase(bin, tc, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
