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

// Embedded copy of testcasesE.txt so the verifier is self-contained.
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
+ 9 6
? 26 36
? 12 35
+ 13 1
? 15 21
? 21 20
+ 13 7
? 20 30
? 19 32
+ 20 17
? 34 31
? 22 33
+ 1 19
? 30 39
? 24 29
+ 12 17
? 27 35
? 37 34
+ 9 12
? 25 25
? 22 27
+ 15 11
? 37 20
? 28 35
+ 20 10
? 24 20
? 20 34
+ 17 6
? 23 32
? 36 30
+ 18 5
? 18 35
? 30 34
+ 15 12
? 17 26
? 35 26
+ 3 16
? 21 22
? 33 37
+ 13 20
? 25 26
? 19 33
+ 1 18
? 32 30
? 17 26
+ 20 10
? 30 34
? 33 29
+ 3 2
? 19 31
? 29 40
+ 14 13
? 37 40
? 24 29
+ 16 13
? 34 20
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

	var cases []testcase
	for scanner.Scan() {
		q, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", scanner.Text(), err))
		}
		ops := make([]op, 0, q)
		for i := 0; i < q; i++ {
			if !scanner.Scan() {
				panic(fmt.Sprintf("missing op marker at operation %d", i+1))
			}
			marker := scanner.Text()
			if marker != "+" && marker != "?" && marker != "++" && marker != "+?" && marker != "?+" && marker != "??" {
				panic(fmt.Sprintf("invalid op marker %q", marker))
			}
			add := marker[0] == '+'
			x := nextInt()
			y := nextInt()
			ops = append(ops, op{add: add, x: int64(x), y: int64(y)})
		}
		cases = append(cases, testcase{ops: ops})
	}

	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	return cases
}

// solve replicates the logic from 1101E.go for a single testcase.
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
