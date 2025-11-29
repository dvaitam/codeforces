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

const testcasesD = `4 77680 71335 17096 48492
10 62137 82016 76135 8590 79379 1727 61505 33996 72194 30716
4 94000 61640 70908 72043
8 52055 83765 19743 30400 83214 19875 68576 51111
1 88005
2 20894 99384
10 5610 39489 4066 35316 61966 77957 94219 50806 93604 55961
7 95438 75618 58279 17585 47911 12775 4705
3 64867 28442 33816
7 82138 39458 55202 66487 50578 75240 45996
9 76688 53423 76581 30461 44142 89390 3758 36660 79407
3 91570 42782 71012
10 74596 13643 93563 85921 27674 82967 75176 35009 37351 16311
2 63178 83725
8 11604 45101 8732 53802 19763 2639 38522 55988
7 15588 5794 79299 80550 99830 5892 49521
10 43380 72203 36580 66248 30928 4722 40591 950 10090 14173
10 70201 4114 25873 53471 38224 80017 34522 20475 90406 5564
6 41136 47213 18132 49519 49383 60348
9 50621 84397 78075 89258 73305 13446 81283 66459 35561
7 83139 94414 93799 31149 39466 57338 33849
9 39715 71885 44423 1503 54424 76019 41273 2630 49352
10 77230 82863 17469 7876 83040 82228 43576 61114 46259 89036
6 79807 92663 36561 96737 64162 2909
10 7941 88609 2787 48389 32917 82311 59811 39144 77684 78835
6 23256 47706 24282 40984 99368 48387
10 34622 39376 49437 13747 3530 74615 89620 96362 17227 40636
9 29170 85686 35309 31288 42965 24564 88877 57049 85138
2 13350 78740
6 43747 88466 29422 57467 22190 10480
6 97266 85229 28578 74503 59128 35470
4 15851 4446 69421 25011
6 75346 24065 36517 44580 84142 11213
10 45258 77266 17000 55219 38267 67949 35545 60913 45399 83125
7 38065 55022 74494 53678 4660 54169 20448
4 613 62569 81615 66867
7 73254 94014 29111 4238 97719 59858 98749
9 37887 71286 44706 29816 8923 77161 37621 15730 32053
1 4608
9 26015 56348 75624 6471 1724 63053 97689 15837 22522
9 39308 31334 86889 2603 68810 70387 54233 6981 80239
2 44743 16439
5 70909 62549 8045 46118 28943
4 16021 70077 15627 22454
4 35884 16842 985 63899
10 52467 6546 99156 35574 32541 35206 81016 69111 68116 55443
1 61985
6 239 7189 16630 6055 16337 6533
2 63285 4328
2 67557 65815
8 41407 20585 41234 9409 46046 50569 84794 51072
10 39863 47297 34726 25047 43095 56192 16222 16728 72815 459
7 10474 74284 23405 5633 48927 60402 79226
9 49844 83452 5695 81650 56571 6957 48815 82242 65033
6 55119 90982 54828 60405 2351 32125
4 70235 35403 91119 77321
2 55692 29419
7 17089 3684 42674 49053 73281 34366 15920
8 90514 16142 95884 86809 69493 49341 87500 14240
6 73906 69728 13519 77040 93925 644
8 18813 30931 50946 5810 69103 12034 73968 13017
7 23484 3080 44760 15898 3340 15087 88283
8 91273 37291 75902 39242 11639 4755 73879 67022
9 93739 31237 13993 72676 98170 13097 72526 8015 72107
6 73923 23650 10138 31737 23565 84670
4 59527 80739 91767 98770
7 33149 48167 78573 51990 45914 72961 54818
2 49196 65578
4 54112 97966 21054 54432
10 99107 76000 88325 67778 89859 63404 20472 84275 52568 19574
3 12566 65270 98064
8 91574 67799 58073 76871 94231 24395 17865 35045
4 19211 76760 67585 41269
4 90590 70523 38781 88001
7 78019 76623 76633 35020 28523 40270 3033
5 62849 50169 26298 22578 74698
6 31305 42212 63256 18825 54844 91431
8 91947 78514 26970 61350 76091 85512 73024 3641
8 94466 9488 52465 96154 6011 61264 30096 30774
2 28493 33316
4 24856 33900 18021 24536
10 92372 88494 4820 33426 22246 5903 41082 24023 55484 11924
2 15464 12144
5 38244 4737 46750 59293 76063
6 904 3853 43885 43448 57178 49766
8 10221 27539 84452 76692 97305 64225 51240 16441
9 41789 15622 36015 9995 87221 56698 14751 57466 69140
5 12703 69167 91844 49052 88924
6 99265 59038 38738 86889 88591 87815
5 14048 98855 44373 88190 74194
9 68903 14851 87411 64735 66676 46160 7803 94123 38592
10 97291 23854 84652 84630 95710 82693 19602 23477 48600 85902
8 16141 14186 73373 18533 43456 84540 94436 85144
10 55070 72736 39397 84869 24499 59994 63227 40954 23120 92468
2 14081 93827
3 98893 72580 71191
10 96947 51363 47085 13090 34846 35501 50252 7010 17896 5521
8 66155 35557 32402 91155 67467 46415 43630 52868
8 70977 9075 46255 65276 14699 19870 35455 77350
2 89325 14771`

