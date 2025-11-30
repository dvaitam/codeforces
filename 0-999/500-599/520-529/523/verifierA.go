package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testcase struct {
	w, h  int
	grid  []string
}

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `100
4 4
.***
***.
.*..
*.*.
1 3
*
.
*
4 3
.***
...*
.**.
3 1
...
2 5
*.
.*
*.
**
.*
5 2
**.**
.*...
1 5
*
*
.
.
.
2 1
.*
5 3
..***
**.*.
**...
3 1
.*.
3 4
...
...
...
.*.
3 1
...
2 1
*.
1 1
*
5 1
*...*
3 4
..*
..*
.**
*..
1 2
.
*
5 3
.*..*
*****
..*.*
1 5
*
.
.
*
*
5 3
*.***
...*.
..***
1 4
*
.
.
*
1 3
.
*
*
5 5
..***
*.*..
..***
.....
...**
2 1
**
1 1
*
4 1
*.*.
2 3
..
..
**
3 5
.**
**.
.**
...
***
3 1
*..
3 2
.**
*.*
1 4
.
.
*
.
5 1
*****
1 1
*
4 3
*.*.
**.*
*..*
1 1
.
2 2
.*
..
4 5
***.
*.*.
*.*.
*...
*..*
4 3
....
..**
*..*
2 5
**
**
**
..
.*
3 3
..*
.**
*..
5 1
....*
1 4
*
.
.
*
3 5
*..
.*.
***
*..
..*
5 5
**.**
****.
...**
..**.
*..**
1 5
*
*
*
.
.
4 2
...*
.*.*
4 1
.**.
5 3
.**..
.***.
**...
4 3
....
**..
**..
4 1
.*..
4 4
*.**
*.*.
.*.*
*...
3 1
...
2 1
.*
3 1
...
4 1
***.
1 2
*
*
4 3
**..
.*.*
...*
5 2
..***
...*.
1 2
*
*
5 4
.***.
.**.*
*.***
.***.
2 1
*.
5 5
.**.*
*..**
.**.*
.*.**
..*.*
5 5
.**.*
.**..
*..**
**.**
*....
3 3
*..
..*
**.
5 5
..***
***..
..***
.****
.****
1 4
*
.
*
.
5 2
.***.
..*.*
1 1
.
3 3
..*
..*
***
2 3
**
.*
.*
2 2
.*
**
4 4
*..*
..*.
...*
...*
5 1
****.
5 5
*****
.**..
..*.*
...*.
*.*.*
5 1
.****
3 1
..*
1 2
*
.
2 3
**
**
**
2 1
*.
1 4
*
*
.
*
4 4
*...
...*
.*.*
*...
4 2
.***
*.*.
3 4
***
...
..*
..*
1 5
*
.
.
*
*
1 5
*
.
*
.
.
1 5
*
*
.
*
*
3 5
...
**.
.*.
...
.*.
2 2
.*
*.
4 1
.***
3 1
..*
4 5
*.*.
*.**
*...
*.**
.**.
3 1
*..
3 1
*..
2 1
..
3 1
.**
2 2
**
.*
3 2
***
.**
5 5
*..**
.....
*....
...*.
*..**
4 1
.*.*
3 2
*.*
.**`

func parseTestcases() ([]testcase, error) {
	reader := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read t: %v", err)
	}
	cases := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		var w, h int
		if _, err := fmt.Fscan(reader, &w, &h); err != nil {
			return nil, fmt.Errorf("case %d header: %v", i+1, err)
		}
		lines := make([]string, h)
		for j := 0; j < h; j++ {
			if _, err := fmt.Fscan(reader, &lines[j]); err != nil {
				return nil, fmt.Errorf("case %d line %d: %v", i+1, j+1, err)
			}
		}
		cases = append(cases, testcase{w: w, h: h, grid: lines})
	}
	return cases, nil
}

// Embedded solver logic from 523A.go.
func solve(w, h int, lines []string) string {
	img := make([][]byte, h)
	for i := 0; i < h; i++ {
		img[i] = []byte(lines[i])
	}
	outRows := 2 * w
	outCols := 2 * h
	var b strings.Builder
	rowBuf := make([]byte, outCols)
	for r2 := 0; r2 < outRows; r2++ {
		r := r2 / 2
		for c2 := 0; c2 < outCols; c2++ {
			c := c2 / 2
			rowBuf[c2] = img[c][r]
		}
		b.Write(rowBuf)
		if r2+1 < outRows {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expect := strings.TrimSpace(solve(tc.w, tc.h, tc.grid))

		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", tc.w, tc.h)
		for _, line := range tc.grid {
			input.WriteString(line)
			input.WriteByte('\n')
		}

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n%s", idx+1, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
