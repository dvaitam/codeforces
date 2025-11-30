package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesC = `4 25 110251 10613 67874 134028 127384 106152 79512 124938
3 19 3579 8269 2282 4618 2290 1554
5 26 131338 279217 315571 77051 162607 51781 38663 173119 247539 293504
1 12 223 162
5 21 6701 18106 15631 14507 28355 17084 8536 2041 26380 17980
1 3 4 1
5 16 1365 1000 1333 258 783 2325 909 978 584 2225
4 3 1 3 5 4 1 3 5 3
1 18 2726 6673
5 7 19 29 6 25 21 16 19 12 13 12
1 20 8522 15615
1 3 2 5
1 27 42074 470902
5 22 25644 46222 34379 18064 34197 15434 14104 44531 38654 27488
5 9 58 64 85 83 46 11 42 79 15 63
5 21 10987 27692 6239 7964 532 23970 8882 3839 23113 7225
3 26 89382 174345 223416 32607 52748 76734
2 2 3 1 3 1
1 1 1 2
5 27 301966 62755 205107 47990 194060 437175 514074 60842 19080 317459
1 7 12 8
4 7 4 2 28 7 17 5 15 5
3 12 224 93 32 258 240 21
5 4 7 4 5 6 8 3 4 1 3 8
3 17 2054 961 3624 1433 109 3864
4 29 1193516 1066612 653150 749002 815042 526244 321730 1175664
1 15 1519 162
3 24 5990 71346 36816 17674 31475 99899
4 12 313 148 345 184 303 325 318 68
3 13 425 83 2 609 197 343
2 8 15 41 29 25
5 28 434537 33078 421808 735922 595038 438543 809696 694363 743447 49053
2 15 131 531 1437 323
4 17 3992 1 319 4052 2671 2557 3825 409
4 7 6 9 1 26 27 21 1 14
1 23 310 69252
5 4 4 2 4 5 5 3 2 8 7 2
1 9 58 15
3 5 11 9 11 13 6 2
2 9 3 6 6 27
3 18 2579 3006 4649 345 6137 5746
5 21 16204 23343 21108 15030 20977 14269 12205 28555 17627 5843
2 13 602 299 10 142
2 9 43 44 48 12
3 25 162671 9350 10802 70680 42956 39169
5 10 93 102 141 34 76 30 123 62 13 79
2 28 548458 764057 74306 317374
4 27 172240 156881 217434 56965 52115 294037 476399 252273
4 11 216 205 208 88 32 123 30 180
4 14 39 310 344 753 704 916 160 942
2 21 18498 12307 26415 20935
1 3 1 2
2 2 2 1 1 2
5 17 2375 3675 4003 1780 3466 686 3018 1804 2138 1365
4 7 23 8 5 2 34 29 13 8
4 13 263 213 44 222 150 108 203 470
4 12 280 78 54 306 250 76 289 208
4 29 1093096 1038942 676249 1045307 1045738 1331828 423607 1138354
5 30 917601 40705 1426988 1334877 1349832 148765 622296 1077441 653915 1589663
5 10 121 17 22 133 11 17 58 34 11 77
1 25 117595 86664
2 26 78025 241558 194680 264766
4 29 1111088 1053670 70503 1203498 190101 1086322 1258590 160135
4 29 432234 607470 1122850 1255386 876276 1011276 814889 1273873
5 8 55 52 2 43 1 48 12 20 33 37
3 11 17 127 221 68 212 78
4 13 393 64 168 131 245 294 343 57
1 16 1712 578
4 29 1262416 171169 317454 739817 862224 73778 1282977 977858
4 15 97 208 965 1594 311 42 67 1226
5 5 11 6 2 12 9 11 6 4 7 13
4 4 1 8 6 2 5 3 7 5
1 17 1550 313
4 15 761 1550 391 933 731 1295 155 92
1 16 1047 110
5 22 37307 37474 14155 15051 6127 41120 32927 45778 34322 27533
5 10 30 38 110 109 22 27 107 17 26 107
2 24 4025 103723 58553 56504
4 1 2 1 2 1 2 1 1 2
1 12 179 92
1 27 510827 120845
3 3 5 2 2 1 2 1
1 10 95 7
5 8 55 10 12 30 8 31 23 46 17 9
1 7 24 22
4 10 76 142 84 48 21 27 137 79
2 13 151 129 229 324
5 8 16 49 12 19 24 27 43 3 9 39
1 13 80 75
2 14 307 564 427 759
2 19 6918 4883 10437 5810
1 8 29 41
3 21 17342 1896 12335 13390 277 13672
3 15 418 762 601 965 187 380
1 9 15 72
5 23 20170 58480 52269 24299 55274 56586 22904 32501 59435 44616
5 5 6 8 11 13 2 8 13 4 5 1
4 20 15140 255 7168 9784 3751 9869 5120 13901
4 3 4 2 5 4 3 1 1 3
1 1 2 1`

