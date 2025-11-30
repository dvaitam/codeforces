package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const testcasesRaw = `241 77631706
1743915 0
9 6
7083 78353374
7 124158683
59786 0736625851
91286570 04999622
403883685 5748906828
460759838 7565088
6790328921 401107043
29254 22
92447 771046
734148 9603
7 027
783873806 5604230411
52694 080
4972980635 39
4716487 5
7402359256 4168
978310 22
53099 6
7 490007020
1840185 89
8 1569
79059 90
6 10651
62138620 424
481121 1
54137 46
95468 1
9 22440
8052 96
143 33374
94421770 9462735393
3 3
8198905 7
932800199 0
7360863 6
627365 7
70023620 3
3 34364536
992269888 6
230418 6
7558 18
8 945444520
736097 5488445335
6 29980
4 73
94182 9
15 6494
7127878342 37
17069 0
1447553767 7
778 5
21 8285793
1 8
508 0534255617
6 4039595495
28602541 1412418
8 23445
229488397 2
426 727733
464 1
842 77
39 30882841
7 56
8 62
8 6
26424 6981562614
8 7829
3070074664 6
9488 109338
7 5251287
6 9
23 7245
8676 7
43919 594
47 3165913
763 15
84302 5
0 0
4694979 9750004
1 55
96779 6
74733 51881185
134 674437
5 43
46 30575099
6 579264
2824727 61
1984782720 7036235699
4292 3
73457 772
69083 62
6 0
6 5868316
12772749 2
3312856 2
56196 26028
149928549 0
23 33
8372 00
2004 3043300995
4970 4053`

type testB struct {
	a string
	s string
}

func parseTestcases() ([]testB, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testB
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 fields got %d", i+1, len(parts))
		}
		cases = append(cases, testB{a: parts[0], s: parts[1]})
	}
	return cases, nil
}

func solveB(tc testB) string {
	digits := []byte(tc.s)
	sort.Slice(digits, func(i, j int) bool { return digits[i] > digits[j] })
	res := []byte(tc.a)
	j := 0
	for i := 0; i < len(res) && j < len(digits); i++ {
		if digits[j] > res[i] {
			res[i] = digits[j]
			j++
		}
	}
	return string(res)
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

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer cleanup()

	cases, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%s\n%s\n", tc.a, tc.s)
		expect := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
