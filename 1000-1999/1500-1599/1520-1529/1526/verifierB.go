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
779761980
594763610
148499304
181712557
164248358
222762668
176599970
245081029
31914298
564040371
145552270
879375853
529399916
383298035
966215453
946902837
669317633
790143275
309280043
769279261
851719131
359888332
734646225
127214240
861744124
981538126
917631182
899894001
439018657
830454050
926873475
281071794
965527025
169994475
367728253
674284347
538547545
356256626
572445491
746419701
155242128
979432200
403876137
794353359
798626699
596440875
327751486
251831608
945422466
982284703
406482783
372079637
419043925
506025819
549268771
328735203
440342795
436259808
872889934
104918296
733834940
165371783
158533919
5279891
626519821
639381954
680807467
756408062
946490644
555253265
114690588
748233342
836747993
695093179
749547681
220031306
639720224
691811953
647607588
551314963
119383944
286344700
739317679
785892721
659605102
181072887
403570941
91672866
851850666
38677342
11259105
123908203
958219686
388873341
877637864
775739146
516152881
341829664
116113533
725822208`

func solve(x int64) string {
	if 111*(x%11) <= x {
		return "YES"
	}
	return "NO"
}

type testCase struct {
	x int64
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesRaw)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	if len(fields) != t+1 {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(fields)-1)
	}
	cases := make([]testCase, t)
	for i := 0; i < t; i++ {
		v, err := strconv.ParseInt(fields[i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse case %d: %v", i+1, err)
		}
		cases[i] = testCase{x: v}
	}
	return cases, nil
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

	for i, tc := range cases {
		expected := solve(tc.x)
		input := fmt.Sprintf("1\n%d\n", tc.x)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.ToUpper(strings.TrimSpace(got)) != expected {
			fmt.Printf("case %d failed\nexpected: %s\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
