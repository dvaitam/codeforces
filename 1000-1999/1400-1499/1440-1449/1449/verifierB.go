package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesB.txt so the verifier is self-contained.
const testcasesRaw = `100
374
584
568
205
964
517
424
497
833
366
425
355
2
552
554
639
806
628
340
470
615
29
824
236
651
182
564
599
186
882
94
818
565
817
872
837
954
262
34
862
967
690
73
86
889
18
464
15
773
774
288
256
276
113
817
640
190
353
298
72
172
164
262
541
975
173
673
280
664
729
302
466
720
330
509
486
117
25
320
396
352
432
816
193
265
112
260
922
748
523
215
989
621
443
837
999
22
231
19
407`

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func parseTestcases() ([]int, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("invalid t: %v", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("expected %d cases, found %d numbers", t, len(fields)-1)
	}
	nums := make([]int, t)
	for i := 0; i < t; i++ {
		v, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse n%d: %v", i+1, err)
		}
		nums[i] = v
	}
	return nums, nil
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(cases))
	for _, n := range cases {
		fmt.Fprintf(&input, "%d\n", n)
	}
	inputStr := input.String()

	var expected strings.Builder
	for i, n := range cases {
		if i > 0 {
			expected.WriteByte('\n')
		}
		if isPrime(n) {
			expected.WriteString("YES")
		} else {
			expected.WriteString("NO")
		}
	}
	expectedStr := expected.String()

	got, err := runCandidate(bin, inputStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expectedStr) {
		fmt.Printf("verifier failed\nexpected:\n%s\n\ngot:\n%s\n", expectedStr, got)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