const limitD = 10000005

var spf []int32
var primes []int32

func initSieve() {
	spf = make([]int32, limitD)
	primes = make([]int32, 0, 700000)
	for i := 2; i < limitD; i++ {
		if spf[i] == 0 {
			spf[i] = int32(i)
			primes = append(primes, int32(i))
		}
		for _, p := range primes {
			t := int32(i) * p
			if t >= limitD {
				break
			}
			spf[t] = p
			if p == spf[i] {
				break
			}
		}
	}
}

func referenceSolveD(arr []int) ([]int, []int) {
	a := make([]int, len(arr))
	b := make([]int, len(arr))
	for i, x := range arr {
		if x%2 == 0 {
			for x%2 == 0 {
				x /= 2
			}
			if x == 1 {
				a[i], b[i] = -1, -1
			} else {
				a[i], b[i] = 2, x
			}
			continue
		}
		var pa, pb int32 = -1, -1
		xx := int32(x)
		for xx > 1 {
			t := spf[int(xx)]
			if pa == -1 {
				pa = t
			} else if pb == -1 {
				pb = t
			}
			for xx%t == 0 {
				xx /= t
			}
			if pa != -1 && pb != -1 {
				break
			}
		}
		if pa == -1 || pb == -1 {
			a[i], b[i] = -1, -1
		} else {
			a[i], b[i] = int(pa), int(pb)
		}
	}
	return a, b
}

type testCaseD struct {
	n   int
	arr []int
}

func parseTestcasesD() ([]testCaseD, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesD))
	scan.Split(bufio.ScanWords)
	cases := make([]testCaseD, 0)
	for {
		if !scan.Scan() {
			break
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("parse n: %w", err)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("missing value %d for case %d", i+1, len(cases)+1)
			}
			v, err := strconv.Atoi(scan.Text())
			if err != nil {
				return nil, fmt.Errorf("parse value %d case %d: %w", i+1, len(cases)+1, err)
			}
			arr[i] = v
		}
		cases = append(cases, testCaseD{n: n, arr: arr})
	}
	if err := scan.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	cases, err := parseTestcasesD()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	initSieve()

	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		wantA, wantB := referenceSolveD(tc.arr)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 2 {
			fmt.Printf("case %d failed: expected two lines got %d\n", idx+1, len(lines))
			os.Exit(1)
		}
		partsA := strings.Fields(lines[0])
		partsB := strings.Fields(lines[1])
		if len(partsA) != tc.n || len(partsB) != tc.n {
			fmt.Printf("case %d failed: wrong output length\n", idx+1)
			os.Exit(1)
		}
		for i := 0; i < tc.n; i++ {
			ga, err := strconv.Atoi(partsA[i])
			if err != nil {
				fmt.Printf("case %d failed: bad integer %q\n", idx+1, partsA[i])
				os.Exit(1)
			}
			gb, err := strconv.Atoi(partsB[i])
			if err != nil {
				fmt.Printf("case %d failed: bad integer %q\n", idx+1, partsB[i])
				os.Exit(1)
			}
			if ga != wantA[i] || gb != wantB[i] {
				fmt.Printf("case %d failed at index %d\ninput:\n%sexpected:\n%v\n%v\ngot:\n%v\n%v\n", idx+1, i, input, wantA, wantB, partsA, partsB)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
