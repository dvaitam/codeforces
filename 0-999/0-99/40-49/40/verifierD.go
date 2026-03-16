package main

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
2
13
98259791
7483378876232860129
40
796669725
273
646869589
93504925899139
411771516
4661
99
69
84830401198
36
9420555266
869819846346484
506950699205755947
9
5
749952525
944610924834
52611955
72153974
10835924
19592648475
46960623079868
7848531
4130083690071284308
609152487053318123
207960434
88607500201017018875
515669
5435612806
920
798609605756
70338491636205
841718615981
723608191620510174
4109883181808592132
79645965
616836269987262217
8792243298538469
9434047632953572679
7980716
73
13433422
4205261111440579500
56671397628
14161784185
744159881785
49
225711
259684277421128896
144602078438
567815712491
9123
6629926388292345407
65894784709301
72873383
87
9422
190059381780555
219015313637
10613267718
37407767270455
859630345278323
2
6820219277200675600
60670331
7325915190447
7384055510
19901010280703
53757506496142618
8586268525673
754289670220711
39090775501
6159417
330
29013
3388689583
291019013411763
17569714976318074155
7185619
33505367143523397
653967679040210248
7033476057342174
71327465264778
354407554380028280
3755
72304650981744376
500602932
85824020070778
86582421281657744748
95280
73799405998172
416150877
4526650
305596840
64617
`

var rawTestcases = strings.Fields(testcasesRaw)

// Embedded correct solver for 40/D
// Returns (yes bool, year int, values []string)
func solve40D(s string) (bool, int, []string) {
	A := new(big.Int)
	A.SetString(s, 10)

	base := big.NewInt(12)
	tmp := new(big.Int).Set(A)
	digits := make([]int, 0)

	for tmp.Sign() > 0 {
		q := new(big.Int)
		r := new(big.Int)
		q.QuoRem(tmp, base, r)
		digits = append(digits, int(r.Int64()))
		tmp = q
	}

	if len(digits) == 0 {
		digits = append(digits, 0)
	}

	nzPos := make([]int, 0, 2)
	nzVal := make([]int, 0, 2)
	for i, d := range digits {
		if d != 0 {
			nzPos = append(nzPos, i)
			nzVal = append(nzVal, d)
		}
	}

	yes := false
	year := 0

	if len(nzPos) == 1 && nzVal[0] == 2 {
		year = 2*nzPos[0] + 1
		yes = true
	} else if len(nzPos) == 2 && nzVal[0] == 1 && nzVal[1] == 1 {
		year = nzPos[0] + nzPos[1] + 1
		yes = true
	}

	if !yes {
		return false, 0, nil
	}

	sExp := year - 1
	pow12 := make([]*big.Int, sExp+1)
	pow12[0] = big.NewInt(1)
	for i := 1; i <= sExp; i++ {
		pow12[i] = new(big.Int).Mul(pow12[i-1], base)
	}

	values := make([]string, 0)
	for q := sExp / 2; q >= 0; q-- {
		p := sExp - q
		v := new(big.Int).Add(pow12[p], pow12[q])
		if v.Cmp(A) != 0 {
			values = append(values, v.String())
		}
	}

	return true, year, values
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// Validate candidate output for multi-answer problem
func validate(input string, candOutput string) error {
	yes, year, values := solve40D(input)

	lines := strings.Split(strings.TrimSpace(candOutput), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	if !yes {
		if len(lines) != 1 || lines[0] != "NO" {
			return fmt.Errorf("expected NO, got: %s", candOutput)
		}
		return nil
	}

	// Expected: YES, then numGroups, then year, then numValues, then values
	if len(lines) < 1 || lines[0] != "YES" {
		return fmt.Errorf("expected YES, got first line: %s", lines[0])
	}

	if len(lines) < 2 {
		return fmt.Errorf("missing number of groups line")
	}
	numGroups, err := strconv.Atoi(lines[1])
	if err != nil || numGroups != 1 {
		return fmt.Errorf("expected 1 group, got: %s", lines[1])
	}

	if len(lines) < 3 {
		return fmt.Errorf("missing year line")
	}
	candYear, err := strconv.Atoi(lines[2])
	if err != nil || candYear != year {
		return fmt.Errorf("expected year %d, got: %s", year, lines[2])
	}

	if len(lines) < 4 {
		return fmt.Errorf("missing numValues line")
	}
	numValues, err := strconv.Atoi(lines[3])
	if err != nil || numValues != len(values) {
		return fmt.Errorf("expected %d values, got: %s", len(values), lines[3])
	}

	// Check that candidate lists the same set of values
	expectedSet := make(map[string]bool)
	for _, v := range values {
		expectedSet[v] = true
	}

	if len(lines) != 4+numValues {
		return fmt.Errorf("expected %d value lines, got %d total lines", numValues, len(lines)-4)
	}

	candSet := make(map[string]bool)
	for i := 4; i < len(lines); i++ {
		candSet[lines[i]] = true
	}

	for v := range expectedSet {
		if !candSet[v] {
			return fmt.Errorf("missing value %s in candidate output", v)
		}
	}
	for v := range candSet {
		if !expectedSet[v] {
			return fmt.Errorf("unexpected value %s in candidate output", v)
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, tc := range rawTestcases {
		input := tc + "\n"
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := validate(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput: %s\ngot:\n%s\n", idx+1, err, tc, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
