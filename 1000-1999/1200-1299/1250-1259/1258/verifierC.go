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

const testcasesRaw = `924 593
306 578
740 786
436 332
778 100
785 835
761 496
426 945
821 709
157 119
909 868
449 705
399 645
696 300
591 956
156 796
132 548
984 823
422 964
793 173
568 932
205 19
849 726
758 197
748 433
466 781
323 96
884 977
620 160
438 989
990 467
423 692
126 286
246 571
520 276
526 533
639 814
15 32
91 448
371 161
460 380
531 28
609 321
780 419
162 976
211 185
755 722
82 995
870 432
226 635
445 911
10 658
38 211
251 457
40 624
899 743
925 437
265 24
433 786
151 65
129 521
158 736
975 777
339 485
209 449
17 804
96 802
212 817
292 295
370 682
551 702
960 733
968 530
965 692
809 155
884 761
922 857
105 521
939 797
240 333
803 488
903 449
484 332
102 843
370 35
836 747
157 738
849 158
748 918
755 461
303 592
779 377
424 312
894 770
417 727
577 38
151 112
18 337
842 154
635 137`

type testCase struct {
	a int
	b int
}

func parseTestcases(raw string) ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(raw))
	scan.Split(bufio.ScanWords)
	tests := []testCase{}
	for scan.Scan() {
		aStr := scan.Text()
		if !scan.Scan() {
			return nil, fmt.Errorf("odd number of tokens")
		}
		bStr := scan.Text()
		a, err1 := strconv.Atoi(aStr)
		b, err2 := strconv.Atoi(bStr)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("invalid integers")
		}
		tests = append(tests, testCase{a: a, b: b})
	}
	return tests, scan.Err()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	tests, err := parseTestcases(testcasesRaw)
	if err != nil {
		fmt.Println("invalid test data:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		expected := max(tc.a, tc.b)
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d %d\n", tc.a, tc.b))
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		got, err := strconv.ParseInt(outStr, 10, 64)
		if err != nil {
			fmt.Printf("Test %d: invalid output %q\n", i+1, outStr)
			os.Exit(1)
		}
		if got != int64(expected) {
			fmt.Printf("Test %d failed: expected %d got %d\n", i+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
