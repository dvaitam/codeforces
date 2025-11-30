package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesA1.txt (100 values, one per line).
const testcasesRaw = `140892
596854
888599
841236
800876
66173
267460
123647
519502
797927
471326
495186
683245
398056
827037
220154
98419
511555
29725
936711
656330
188792
889010
66125
529659
606429
563018
393970
222165
735296
209728
397871
579747
212532
240573
343665
309316
986468
954454
384108
357264
439856
884793
23540
101470
284298
372586
322814
578599
290739
474131
943480
419038
435442
546104
410940
398169
330289
959752
897782
523968
346956
37896
278968
759484
65062
94550
115963
901727
114245
80526
773365
52711
248444
971882
967626
680774
495399
787781
333539
286249
253173
272515
718908
921538
861332
771936
127360
269581
86277
79622
735698
100148
176275
336917
261150
521358
949953
104825
971781`

type testCase struct {
	n int64
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		v, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d parse: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: v})
	}
	return cases, nil
}

// Embedded solver logic from 1786A1.go.
func expected(n int64) (int64, int64) {
	var a, b int64
	step := int64(1)
	for n > 0 {
		take := step
		if take > n {
			take = n
		}
		if step%4 == 1 || step%4 == 0 {
			a += take
		} else {
			b += take
		}
		n -= take
		step++
	}
	return a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
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

	for idx, tc := range cases {
		wantA, wantB := expected(tc.n)
		input := fmt.Sprintf("1\n%d\n", tc.n)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		parts := strings.Fields(strings.TrimSpace(gotStr))
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected two numbers got %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		aGot, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse a: %v\n", idx+1, err)
			os.Exit(1)
		}
		bGot, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse b: %v\n", idx+1, err)
			os.Exit(1)
		}
		if aGot != wantA || bGot != wantB {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %d %d\n", idx+1, wantA, wantB, aGot, bGot)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
