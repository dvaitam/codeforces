package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `980
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	n int
}

func solveCase(tc testCase) (int, []int) {
	n := tc.n
	result := make([]int, 0)
	prod := 1 % n
	for i := 1; i < n; i++ {
		if gcd(i, n) == 1 {
			result = append(result, i)
			prod = prod * i % n
		}
	}
	if prod != 1 {
		filtered := make([]int, 0, len(result)-1)
		for _, v := range result {
			if v != prod {
				filtered = append(filtered, v)
			}
		}
		result = filtered
	}
	return len(result), result
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	res := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("case %d: parse n: %v", idx+1, err)
		}
		res = append(res, testCase{n: n})
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return res, nil
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

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
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

	for i, tc := range cases {
		cnt, values := solveCase(tc)
		expected := strconv.Itoa(cnt)
		if cnt > 0 {
			expected = expected + "\n" + strings.TrimSpace(strings.Trim(fmt.Sprint(values), "[]"))
		}

		input := fmt.Sprintf("%d\n", tc.n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
