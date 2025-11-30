package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesC = `100
980
885
972
871
59
95
88
371
857
175
755
830
687
876
317
259
622
219
623
38
597
699
164
443
655
404
824
742
882
523
974
382
559
960
457
516
276
924
38
893
30
374
478
956
328
931
391
435
915
907
540
170
575
183
243
238
26
182
334
179
141
524
524
370
528
692
575
188
917
458
817
426
754
539
930
932
783
374
810
609
364
372
881
986
458
167
979
774
411
734
758
474
672
545
257
503
287
949
512
514`

type testCase struct {
	input    string
	expected string
}

// solveCase mirrors the logic in 1822C.go.
func solveCase(n int64) string {
	ans := n*(n+2) + 2
	return fmt.Sprintf("%d", ans)
}

func loadCases() ([]testCase, error) {
	fields := strings.Fields(testcasesC)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, fmt.Errorf("bad test count: %w", err)
	}
	pos++
	cases := make([]testCase, 0, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if pos >= len(fields) {
			return nil, fmt.Errorf("case %d: missing n", caseNum+1)
		}
		n, err := strconv.ParseInt(fields[pos], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d: bad n: %w", caseNum+1, err)
		}
		pos++
		input := fmt.Sprintf("1\n%d\n", n)
		cases = append(cases, testCase{
			input:    input,
			expected: solveCase(n),
		})
	}
	return cases, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
