package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases from testcasesA2.txt (one n per line).
const testcasesRaw = `523620
887303
986620
529829
412462
617614
894738
36203
503555
254532
779859
836139
423927
434440
697035
181412
384958
575458
925612
737192
178343
85421
114894
348071
407320
850943
950182
296448
383718
572661
811655
105894
852440
587978
930104
919242
76060
586947
63482
596815
562444
835994
158634
374390
449159
427695
94681
224860
387229
766241
20019
587371
57028
849634
285481
408679
113873
940999
896855
417946
475406
86098
119409
348923
925764
292560
948474
763844
249197
315652
619442
140227
882287
610787
633079
274156
359375
333613
870509
768197
976021
141959
935341
522881
89756
465807
18123
826412
4422
440035
486003
194216
940576
104580
294743
128558
107455
831284
85822
728040
361350
499186`

type testCase struct {
	n int
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
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("line %d parse: %v", idx+1, err)
		}
		cases = append(cases, testCase{n: v})
	}
	return cases, nil
}

// Embedded solver logic from 1786A2.go.
func expected(n int) (aW, aB, bW, bB int) {
	pos := 1
	step := 1
	remaining := n
	for remaining > 0 {
		take := step
		if take > remaining {
			take = remaining
		}
		l := pos
		r := pos + take - 1
		white := (r+1)/2 - (l)/2
		black := take - white
		var alice bool
		if step == 1 {
			alice = true
		} else {
			pair := (step - 2) / 2
			if pair%2 == 1 {
				alice = true
			} else {
				alice = false
			}
		}
		if alice {
			aW += white
			aB += black
		} else {
			bW += white
			bB += black
		}
		pos += take
		remaining -= take
		step++
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA2.go /path/to/binary")
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
		aW, aB, bW, bB := expected(tc.n)
		input := fmt.Sprintf("1\n%d\n", tc.n)
		gotStr, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		parts := strings.Fields(strings.TrimSpace(gotStr))
		if len(parts) != 4 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected four numbers got %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		vals := make([]int, 4)
		for i := 0; i < 4; i++ {
			v, err := strconv.Atoi(parts[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d parse output: %v\n", idx+1, err)
				os.Exit(1)
			}
			vals[i] = v
		}
		if vals[0] != aW || vals[1] != aB || vals[2] != bW || vals[3] != bB {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d %d %d got %v\n", idx+1, aW, aB, bW, bB, vals)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
