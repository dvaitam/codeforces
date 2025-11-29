package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `710 1525
666 831
873 1546
873 1966
390 1039
642 1642
894 968
633 774
241 1534
701 1290
233 1763
93 982
101 1659
649 855
455 796
711 1325
926 986
48 713
816 931
301 1036
271 1821
345 1358
478 640
997 1997
65 1747
723 842
794 1197
821 1564
880 1585
315 1687
572 659
892 1780
457 1900
433 983
811 921
176 842
416 1876
735 1112
119 956
82 1736
579 930
128 1245
746 1938
421 1476
889 1547
760 1636
462 862
928 1628
531 575
996 1045
528 577
688 766
79 1212
903 1510
920 1312
766 1795
591 1435
226 1308
584 845
662 1949
125 1016
611 1029
247 1059
217 1399
106 622
7 273
963 1586
36 1344
925 1949
584 828
352 1680
621 962
414 1167
625 655
380 1963
957 1750
996 1993
93 1357
580 1444
994 1076
48 1459
266 1409
95 1520
692 1724
361 1669
334 1277
912 1007
709 1771
94 1358
471 1495
478 581
633 1382
456 1067
910 1118
882 1155
484 1660
11 182
746 1857
914 1891
849 1437
910 1505
232 449
777 980
189 1126`

type testCaseB struct {
	h, l int
}

func parseTestsB() ([]testCaseB, error) {
	fields := strings.Fields(testcases)
	if len(fields)%2 != 0 {
		return nil, fmt.Errorf("expected even number of ints, got %d", len(fields))
	}
	cases := make([]testCaseB, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		h, err := strconv.Atoi(fields[i])
		if err != nil {
			return nil, fmt.Errorf("parse h: %w", err)
		}
		l, err := strconv.Atoi(fields[i+1])
		if err != nil {
			return nil, fmt.Errorf("parse l: %w", err)
		}
		cases[i/2] = testCaseB{h: h, l: l}
	}
	return cases, nil
}

func solveB(tc testCaseB) float64 {
	h := float64(tc.h)
	l := float64(tc.l)
	return (l*l - h*h) / (2 * h)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestsB()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.h, tc.l)
		expect := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Printf("case %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		diff := math.Abs(val - expect)
		if diff > 1e-6*math.Max(1.0, math.Abs(expect)) {
			fmt.Printf("case %d failed: expected %.10f got %.10f\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
