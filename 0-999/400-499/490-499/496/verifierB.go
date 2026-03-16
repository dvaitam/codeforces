package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)
const testcasesBRaw = `
5 55445
16 6962029072507493
3 864
6 821291
17 89664440644888553
14 20829657086930
12 823570393426
3 973
15 813271664465951
10 0704366660
19 7592945067820195501
7 1870505
13 2462926480222
16 0808625911824345
9 487278209
6 805137
20 37825278081850116959
15 568512502205930
14 04609251607599
9 497620743
13 1511050269057
18 771086408115517028
10 0065282222
8 95076033
10 5235326759
5 69029
1 6
6 200580
2 01
19 9226066953258386126
19 5166376363769144858
1 9
20 73409561802606661797
6 257629
10 8155257803
9 295460486
2 64
13 3522195206967
3 168
18 222323082759451642
11 85826843656
16 8466122099813998
4 4926
3 001
12 751759473288
3 820
6 869376
9 092627906
3 965
8 87079148
16 8674238524927117
13 9681473142512
2 29
18 300689717855103374
9 307585114
19 6356234375469216587
8 55645641
16 4172532573166787
19 9368473580175636090
14 14352125030068
1 2
4 9931
15 308961823366706
7 6094095
12 572981411402
20 26395368818171877278
11 26461985373
12 760708673634
16 7237457928134810
5 50592
15 636732605835288
20 96297939703795830002
8 79318645
16 3258576757123126
15 835006316996830
5 87468
18 119703695410390985
5 74197
9 366853941
2 25
19 9718765601692201630
9 461545938
18 058802704525198614
19 1819270835183236397
12 245833752480
3 545
2 19
1 4
4 2438
18 571688611975243640
1 0
9 155840382
9 906649345
18 771652803760314122
15 058618374466702
17 53323422848023623
18 292384073062686368
17 71218010367465641
6 933008
3 737
`


type testCase struct {
	n int
	s string
}

func parseTestcases() ([]testCase, error) {
	var cases []testCase
	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("bad line: %s", line)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		s := fields[1]
		if len(s) != n {
			return nil, fmt.Errorf("len mismatch")
		}
		cases = append(cases, testCase{n: n, s: s})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func cmpStrings(a, b string) bool {
	i := 0
	for i < len(a) && a[i] == '0' {
		i++
	}
	a = a[i:]
	j := 0
	for j < len(b) && b[j] == '0' {
		j++
	}
	b = b[j:]
	if len(a) != len(b) {
		return len(a) < len(b)
	}
	return a < b
}

func solveCase(tc testCase) string {
	n := tc.n
	s := tc.s
	best := strings.Repeat("9", n)
	digits := make([]int, n)
	for i := 0; i < n; i++ {
		digits[i] = int(s[i] - '0')
	}
	for rot := 0; rot < n; rot++ {
		rotDigits := make([]int, n)
		for i := 0; i < n; i++ {
			rotDigits[i] = digits[(i-rot+n)%n]
		}
		for add := 0; add < 10; add++ {
			b := make([]byte, n)
			for i := 0; i < n; i++ {
				b[i] = byte((rotDigits[i]+add)%10) + '0'
			}
			cand := string(b)
			if cmpStrings(cand, best) {
				best = cand
			}
		}
	}
	return best
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d\n%s\n", tc.n, tc.s)
		expected := solveCase(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
