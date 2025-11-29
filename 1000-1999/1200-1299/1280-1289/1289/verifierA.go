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

// Embedded testcases from testcasesA.txt.
const testcasesRaw = `729 -212
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

type testCase struct {
	a int
	b int
}

// referenceSolution embeds 1289A.go logic (sum of two integers).
func referenceSolution(a, b int) int {
	return a + b
}

func parseTestcases() []testCase {
	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	tests := make([]testCase, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			panic("invalid testcase line")
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			panic("invalid integers in testcase")
		}
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := parseTestcases()
	for idx, tc := range tests {
		expected := referenceSolution(tc.a, tc.b)
		input := fmt.Sprintf("%d %d\n", tc.a, tc.b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(out)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx+1, out)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
