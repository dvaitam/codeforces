package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded test cases from testcasesA.txt. Each line is a single input string.
const testcasesRaw = `6604876475938242194892411578
5659
387784080160975351393
28711587
4841
583989471965934232
94
1122018684833969
775159179
533041352560123098910
3991
1510903217300
691413145620870916
4579230
22584197207698456428071508
237594599
46610
35233769606960271427
78900754706381206
6503008913193442176104
567
7
436191909245673732
2
28406054255509276
2457533532183
6042
584093185
29036
41659252
539
5592
24196837510
10674844855889707893234512
758686
7860
1358258973563579705
1138136
35608867636
4091403609316988
29654911
3366765908723944510269961
460243257169210358
1663528216083703438415
3
35
6576020356469628311782303
8241
683849750
24860485216624734
92369
42932514993
720827
617851291368218
8758374
989
526260402291
7199110181
865300237672455
8
411820288269040
31636987490180663941728
507237873
1225858144302185
58080
609701555803294
923238689291747295902
48340
63891054570
9984538849
3928055496239408504
5149389097973410604
16491374255946650
63543510666359573
7180842598501
792006003
2677597068
36194023659412902
133640781
4027783875903954790348114
2108526769653
31329371750767067
305520129032
6671
60010896481
6581671867
91621495830367511298557
4
86554431155
67715160899518
6932033009445404578996714
430624047002
1520674707143785
1
1358
4306
4990
15223946240803016
15204037050217264
78401
356548667955265315628143
250063852
2862`

type testCase struct {
	line string
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

func parseTestcases() []testCase {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cases = append(cases, testCase{line: line})
	}
	return cases
}

// Embedded solver logic from 175A.go.
func parseSegment(s string) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	if len(s) > 1 && s[0] == '0' {
		return 0, false
	}
	if len(s) > 7 {
		return 0, false
	}
	if len(s) == 7 && s > "1000000" {
		return 0, false
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return v, true
}

func expected(tc testCase) string {
	s := tc.line
	n := len(s)
	maxSum := -1
	for i := 1; i <= n-2; i++ {
		for j := i + 1; j <= n-1; j++ {
			a := s[:i]
			b := s[i:j]
			c := s[j:]
			va, ok := parseSegment(a)
			if !ok {
				continue
			}
			vb, ok := parseSegment(b)
			if !ok {
				continue
			}
			vc, ok := parseSegment(c)
			if !ok {
				continue
			}
			sum := va + vb + vc
			if sum > maxSum {
				maxSum = sum
			}
		}
	}
	return fmt.Sprintf("%d", maxSum)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	cases := parseTestcases()

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		want := expected(tc)
		got, err := runCandidate(bin, tc.line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
