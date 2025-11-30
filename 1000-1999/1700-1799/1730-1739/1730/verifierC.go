package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCase struct {
	s string
}

// solve embeds the logic from 1730C.go.
func solve(s string) string {
	digits := []byte(s)
	n := len(digits)
	minRight := byte('9')
	for i := n - 1; i >= 0; i-- {
		orig := digits[i]
		if orig > minRight && orig < '9' {
			digits[i] = orig + 1
		}
		if orig < minRight {
			minRight = orig
		}
	}
	sort.Slice(digits, func(i, j int) bool { return digits[i] < digits[j] })
	return string(digits)
}

// Embedded copy of testcasesC.txt.
const testcaseData = `
155257803
29546
4
606463522
95
069
7116882
232
0827
945164
585
268436567
466122099
139981492
1001575
75
4732881820
869
7640
2627906196
387079
48
86742385
492
11769681
73142
120298
0068
7178551033
44307585
14
6356234375
69216
873556
56417
17253
573
66
87993684
35801756
6090
1435212
030068
2
99
1730
961823366
06360940
5557298141
40
926
9536
818171877
785
646
98
373576
7
8
7363477
374
792813
81025
5
2763673260
835288
9629793970
7958
0002
7931
645732585
67571231
678
5006
1699
8302874
8811970
6954
03
0985274197
36685
9410
599
18765601
9220163
4
61545
3880588027
4
251986
49
81
2708351832
6397
245833
52480154
019041
438
571688611
7524364000
`

var expectedOutputs = []string{
	"023366689",
	"24669",
	"4",
	"022456777",
	"59",
	"069",
	"1127899",
	"224",
	"0279",
	"145679",
	"559",
	"235567779",
	"023357799",
	"112459999",
	"0012558",
	"58",
	"0233458999",
	"699",
	"0578",
	"0133677899",
	"047899",
	"48",
	"23557899",
	"259",
	"11177899",
	"12458",
	"022389",
	"0068",
	"0223366889",
	"04555589",
	"14",
	"2334556778",
	"13679",
	"355689",
	"15677",
	"12368",
	"368",
	"66",
	"34789999",
	"01456689",
	"0079",
	"1123456",
	"000468",
	"2",
	"99",
	"0248",
	"123366799",
	"00045779",
	"1135666899",
	"05",
	"269",
	"3669",
	"111778999",
	"589",
	"467",
	"89",
	"335688",
	"7",
	"8",
	"3347778",
	"348",
	"133899",
	"02259",
	"5",
	"0334477788",
	"246889",
	"0347889999",
	"5889",
	"0002",
	"1489",
	"245556789",
	"11346788",
	"678",
	"0066",
	"1699",
	"0244899",
	"0228999",
	"4679",
	"03",
	"0135678999",
	"35779",
	"0259",
	"599",
	"01267789",
	"0133379",
	"4",
	"14567",
	"0024679999",
	"4",
	"136699",
	"49",
	"19",
	"0123446899",
	"3779",
	"233569",
	"01345669",
	"001259",
	"358",
	"111677899",
	"0003455678",
}

func loadTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcaseData), "\n")
	tests := make([]testCase, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tests = append(tests, testCase{s: line})
		_ = i
	}
	return tests, nil
}

func runCase(bin string, tc testCase, expected string) error {
	input := fmt.Sprintf("1\n%s\n", tc.s)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