// Embedded solution logic from 232C.go.
func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func endpointDistances(k int, u uint64, f []uint64, distSE []int64) (int64, int64) {
	ks := make([]int, 0, k+1)
	us := make([]uint64, 0, k+1)
	for k > 1 {
		ks = append(ks, k)
		us = append(us, u)
		if u <= f[k-1] {
			k = k - 1
		} else {
			u = u - f[k-1]
			k = k - 2
		}
	}
	var ds, de int64
	if k == 0 {
		ds, de = 0, 0
	} else if k == 1 {
		if u == 1 {
			ds, de = 0, 1
		} else {
			ds, de = 1, 0
		}
	}
	for i := len(ks) - 1; i >= 0; i-- {
		ck := ks[i]
		cu := us[i]
		if cu <= f[ck-1] {
			a, b := ds, de
			ds = a
			de = minInt64(a, b) + 1 + distSE[ck-2]
		} else {
			a2, b2 := ds, de
			ds = a2 + 1
			de = b2
		}
	}
	return ds, de
}

func solveCase(t, n int, queries [][2]uint64) string {
	maxU := uint64(10000000000000000) // 1e16
	f := make([]uint64, n+1)
	f[0], f[1] = 1, 2
	for i := 2; i <= n; i++ {
		v := f[i-1] + f[i-2]
		if v > maxU {
			f[i] = maxU + 1
		} else {
			f[i] = v
		}
	}

	distSE := make([]int64, n+1)
	if n >= 0 {
		distSE[0] = 0
	}
	if n >= 1 {
		distSE[1] = 1
	}
	for i := 2; i <= n; i++ {
		a := 1 + distSE[i-2]
		b := distSE[i-1] + 1 + distSE[i-2]
		if a < b {
			distSE[i] = a
		} else {
			distSE[i] = b
		}
	}

	var sb strings.Builder
	for idx, q := range queries {
		u, v := q[0], q[1]
		if u == v {
			sb.WriteString("0")
		} else {
			kk := n
			var ans int64
			for {
				if kk == 1 {
					ans = 1
					break
				}
				if u <= f[kk-1] && v <= f[kk-1] {
					kk--
					continue
				}
				if u > f[kk-1] && v > f[kk-1] {
					u -= f[kk-1]
					v -= f[kk-1]
					kk -= 2
					continue
				}
				if u > f[kk-1] {
					u, v = v, u
				}
				du1, du2 := endpointDistances(kk-1, u, f, distSE)
				v2 := v - f[kk-1]
				dv1, _ := endpointDistances(kk-2, v2, f, distSE)
				ans = minInt64(minInt64(du1, du2)+1+dv1, 1<<60)
				break
			}
			sb.WriteString(strconv.FormatInt(ans, 10))
		}
		if idx+1 != len(queries) {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

type testCase struct {
	t       int
	n       int
	queries [][2]uint64
}

func parseCases() ([]testCase, error) {
	data := strings.TrimSpace(testcasesC)
	if data == "" {
		return nil, fmt.Errorf("no testcases provided")
	}
	lines := strings.Split(data, "\n")
	res := make([]testCase, 0, len(lines))
	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d malformed", i+1)
		}
		t, err1 := strconv.Atoi(fields[0])
		n, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("line %d parse error", i+1)
		}
		if len(fields) != 2+2*t {
			return nil, fmt.Errorf("line %d expected %d pairs got %d numbers", i+1, t, len(fields)-2)
		}
		queries := make([][2]uint64, t)
		idx := 2
		for j := 0; j < t; j++ {
			u, errU := strconv.ParseUint(fields[idx], 10, 64)
			v, errV := strconv.ParseUint(fields[idx+1], 10, 64)
			if errU != nil || errV != nil {
				return nil, fmt.Errorf("line %d query %d parse error", i+1, j+1)
			}
			queries[j] = [2]uint64{u, v}
			idx += 2
		}
		res = append(res, testCase{t: t, n: n, queries: queries})
	}
	return res, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.t))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.FormatUint(q[0], 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatUint(q[1], 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: verifierC /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expect := solveCase(tc.t, tc.n, tc.queries)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n got: %s\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
