package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `3 583 1404
2 262 325
8 484 872
4 97 111
7 444 446
8 273 507
10 968 1008
1 23 43
9 10 912
7 703 811
1 541 736
8 962 1245
4 354 527
4 780 928
1 427 996
2 191 494
2 761 1017
7 520 1206
4 311 566
9 403 438
1 313 1001
5 536 937
6 934 1205
8 611 1185
8 860 1028
10 91 559
5 761 1181
7 14 465
6 201 1153
10 518 1061
5 214 659
7 943 1148
6 83 330
6 657 1014
8 654 918
10 707 1123
4 639 1204
5 194 847
9 773 1209
3 529 1177
6 315 1196
5 131 489
10 288 918
6 370 939
10 362 906
2 301 719
7 484 1071
2 736 1007
6 147 961
3 109 615
3 671 1098
6 483 873
10 423 1194
7 665 1008
6 364 466
10 381 634
10 575 1192
8 685 761
8 11 972
4 164 743
3 355 910
6 594 681
1 218 369
8 724 789
7 476 596
4 553 858
5 799 991
10 830 1027
10 803 823
1 283 722
1 844 1001
2 679 876
5 251 800
9 482 812
5 632 744
2 214 939
7 109 301
9 466 1179
9 601 636
6 84 394
5 326 759
7 109 273
3 61 1014
6 5 56
10 259 768
2 315 1142
8 341 769
10 49 160
2 440 468
3 477 721
3 339 634
1 165 1000
3 754 991
8 849 1048
7 390 1046
9 838 918
7 58 874
8 240 726
4 395 1207
6 415 717
5 675 1054
7 395 807
7 339 1012
8 580 1145
2 317 1124
1 103 233
1 198 277
1 564 1106`

type testCase struct {
	n int
	l int64
	r int64
}

func solveB(n int, l, r int64) (bool, []int64) {
	ans := make([]int64, n)
	for i := 1; i <= n; i++ {
		ii := int64(i)
		x := ((l-1)/ii + 1) * ii
		if x > r {
			return false, nil
		}
		ans[i-1] = x
	}
	return true, ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(testcasesRaw, "\n")
	var cases []testCase
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields got %d", i+1, len(fields))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", i+1, err)
		}
		l, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse l: %v", i+1, err)
		}
		r, err := strconv.ParseInt(fields[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: parse r: %v", i+1, err)
		}
		cases = append(cases, testCase{n: n, l: l, r: r})
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
		ok, arr := solveB(tc.n, tc.l, tc.r)
		var expected string
		if !ok {
			expected = "NO"
		} else {
			var sb strings.Builder
			sb.WriteString("YES\n")
			for i, v := range arr {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			expected = strings.TrimSpace(sb.String())
		}
		var inputSB strings.Builder
		inputSB.WriteString("1\n")
		inputSB.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.l, tc.r))
		got, err := runBinary(bin, inputSB.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
