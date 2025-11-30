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
-83 -814
-836 -345
792 40
910 2
-777 -383
128 -404
447 -745
121 -319
668 888
106 -584
973 637
235 120
203 -411
-89 -813
221 634
-212 -351
178 -505
-406 -624
-613 682
-618 -933
-424 396
327 36
-105 -653
572 46
-882 -596
-63 -356
675 -640
929 -481
-12 105
-113 -632
422 -390
40 566
453 -272
206 691
1000 -237
192 120
-825 -657
-854 964
330 -602
-984 139
-562 988
-498 324
625 996
-701 570
-616 -941
802 952
754 1000
-961 109
277 969
-984 957
-907 -72
-756 651
-957 938
352 63
-548 -63
755 -642
147 738
875 -485
388 -651
-592 1000
958 -780
128 610
-583 704
-52 146
-273 160
920 46
-406 437
-122 182
-882 202
380 970
522 988
-848 568
9 -794
-273 155
-318 -971
606 -442
762 -695
100 660
-611 1000
1000 396
520 -774
-947 1000
615 664
200 431
21 278
-927 -896
-593 -552
-24 975
40 443
923 -923
258 -548
421 542
-947 509
152 530
-464 -504
-100 923
-726 -711
-458 710
-364 -917
-643 1000
400 400
624 -169
-276 584
`

type testCase struct {
	a, b  int64
	input string
	want  string
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func loadCases() ([]testCase, error) {
	scan := strings.Split(testcaseData, "\n")
	cases := []testCase{}
	for idx, line := range scan {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: expected 2 numbers, got %d", idx+1, len(parts))
		}
		a, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad a: %w", idx+1, err)
		}
		b, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: bad b: %w", idx+1, err)
		}
		var sb strings.Builder
		sb.WriteString(strconv.FormatInt(a, 10))
		sb.WriteByte(' ')
		sb.WriteString(strconv.FormatInt(b, 10))
		sb.WriteByte('\n')
		cases = append(cases, testCase{
			a:     a,
			b:     b,
			input: sb.String(),
			want:  strconv.FormatInt(gcd(a, b), 10),
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
		fmt.Println("usage: verifierB /path/to/binary")
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
		if got != tc.want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, tc.want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
