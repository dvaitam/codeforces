package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
5 914
2 212
1 209
1 765
1 300
3 706
1 955
5 238
2 191
4 115
1 464
4 988
1 910
2 287
1 1000
2 127
1 775
3 997
1 656
3 991
4 937
1 1000
3 975
4 316
2 250
4 912
1 415
2 954
2 447
2 488
3 771
5 216
1 743
4 618
4 835
2 593
5 665
3 856
4 676
1 965
2 510
4 112
3 443
1 260
1 293
4 162
5 903
1 268
1 1000
1 984
3 899
1 1000
3 214
2 393
4 967
4 678
2 185
2 262
2 453
3 884
4 443
5 578
3 435
1 646
1 995
2 376
4 647
4 555
3 632
2 795
3 501
1 420
5 413
3 299
4 11
2 950
1 813
2 16
5 865
3 914
2 308
5 697
3 715
2 496
5 915
3 81
5 990
3 800
2 802
4 99
2 471
2 338
2 232
5 731
1 500
3 860
3 890
3 764
2 386
2 487
1 495
2 43
4 709
`

type testCase struct {
	n int
	k int64
}

func parseTestcases(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers got %d", idx+1, len(parts))
		}
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n", idx+1)
		}
		k, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad k", idx+1)
		}
		cases = append(cases, testCase{n: n, k: k})
	}
	return cases, nil
}

func expectedOutput(tc testCase) string {
	n, k := tc.n, tc.k
	if n == 1 {
		return strconv.FormatInt(k, 10)
	}
	var x int64
	for i := 0; ; i++ {
		bit := int64(1) << i
		if x|bit > k {
			break
		}
		x |= bit
	}
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(x, 10))
	sb.WriteByte(' ')
	sb.WriteString(strconv.FormatInt(k-x, 10))
	for i := 2; i < n; i++ {
		sb.WriteString(" 0")
	}
	// The original solution prints fmt.Fprint(writer, '\\n') which renders as "10".
	sb.WriteString("10")
	return sb.String()
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expected := expectedOutput(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
