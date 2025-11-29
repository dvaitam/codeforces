package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	a int
	b int
}

const testcasesRaw = `100
729 -212
552 823
-139 -918
-470 977
47 -5
-171 880
605 699
-379 982
-24 -267
194 826
859 -553
33 -715
-423 -714
547 -806
266 637
-487 863
90 444
659 232
847 -700
-365 -798
494 -849
840 741
400 -324
-34 146
-794 -276
-111 -353
251 311
869 -582
979 131
-24 -94
772 67
-467 -873
648 881
123 875
-972 -809
473 720
-184 454
689 607
368 280
-998 253
10 695
776 -318
-501 495
-334 441
782 -872
-609 878
162 -546
-512 645
981 -709
644 112
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
254 344
-468 -25
-859 -817
390 551
-734 795
-694 891
-921 725
-836 839
432 890
698 107
399 -199
715 444
74 -436
68 662
-518 739
-560 833
391 207
690 945
-142 187
-437 -78
8 352
313 434
877 624
-269 -832
-336 254
-764 -4
202 290
-314 730
-611 -503
-967 498`

func parseTestcases(raw string) []testCase {
	fields := strings.Fields(raw)
	if len(fields) < 1 {
		panic("no testcase data")
	}
	t, err := strconv.Atoi(fields[0])
	if err != nil {
		panic(fmt.Sprintf("bad testcase count: %v", err))
	}
	fields = fields[1:]
	if len(fields)%2 != 0 {
		panic("odd number of values in testcases")
	}
	if len(fields)/2 != t {
		panic(fmt.Sprintf("count mismatch: header %d pairs %d", t, len(fields)/2))
	}
	res := make([]testCase, 0, t)
	for i := 0; i < len(fields); i += 2 {
		a, err := strconv.Atoi(fields[i])
		if err != nil {
			panic(fmt.Sprintf("bad a at pair %d: %v", i/2+1, err))
		}
		b, err := strconv.Atoi(fields[i+1])
		if err != nil {
			panic(fmt.Sprintf("bad b at pair %d: %v", i/2+1, err))
		}
		res = append(res, testCase{a: a, b: b})
	}
	return res
}

func expected(tc testCase) int {
	return tc.a + tc.b
}

func runCandidate(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases(testcasesRaw)

	var input strings.Builder
	input.WriteString(strconv.Itoa(len(tests)))
	input.WriteByte('\n')
	for _, tc := range tests {
		input.WriteString(fmt.Sprintf("%d %d\n", tc.a, tc.b))
	}
	want := make([]string, len(tests))
	for i, tc := range tests {
		want[i] = strconv.Itoa(expected(tc))
	}
	gotStr, err := runCandidate(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	gotLines := strings.Fields(gotStr)
	if len(gotLines) != len(want) {
		fmt.Fprintf(os.Stderr, "expected %d lines got %d\n", len(want), len(gotLines))
		os.Exit(1)
	}
	for i := range want {
		if strings.TrimSpace(gotLines[i]) != want[i] {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, want[i], gotLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
