package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt.
const testcasesCData = `
2 10 13 20 3 2
9 25 45 11 38 1
14 42 59 29 35 6
17 37 39 30 24 1
11 38 52 36 11 9
6 24 32 21 6 1
6 17 34 33 12 5
18 32 47 24 34 4
19 44 56 26 11 8
15 51 59 32 18 16
17 52 64 23 30 15
19 57 72 15 43 16
11 24 44 20 31 3
10 45 63 40 33 17
19 48 58 32 14 12
17 43 63 22 51 2
1 16 20 18 2 1
19 36 40 16 18 3
7 13 27 24 2 1
12 26 34 6 1 6
4 11 12 24 1 1
9 20 26 1 17 2
13 53 55 3 10 8
1 26 46 19 8 11
11 45 46 36 29 10
20 25 34 40 28 2
5 38 46 7 21 3
1 32 37 26 19 9
16 51 62 22 56 5
9 28 48 36 2 7
5 11 20 11 5 1
6 15 30 3 17 2
8 25 40 6 17 2
19 36 56 28 17 6
9 45 46 25 3 5
14 27 31 6 24 5
8 17 21 15 6 1
4 20 21 30 11 5
15 37 55 14 14 7
14 44 61 38 38 1
2 31 48 7 12 10
16 42 43 40 8 9
12 33 45 27 2 5
4 13 23 29 1 2
2 31 47 38 14 8
20 27 28 24 1 2
10 17 25 8 7 2
19 45 58 9 30 12
12 40 44 8 8 5
3 27 40 2 7 4
20 53 55 23 19 16
15 27 39 34 31 3
16 45 61 19 44 14
13 30 36 17 20 8
18 48 51 7 37 10
3 28 34 27 5 9
3 11 13 25 5 1
8 32 47 19 34 3
4 16 34 22 4 4
17 35 52 11 11 5
15 33 46 10 37 6
15 46 47 12 25 13
13 48 50 26 18 16
9 38 54 22 36 6
3 20 38 25 26 4
1 24 39 30 23 9
6 15 16 37 4 4
20 47 54 36 25 2
7 27 46 40 32 4
5 8 28 17 16 1
17 56 62 14 46 15
3 28 29 5 18 8
19 53 64 33 18 15
15 19 22 26 6 1
9 20 22 25 16 2
15 36 41 36 19 1
15 18 30 25 18 1
19 50 57 32 20 14
5 38 56 17 5 10
11 33 44 34 26 5
3 38 45 34 20 13
5 40 43 15 3 10
15 53 61 4 18 17
4 14 27 21 7 3
12 19 30 11 12 2
16 47 57 29 9 8
7 27 38 16 7 3
16 31 43 9 23 2
5 22 31 22 13 7
9 44 63 26 48 11
10 47 67 20 24 3
13 46 52 29 23 9
16 24 30 9 13 2
1 10 22 5 12 2
14 17 35 39 16 1
13 50 60 10 41 16
12 35 42 10 7 8
7 31 40 24 27 3
9 17 28 16 8 1
20 25 36 10 4 2
`

// Embedded reference solver for 248C.
const embeddedSolver248C = `package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var y1, y2, yw, xb, yb, r int64
	fmt.Fscan(in, &y1, &y2, &yw, &xb, &yb, &r)

	y1f := float64(y1)
	y2f := float64(y2)
	ywf := float64(yw)
	xbf := float64(xb)
	ybf := float64(yb)
	rf := float64(r)

	h := ywf - rf
	a := 2*h - y2f
	u := 2*h - y1f - rf

	dx := xbf
	dy := ybf - a
	rr := rf * rf
	d2 := dx*dx + dy*dy

	if d2 <= rr {
		fmt.Fprintln(out, -1)
		return
	}

	s := math.Sqrt(d2 - rr)
	gt := a + (-rr*dy+rf*dx*s)/(dx*dx-rr)

	l := math.Max(a+rf, gt)
	if l >= u-1e-12 {
		fmt.Fprintln(out, -1)
		return
	}

	gp := (l + u) / 2
	xw := xbf * (gp - h) / (gp - ybf)

	fmt.Fprintf(out, "%.15f\n", xw)
}
`

type testCase struct {
	y1, y2, yw, xb, yb, r float64
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesCData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 6 {
			return nil, fmt.Errorf("line %d: expected 6 numbers, got %d", idx+1, len(fields))
		}
		vals := make([]float64, 6)
		for i, f := range fields {
			v, err := strconv.ParseFloat(f, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: parse float: %w", idx+1, err)
			}
			vals[i] = v
		}
		cases = append(cases, testCase{
			y1: vals[0], y2: vals[1], yw: vals[2],
			xb: vals[3], yb: vals[4], r: vals[5],
		})
	}
	return cases, nil
}

func buildReference() (string, error) {
	tmpSrc, err := os.CreateTemp("", "248C-src-*.go")
	if err != nil {
		return "", err
	}
	srcPath := tmpSrc.Name()
	if _, err := tmpSrc.WriteString(embeddedSolver248C); err != nil {
		tmpSrc.Close()
		os.Remove(srcPath)
		return "", err
	}
	tmpSrc.Close()
	defer os.Remove(srcPath)

	tmp, err := os.CreateTemp("", "248C-ref-*")
	if err != nil {
		return "", err
	}
	binPath := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(binPath)
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return binPath, nil
}

func runProgram(bin string, tc testCase) (float64, error) {
	input := fmt.Sprintf("%d %d %d %d %d %d\n", int(tc.y1), int(tc.y2), int(tc.yw), int(tc.xb), int(tc.yb), int(tc.r))
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
	}
	str := strings.TrimSpace(out.String())
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("parse output %q: %v", str, err)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProgram(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp == -1 {
			if got != -1 {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %f\n", i+1, got)
				os.Exit(1)
			}
			continue
		}
		if math.Abs(got-exp) > 1e-6*math.Max(1.0, math.Abs(exp)) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.10f got %.10f\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
