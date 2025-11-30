package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `
4 3 005 411 998 679 807
2 3 978 207 151
2 2 37 78 89
2 3 546 659 051
3 2 86 44 92 51
2 3 462 914 865
2 4 8168 5054 2357
2 2 22 14 18
1 4 5929 6222
2 4 7065 3794 7383
3 3 359 774 688 623
2 3 075 818 141
2 3 782 613 750
4 2 68 53 61 53 05
1 4 5220 4727
4 4 0104 3289 8614 3410 3697
1 2 79 80
2 4 2460 9539 6218
3 2 88 88 80 67
1 4 6540 5153
1 4 5205 8527
4 4 2217 0430 3054 8687 4034
3 2 54 15 66 76
3 2 77 58 41 61
4 4 2845 1154 4796 2757 0596
3 4 0165 8202 9702 1355
4 4 0927 5571 9285 6543 1027
4 4 1447 3947 3121 7271 5518
3 3 225 831 323 705
3 3 829 114 678 669
1 2 51 77
3 2 89 22 68 01
2 3 225 358 414
2 4 8424 9818 2992
2 4 9959 0010 9439
4 4 0784 4736 4710 2767 7359
2 3 556 258 815
2 3 147 321 046
4 2 25 95 32 78 77
3 3 016 873 395 004
2 4 2111 7432 7349
2 4 9546 5096 8230
4 4 7909 8813 2008 8702 8378
2 3 299 216 093
2 2 85 49 05
4 3 351 041 177 075 190
4 2 60 52 67 59 17
4 4 3170 3508 6349 1673 8990
2 4 3130 1311 0471
2 3 484 547 324
4 2 25 40 18 47 78
4 3 155 096 611 998 932
3 4 0774 2154 0761 3514
2 3 963 632 486
3 2 66 91 68 46
1 3 209 203
1 2 59 75
4 4 5323 9362 3483 6196 2679
3 4 9137 2966 1727 8011
1 2 80 56
2 2 26 00 19
2 2 76 45 57
2 2 45 52 61
3 2 33 37 93 74
1 4 9620 0101
2 2 98 69 92
1 3 160 184
1 2 49 48
2 2 12 36 40
4 2 77 86 56 21 48
2 2 54 24 94
4 4 2050 1392 6055 8473 5419
3 2 08 92 58 42
2 2 08 70 06
3 3 041 218 637 117
3 2 90 08 76 47
2 4 7701 8674 2921
4 4 2638 0748 7230 0474 0122
4 4 8361 6677 7336 3550 6598
3 4 3199 3159 2326 6719
1 4 8871 1234
2 2 98 14 06
3 2 56 90 60 60
4 4 7446 2512 3102 4396 3037
2 4 0816 7099 1509
4 3 095 857 137 863 280
3 3 502 352 737 198
4 2 81 99 83 36 18
2 4 1394 0400 0079
2 2 16 57 44
4 2 81 25 09 87 15
4 4 2159 3230 5214 6404 0604
4 2 62 23 16 09 30
2 3 073 374 852
2 2 45 69 71
1 3 329 223
3 3 977 314 971 496
2 3 329 930 860
`

type testCase struct {
	n      int
	m      int
	src    []string
	target string
}

type info struct{ l, r, id int }

func solveCase(tc testCase) (string, []info) {
	map2 := make(map[string]info)
	map3 := make(map[string]info)
	for i, s := range tc.src {
		idx := i + 1
		for j := 0; j+1 < tc.m; j++ {
			key := s[j : j+2]
			map2[key] = info{j + 1, j + 2, idx}
		}
		for j := 0; j+2 < tc.m; j++ {
			key := s[j : j+3]
			map3[key] = info{j + 1, j + 3, idx}
		}
	}

	s := tc.target
	m := tc.m
	valid := make([]bool, m)
	prev := make([]int, m)
	seg := make([]info, m)
	for i := range prev {
		prev[i] = -1
	}
	if m >= 2 {
		if v, ok := map2[s[0:2]]; ok {
			valid[1] = true
			seg[1] = v
		}
	}
	if m >= 3 {
		if v, ok := map3[s[0:3]]; ok {
			valid[2] = true
			seg[2] = v
		}
	}
	for i := 3; i < m; i++ {
		key2 := s[i-1 : i+1]
		if v, ok := map2[key2]; ok && valid[i-2] {
			valid[i] = true
			prev[i] = i - 2
			seg[i] = v
		} else {
			key3 := s[i-2 : i+1]
			if v3, ok3 := map3[key3]; ok3 && valid[i-3] {
				valid[i] = true
				prev[i] = i - 3
				seg[i] = v3
			}
		}
	}
	if m == 0 || !valid[m-1] {
		return "-1", nil
	}
	var res []info
	for i := m - 1; i >= 0; {
		res = append(res, seg[i])
		if prev[i] < 0 {
			break
		}
		i = prev[i]
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return fmt.Sprintf("%d", len(res)), res
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: not enough fields", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		expected := 2 + n + 1
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d fields got %d", idx+1, expected, len(fields))
		}
		src := make([]string, n)
		for i := 0; i < n; i++ {
			src[i] = fields[2+i]
			if len(src[i]) != m {
				return nil, fmt.Errorf("line %d: source %d length %d != m %d", idx+1, i, len(src[i]), m)
			}
		}
		target := fields[len(fields)-1]
		if len(target) != m {
			return nil, fmt.Errorf("line %d: target length %d != m %d", idx+1, len(target), m)
		}
		cases = append(cases, testCase{n: n, m: m, src: src, target: target})
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
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
		firstLine, segments := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, s := range tc.src {
			sb.WriteString(s)
			sb.WriteByte('\n')
		}
		sb.WriteString(tc.target)
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := firstLine
		if firstLine != "-1" {
			var segBuf strings.Builder
			for _, seg := range segments {
				segBuf.WriteString(fmt.Sprintf("%d %d %d\n", seg.l, seg.r, seg.id))
			}
			expected = strings.TrimSpace(firstLine + "\n" + strings.TrimSpace(segBuf.String()))
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
