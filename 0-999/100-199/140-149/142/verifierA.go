package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// divisors returns all positive divisors of n.
func divisors(n int64) []int64 {
	var ds []int64
	lim := int64(math.Sqrt(float64(n)))
	for i := int64(1); i <= lim; i++ {
		if n%i == 0 {
			ds = append(ds, i)
			if j := n / i; j != i {
				ds = append(ds, j)
			}
		}
	}
	return ds
}

// solve mirrors 142A.go logic.
func solve(n int64) (int64, int64) {
	const inf64 int64 = 1<<63 - 1
	minStolen := inf64
	maxStolen := int64(0)
	for _, x := range divisors(n) {
		m := n / x
		for _, y := range divisors(m) {
			z := m / y
			s := 2*x*y + 2*x*z + 4*x + y*z + 2*y + 2*z + 4
			if s < minStolen {
				minStolen = s
			}
			if s > maxStolen {
				maxStolen = s
			}
		}
	}
	return minStolen, maxStolen
}

// Embedded testcases from testcasesA.txt.
const testcaseData = `
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
16
17
18
19
20
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
`

func parseTestcases() ([]int64, error) {
	fields := strings.Fields(testcaseData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases found")
	}
	tests := make([]int64, 0, len(fields))
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("bad testcase at position %d: %v", i+1, err)
		}
		tests = append(tests, val)
	}
	return tests, nil
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]

	tests, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for idx, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		expMin, expMax := solve(n)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		var gMin, gMax int64
		if _, err := fmt.Fscan(strings.NewReader(got), &gMin, &gMax); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, got)
			os.Exit(1)
		}
		if gMin != expMin || gMax != expMax {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %d %d\ngot: %d %d\n", idx+1, input, expMin, expMax, gMin, gMax)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
