package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcaseData = `
9 220 -266 1 -67 -368 -975 963 -550 137
3 353 800 13
8 118 -358 442 736 -839 -470 -719 237
7 440 -609 682 -352 608 -403 -214
1 -573
1 -355
4 -298 767 -97 368
4 -468 -296 375 -668
5 -966 -271 174 104 -881
3 -277 -955 5
1 -950
4 -908 -975 -538 936
6 -864 681 -873 -294 990 359
7 -722 827 888 -557 -81 -109 -710
6 -362 -638 329 -327 491 608
7 -217 -981 -162 824 -460 93 88
8 554 -915 156 887 -749 -163 953 -201
3 -994 25 -717
10 779 355 52 720 500 432 -697 -837 -325 -514
3 -495 937 -955
4 693 2 627 -321
7 -209 -416 -582 833 -835 843 600
1 -245
10 -166 -157 -199 -266 -380 787 134 982 432 -59
2 -405 20
4 457 193 -386 359
5 -469 901 334 -815 -836
4 805 392 794 557
10 211 -109 -226 886 688 -748 907 -566 -590 687
9 -357 -837 283 662 691 -197 942 654 -34
4 20 513 742 60
4 265 -212 590 844
6 522 -251 48 -65 640 -338
2 -328 -624
2 -554 315
10 -824 -364 -6 -223 968 460 -216 264 608 -253
9 401 -654 346 512 871 970 80 289 135
4 -666 956 158 451
5 -202 952 140 772 8
10 172 -564 -96 756 -57 404 -504 -399 170 -779
10 -268 -284 -608 834 602 34 -835 -689 298 710
8 178 398 395 477 365 478 190 -763
3 -621 -533 -516
6 -355 -268 -109 -711 455 142
1 184
9 586 896 886 63 963 -544 -80 676 -567
7 762 -353 -908 657 610 -580 -92
3 -383 814 708
3 -711 -849 -481
7 -890 765 -436 -79 -101 990 758
7 -21 240 -209 727 890 764 960
10 302 907 -72 -821 -572 652 224 586 -385 -122
10 -203 420 -423 -815 963 -764 759 582 598 -260
8 46 174 -220 -46 362 933 98 -552
8 -369 233 467 -768 -153 298 922 361
9 -230 926 37 145 23 -627 765 -786 294
9 -874 161 531 590 -485 -129 -559 220 63
8 579 14 309 -265 -388 166 4 849
5 -796 -253 197 -585 998
6 -494 199 756 -800 519 300
2 -376 -183
2 96 491
2 385 582
7 318 348 785 -321 778 418 677
8 836 -631 -827 182 -290 883 911 -993
6 514 -154 281 -592 917 -115
1 870
5 -845 150 512 48 291
4 171 -60 813 -630
5 -885 718 248 344 114
8 -681 88 22 -45 -987 -906 984 -66
2 277 201
6 486 362 -634 439 -199 123
7 263 -503 -76 270 888 841 -109
2 -920 -191
6 32 968 206 -504 339 299
1 -951
1 700
4 -169 880 270 490
8 241 -280 -292 533 -6 164 927 -444
3 256 447 -226
1 -983
3 -403 214 -835
10 -728 925 -879 17 -378 1000 -744 389 -980 965
7 -518 145 -100 949 215 -47 444
9 555 465 -869 -218 -616 -357 201 -12 254
5 -697 766 823 356 179
7 780 559 100 533 -180 986 858
3 876 339 751
7 -451 672 -421 -127 781 -180 -213
4 643 978 397 812
2 969 -154
6 -311 394 755 -829 -279 -189
3 -890 445 11
5 240 -909 965 -304 726
8 980 179 -778 -325 539 -559 -115 778
2 -550 144
7 100 232 170 -438 25 583 -73
3 -286 919 717
8 -53 -716 162 974 289 -576 -287 628
10 -997 29 485 -799 -939 570 143 868 331 744
6 654 577 -828 777 264 -932
`

type testCase struct {
	input    string
	expected string
}

func solve(arr []int64) int64 {
	var sum int64
	for _, v := range arr {
		sum += v
	}
	return sum
}

func loadCases() ([]testCase, error) {
	lines := strings.Split(testcaseData, "\n")
	cases := []testCase{}
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: bad n: %w", idx+1, err)
		}
		if len(fields)-1 != n {
			return nil, fmt.Errorf("line %d: expected %d numbers, got %d", idx+1, n, len(fields)-1)
		}
		arr := make([]int64, n)
		valStr := make([]string, n)
		for i := 0; i < n; i++ {
			v, err := strconv.ParseInt(fields[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("line %d: bad value: %w", idx+1, err)
			}
			arr[i] = v
			valStr[i] = fields[i+1]
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		sb.WriteString(strings.Join(valStr, " "))
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			input:    sb.String(),
			expected: strconv.FormatInt(solve(arr), 10),
		})
	}
	return cases, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD1 /path/to/binary")
		os.Exit(1)
	}
	cases, err := loadCases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load cases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		got, err := run(os.Args[1], tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", idx+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
