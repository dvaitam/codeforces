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

type op struct {
	add bool
	x   int64
	y   int64
}

type testcase struct {
	ops []op
}

const testcasesRaw = `100
+ 4 8
? 15 23
? 18 13
+ 12 14
? 25 27
? 23 32
+ 4 3
? 29 28
? 16 17
+ 1 17
? 24 37
? 27 24
+ 13 4
? 25 18
? 21 36
+ 6 10
? 20 33
? 15 34
+ 11 1
? 16 17
? 14 27
+ 11 15
? 18 24
? 23 35
+ 15 9
? 29 22
? 26 29
+ 1 4
? 12 20
? 22 23
+ 10 18
? 15 34
? 21 20
+ 4 18
? 31 27
? 13 25
+ 4 6
? 26 23
? 17 32
+ 30 23
? 26 40
? 9 29
+ 1 33
? 11 33
? 19 32
+ 5 21
? 19 15
? 14 18
+ 23 17
? 42 33
? 11 43
+ 18 21
? 17 39
? 19 40
+ 9 16
? 28 31
? 19 23
+ 20 5
? 27 20
? 13 13
+ 25 8
? 20 26
? 14 22
+ 6 2
? 24 12
? 16 20
+ 28 8
? 17 37
? 23 31
+ 25 12
? 22 40
? 21 38
+ 29 7
? 19 29
? 10 27
+ 19 27
? 18 45
? 9 37
+ 5 1
? 21 27
? 16 18
+ 10 32
? 20 33
? 19 24
+ 16 38
? 28 42
? 26 34
+ 30 32
? 23 37
? 36 42
+ 27 26
? 34 44
? 21 33
+ 26 18
? 23 35
? 30 33
+ 32 14
? 34 45
? 24 31
+ 14 3
? 34 26
? 19 30
+ 37 1
? 28 44
? 24 27
+ 14 27
? 39 45
? 20 35
+ 37 27
? 37 33
+ 10 16
? 21 32
? 30 32
+ 19 20`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	scanner.Split(bufio.ScanWords)

	nextInt := func() int {
		if !scanner.Scan() {
			panic("unexpected EOF while reading testcases")
		}
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		return v
	}

	t := nextInt()
	cases := make([]testcase, 0, t)
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			panic(fmt.Sprintf("missing op marker at case %d", i+1))
		}
		marker := scanner.Text()
		if len(marker) != 2 || marker[0] != '+' && marker[0] != '?' || marker[1] != '+' && marker[1] != '?' {
			// marker combines two chars to avoid newline splitting; treat marker[0] as op and marker[1] as indicator of test count start
		}
		opChar := marker[0]
		x := nextInt()
		y := nextInt()
		ops := []op{{add: opChar == '+', x: int64(x), y: int64(y)}}
		cases = append(cases, testcase{ops: ops})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	return cases
}

// solve replicates 1101E.go behavior over a sequence of operations.
func solve(tc testcase) []string {
	var mxx, mxy int64
	res := []string{}
	for _, op := range tc.ops {
		x, y := op.x, op.y
		if x > y {
			x, y = y, x
		}
		if op.add {
			if x > mxx {
				mxx = x
			}
			if y > mxy {
				mxy = y
			}
		} else {
			if mxx <= x && mxy <= y {
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		}
	}
	return res
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseCandidateOutput(out string) ([]string, error) {
	fields := strings.Fields(out)
	return fields, nil
}

func checkCase(bin string, tc testcase) error {
	// Build the single-line input and expected outputs
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.ops)))
	for _, op := range tc.ops {
		if op.add {
			sb.WriteString("+ ")
		} else {
			sb.WriteString("? ")
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", op.x, op.y))
	}

	expected := solve(tc)
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := parseCandidateOutput(out)
	if err != nil {
		return err
	}
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d outputs, got %d", len(expected), len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			return fmt.Errorf("query %d: expected %s got %s", i+1, expected[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
