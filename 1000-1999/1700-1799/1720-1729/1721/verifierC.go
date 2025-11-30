package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
1
6
8
6
11 17 20 39 43 48
17 21 36 48 57 61
7
18 24 29 33 33 35 47
19 24 40 43 47 47 60
9
2 11 12 12 12 15 16 21 36
6 21 23 27 28 28 32 35 49
9
11 23 24 24 26 29 38 49 49
25 31 37 40 41 43 53 65 65
6
23 30 30 37 43 47
40 44 44 45 52 53
10
18 20 20 31 33 33 34 36 46 50
38 38 39 39 42 44 49 52 57 69
2
22 47
22 53
2
4 37
24 38
5
7 15 38 44 49
19 23 46 51 55
1
28
29
1
24
35
3
2 16 44
4 19 46
1
3
3
6
9 11 12 17 34 48
9 18 23 30 41 52
1
1
12
10
2 8 19 20 22 29 32 41 48 48
19 20 27 28 34 36 48 50 55 56
6
2 7 9 29 34 38
14 22 25 38 39 48
5
2 17 27 39 42
19 21 28 43 47
3
7 11 11
18 21 31
9
3 5 6 15 16 17 29 46 46
12 21 25 25 27 34 42 54 62
1
10
11
7
6 8 11 16 27 33 47
9 11 11 21 34 36 53
1
34
48
8
14 14 20 25 35 42 44 49
25 27 27 36 45 53 60 62
9
2 7 8 12 24 31 34 38 43
17 18 21 23 31 33 41 46 47
5
2 13 29 44 50
3 26 49 59 64
4
1 5 38 40
5 10 49 49
2
15 49
30 55
2
24 37
36 51
3
23 26 49
26 34 52
2
6 40
16 60
7
2 7 14 31 40 43 45
3 22 23 42 47 54 56
5
31 31 34 47 47
43 44 46 54 59
3
17 32 39
34 41 45
10
5 5 7 10 12 23 27 35 37 47
7 8 14 21 25 34 35 45 51 52
9
7 8 10 19 22 28 34 35 49
14 18 24 24 27 41 42 47 60
10
2 10 12 25 29 30 39 47 47 48
13 14 26 37 40 42 47 60 63 67
6
6 22 36 43 46 48
13 39 49 55 58 60
1
21
35
9
2 7 12 14 26 30 42 46 47
20 20 24 26 29 42 52 55 59
10
1 9 13 28 31 32 38 40 40 44
9 25 31 33 38 40 40 45 51 59
9
5 22 30 32 38 43 43 44 49
13 32 38 40 44 49 54 61 62
5
4 9 41 44 47
9 24 53 56 58
3
1 19 36
15 19 47
1
35
47
10
9 14 20 29 31 32 35 42 44 45
16 18 28 39 40 42 51 55 64 65
7
6 14 26 33 34 39 41
18 22 36 42 42 48 53
4
15 30 34 36
23 31 37 39
7
5 14 21 22 23 24 30
16 19 32 34 36 36 38
8
7 11 14 16 18 21 31 41
13 19 22 22 25 27 38 49
9
18 22 25 26 33 39 41 47 48
32 36 37 42 46 50 56 60 67
6
12 17 20 23 26 31
22 26 28 32 36 43
3
2 7 23
7 18 25
7
1 16 21 25 35 35 39
10 29 31 41 45 45 46
8
7 10 10 14 17 22 24 27
12 15 20 20 24 28 29 43
6
4 10 24 40 42 50
9 12 37 50 54 54
6
8 22 34 37 42 50
27 34 38 41 54 65
8
6 21 33 35 38 39 40 40
23 36 43 45 49 52 52 56
8
3 7 8 9 29 33 38 44
8 9 18 20 33 43 46 47
6
2 10 12 15 28 43
4 20 29 29 32 58
4
5 9 31 36
5 13 47 53
1
4
10
9
1 9 16 22 24 32 34 34 44
1 13 25 31 33 35 40 45 48
10
14 22 25 26 33 34 40 41 42 46
19 28 30 38 40 45 48 53 53 55
6
9 10 21 26 28 28
13 18 29 36 38 43
5
15 18 32 34 49
22 28 37 49 49
10
3 9 13 13 14 26 36 38 42 49
9 21 23 27 28 31 42 43 59 61
1
45
51
2
9 41
26 46
2
30 44
50 53
4
11 18 21 45
23 27 36 58
7
3 8 20 30 43 43 47
16 20 23 36 53 55 56
9
9 13 26 28 34 38 47 48 50
14 27 39 40 46 53 56 66 69
4
13 29 31 38
28 41 47 47
2
11 24
30 44
8
9 15 37 40 42 43 44 50
18 21 43 44 45 49 54 56
6
4 17 21 35 43 47
14 23 31 47 48 58
5
9 10 14 37 38
14 26 29 39 58
10
6 10 16 18 20 29 31 32 37 42
13 16 28 35 38 43 48 48 49 50
4
4 17 36 43
11 21 55 55
7
8 19 25 26 26 30 31
14 26 27 32 33 43 46
10
1 4 15 25 26 28 33 35 44 46
7 9 26 37 41 44 46 48 53 60
4
2 2 18 31
3 6 35 38
4
4 16 21 35
20 23 41 44
2
36 42
44 53
3
9 28 46
10 37 62
5
4 23 31 36 49
14 26 47 50 52
10
2 10 18 19 23 24 31 33 39 41
3 20 23 31 34 35 39 41 48 58
9
11 12 16 27 28 28 40 40 50
12 12 33 34 37 38 40 41 57
10
5 5 8 15 15 16 24 26 36 39
8 8 11 17 27 32 32 44 46 56
1
34
50
9
9 9 10 24 26 28 28 30 48
10 15 27 40 40 41 47 50 62
3
9 32 39
20 36 39
5
10 12 27 41 46
28 32 35 55 61
8
13 15 18 23 28 28 49 50
16 28 30 33 41 42 50 58
9
16 17 18 30 30 31 37 44 46
27 33 37 37 44 48 49 54 60
5
8 24 27 33 49
15 36 39 44 62
10
6 25 29 30 30 34 40 40 48 50
17 34 35 40 41 42 45 60 61 64
8
12 17 21 26 27 37 46 47
21 27 29 36 37 48 49 64`

type testCase struct {
	n int
	a []int
	b []int
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func solveCase(a, b []int) (string, string) {
	n := len(a)
	mn := make([]int, n)
	j := 0
	for i := 0; i < n; i++ {
		for j < n && b[j] < a[i] {
			j++
		}
		mn[i] = b[j] - a[i]
	}
	mx := make([]int, n)
	pos := n - 1
	for i := n - 1; i >= 0; i-- {
		mx[i] = b[pos] - a[i]
		if i > 0 && b[i-1] < a[i] {
			pos = i - 1
		}
	}
	sb1 := strings.TrimSpace(strings.Join(intsToStrs(mn), " "))
	sb2 := strings.TrimSpace(strings.Join(intsToStrs(mx), " "))
	return sb1, sb2
}

func intsToStrs(a []int) []string {
	res := make([]string, len(a))
	for i, v := range a {
		res[i] = strconv.Itoa(v)
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	if len(lines)-1 != t*3 {
		return nil, fmt.Errorf("expected %d*3 lines for cases, got %d", t, len(lines)-1)
	}
	var cases []testCase
	for i := 0; i < t; i++ {
		nLine := strings.TrimSpace(lines[1+i*3])
		aLine := strings.TrimSpace(lines[2+i*3])
		bLine := strings.TrimSpace(lines[3+i*3])
		n, err := strconv.Atoi(nLine)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", i+1, err)
		}
		aFields := strings.Fields(aLine)
		bFields := strings.Fields(bLine)
		if len(aFields) != n || len(bFields) != n {
			return nil, fmt.Errorf("case %d: expected %d elements, got %d and %d", i+1, n, len(aFields), len(bFields))
		}
		aVals := make([]int, n)
		bVals := make([]int, n)
		for idx, f := range aFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse a[%d]: %v", i+1, idx, err)
			}
			aVals[idx] = v
		}
		for idx, f := range bFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("case %d: parse b[%d]: %v", i+1, idx, err)
			}
			bVals[idx] = v
		}
		cases = append(cases, testCase{n: n, a: aVals, b: bVals})
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for caseNum, tc := range cases {
		exp1, exp2 := solveCase(tc.a, tc.b)
		var input strings.Builder
		fmt.Fprintf(&input, "1\n%d\n", tc.n)
		fmt.Fprintln(&input, strings.Join(intsToStrs(tc.a), " "))
		fmt.Fprintln(&input, strings.Join(intsToStrs(tc.b), " "))
		got, err := runExe(bin, input.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", caseNum+1, err)
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(got), "\n")
		if len(outLines) < 2 {
			fmt.Printf("case %d: output too short\n", caseNum+1)
			os.Exit(1)
		}
		g1 := strings.TrimSpace(outLines[0])
		g2 := strings.TrimSpace(outLines[1])
		if g1 != exp1 || g2 != exp2 {
			fmt.Printf("case %d failed\nexpected:\n%s\n%s\n got:\n%s\n%s\n", caseNum+1, exp1, exp2, g1, g2)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
