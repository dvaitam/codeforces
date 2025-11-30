package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `5 0487647 0488242 9489241 5781565 3877840
2 60 53
4 39 87 38 71
4 418583989 496593423 094711220 868483396
4 75159179 33041352 12309891 75199161
4 09 17 08 41
3 45 87 63
4 792302 584197 076984 642807
2 084237 085992
4 6109352 6109352 6109350 6109350
3 14278789 07547063 12066503
2 8 3
2 3442176104 1428512400
2 4855 0977
5 823694 845551 823600 894568 730428
2 65461 65475
4 71 60 71 22
5 11 06 18 77 36
2 349263 087317
5 30392 30365 21972 68757 38930
4 508249 694711 013204 508249
5 227586 809189 634896 769930 024894
4 74 46 70 34
4 0 7 9 2
4 0976701 0900992 1853671 9795194
3 4183067 4187510 0740899
3 1886 1269 1611
5 076 075 151 552 051
3 0310 9322 0315
4 035 034 852 118
4 622 583 645 290
5 12 12 96 11 12
4 066065735 066065712 418927626 066758257
2 34079703 34078926
5 73 64 01 73 76
5 14543910 14544851 32164482 38168644 14211669
5 845 637 847 074 845
2 3057608116 1904432943
2 7562566 7585232
5 667633 690603 250239 845670 886977
4 35400025 21639683 35400919 55875408
2 5138 5444
5 75304 10777 75377 13126 37916
5 2 2 4 5 1
4 8347897844 8912633408 0164195384 8341845357
4 2208855902 2281237933 3659927190 5774038277
4 42 35 60 43
4 34524059 25070030 34524105 62376254
3 660648 709166 660629
3 53 23 28
3 69 69 99
2 56007 56273
4 654 604 487 898
4 7 6 0 4
4 8856 4319 3985 8509
2 2554457696 0290725074
3 86 82 86
3 1889664440 1889488855 1889662082
5 708693 582357 708693 342619 708697
2 2716 4465
4 40 43 66 60
5 929450 929420 929450 131870 929450
2 624629 624629
4 022270808 625911824 022254487 278209280
4 37 78 32 37
2 185011695 568512502
3 5 3 4
5 9 2 6 9 9
4 97620 97361 05026 05787
5 08 40 11 70 08
4 0 8 2 3
4 7 6 4 5
3 675 269 290
5 200 800 992 606 203
3 838612 838615 663763
5 7691 4858 7691 3409 7691
2 606 661
5 7225762948 7225752578 3429546048 7225762948 7225763522
2 5206967116 8222323082
5 945164 585826 843656 846612 945198
2 9981 9981
3 1001575 1001575 1001732
2 202869376 202869376
3 2790619 2798707 4878674
3 8524 2711 8524
5 147314251 298300689 178551033 147343075 851149635
5 343 754 316 343 345
5 56417 17253 57316 78799 56484
5 5801 5803 6090 5835 5503
2 6 2
2 9317308961 3366706360
4 9 9 7 9
2 11402 26395
3 8818171 7727852 8818171
4 7357 6070 7363 7352
3 45792813 45792813 10250592
5 3673260 5835288 9962979 9703795 3000237
3 86 45 25
4 67571231 67835006 67571299 68302874
5 811970369 410390985 811919743 811970365 811970369`

type testCase struct {
	n    int
	strs []string
}

func run(bin, input string) (string, error) {
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

// solve mirrors 172A.go logic.
func solve(strs []string) int {
	if len(strs) == 0 {
		return 0
	}
	first := strs[0]
	commonLen := len(first)
	for i := 1; i < len(strs); i++ {
		s := strs[i]
		j := 0
		for j < commonLen && j < len(s) && s[j] == first[j] {
			j++
		}
		commonLen = j
		if commonLen == 0 {
			break
		}
	}
	return commonLen
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	var cases []testCase
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d parse n: %v", idx+1, err)
		}
		if len(fields) != n+1 {
			return nil, fmt.Errorf("line %d expected %d strings got %d", idx+1, n, len(fields)-1)
		}
		cases = append(cases, testCase{n: n, strs: fields[1:]})
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		want := solve(tc.strs)
		var input strings.Builder
		fmt.Fprintf(&input, "%d\n", tc.n)
		for _, s := range tc.strs {
			input.WriteString(s)
			input.WriteByte('\n')
		}
		gotStr, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", idx+1, want, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
