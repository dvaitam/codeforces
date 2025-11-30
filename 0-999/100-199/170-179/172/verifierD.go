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

type testCase struct {
	a int
	n int
}

// solve embeds the logic from 172D.go.
func solve(tc testCase) string {
	a := tc.a
	n := tc.n
	N := a + n - 1
	rem := make([]int32, N+1)
	for i := 1; i <= N; i++ {
		rem[i] = int32(i)
	}
	lim := int(math.Sqrt(float64(N)))
	for p := 2; p <= lim; p++ {
		sq := p * p
		for j := sq; j <= N; j += sq {
			for rem[j]%int32(sq) == 0 {
				rem[j] /= int32(sq)
			}
		}
	}
	var total uint64
	for x := a; x <= N; x++ {
		total += uint64(rem[x])
	}
	return strconv.FormatUint(total, 10)
}

// Embedded copy of testcasesD.txt.
const testcaseData = `
6 30
17 48
38 17
45 36
25 11
16 41
8 44
38 47
50 34
11 32
42 41
16 1
15 10
46 39
44 50
10 41
10 49
38 40
12 5
7 18
7 18
40 30
37 50
17 38
27 3
46 5
10 29
4 23
22 17
41 27
25 2
19 42
16 14
50 6
38 47
10 16
34 14
29 21
1 49
17 27
38 34
36 9
46 22
46 37
42 35
25 13
20 7
29 42
41 18
16 42
17 50
20 6
24 10
29 2
25 42
29 4
41 50
50 19
33 40
50 28
30 40
37 48
34 13
11 1
18 7
25 46
17 13
32 31
17 33
26 9
39 32
18 12
45 1
14 22
17 46
46 7
43 2
28 20
21 40
47 13
22 27
50 9
42 8
16 23
49 47
31 29
12 44
23 28
28 27
24 38
7 6
36 29
9 7
1 34
45 25
42 31
16 26
41 47
7 27
4 44
`

var expectedOutputs = []string{
	"429",
	"1224",
	"451",
	"1451",
	"231",
	"905",
	"833",
	"1862",
	"1515",
	"611",
	"1630",
	"1",
	"131",
	"1633",
	"2235",
	"779",
	"1086",
	"1513",
	"46",
	"192",
	"192",
	"1047",
	"2070",
	"835",
	"39",
	"99",
	"489",
	"231",
	"358",
	"951",
	"27",
	"1074",
	"182",
	"180",
	"1862",
	"183",
	"429",
	"558",
	"805",
	"597",
	"1265",
	"262",
	"809",
	"1529",
	"1308",
	"269",
	"104",
	"1426",
	"548",
	"962",
	"1355",
	"78",
	"168",
	"59",
	"1240",
	"92",
	"2154",
	"729",
	"1407",
	"1187",
	"1327",
	"1899",
	"382",
	"11",
	"98",
	"1463",
	"181",
	"974",
	"710",
	"195",
	"1156",
	"164",
	"5",
	"376",
	"1216",
	"163",
	"54",
	"561",
	"1050",
	"419",
	"645",
	"309",
	"198",
	"423",
	"2312",
	"867",
	"936",
	"626",
	"690",
	"1045",
	"34",
	"878",
	"67",
	"406",
	"900",
	"1139",
	"513",
	"2033",
	"354",
	"795",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers, got %d", i+1, len(parts))
		}
		a, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad a: %v", i+1, err)
		}
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %v", i+1, err)
		}
		tests = append(tests, testCase{a: a, n: n})
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	input := fmt.Sprintf("%d %d\n", tc.a, tc.n)
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
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := loadTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to load testcases:", err)
		os.Exit(1)
	}
	if len(tests) != len(expectedOutputs) {
		fmt.Fprintf(os.Stderr, "testcase/expected mismatch: %d vs %d\n", len(tests), len(expectedOutputs))
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := runCase(bin, tc, expectedOutputs[i]); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
