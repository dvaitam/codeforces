package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Embedded copy of testcasesB.txt.
const testcasesBData = `
906691060
413654000
813847340
955892129
451585302
43469774
278009743
548977049
521760890
434794719
985946605
841597327
891047769
325679555
511742082
384452588
626401696
957413343
975078789
234551095
541903390
149544007
302621085
150050892
811538591
101823754
663968656
858351977
268979134
976832603
571835845
757172937
869964136
646287426
968693315
157798603
333018423
106046332
783650879
79180333
965120264
913189318
734422155
354546568
506959382
601095368
108127102
379880546
466188457
339513622
655934895
687649392
980338160
219556307
593267778
512185346
475338373
929119464
559799207
279701489
66872193
864392047
986194170
589161386
983541587
15077163
100149904
772777020
902041077
428233517
762628806
885670548
842938613
717424033
671374074
1227090
657019496
529975200
889126175
931581387
357701129
261897307
784130655
349185523
755530427
934661371
67628852
205156724
984641620
609360020
238052748
256211902
862585180
153002189
862407392
583031025
481003666
97942385
86378037
343656009
`

type testCase struct {
	x int
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesBData, "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var x int
		if _, err := fmt.Sscan(line, &x); err != nil {
			return nil, fmt.Errorf("line %d: parse int: %w", idx+1, err)
		}
		cases = append(cases, testCase{x: x})
	}
	return cases, nil
}

func shares(d int, has []bool) bool {
	if d == 0 {
		return has[0]
	}
	for d > 0 {
		if has[d%10] {
			return true
		}
		d /= 10
	}
	return false
}

func solve(x int) int {
	s := fmt.Sprint(x)
	has := make([]bool, 10)
	for _, ch := range s {
		has[ch-'0'] = true
	}
	ans := 0
	for i := 1; i*i <= x; i++ {
		if x%i == 0 {
			d1 := i
			d2 := x / i
			if shares(d1, has) {
				ans++
			}
			if d2 != d1 && shares(d2, has) {
				ans++
			}
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse embedded testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range tests {
		input := fmt.Sprintf("%d\n", tc.x)
		want := fmt.Sprintf("%d", solve(tc.x))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
