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

// solution logic from 1913A.go
func solveCase(s string) string {
	for i := 1; i < len(s); i++ {
		aStr := s[:i]
		bStr := s[i:]
		if aStr[0] == '0' || bStr[0] == '0' {
			continue
		}
		aVal, _ := strconv.Atoi(aStr)
		bVal, _ := strconv.Atoi(bStr)
		if bVal > aVal {
			return fmt.Sprintf("%d %d", aVal, bVal)
		}
	}
	return "-1"
}

const testcasesData = `
345
855
5927354
5519
313328
288554
487948
987671
51956
751
11948414
464426
997861
15879
662231
731035
777945
822965
5494
23
7429
99948750
80
415
538666
192199
47
554818
349564
70
13374
83257
55
95
61
2803
993052
48420
593597
14534
539364
36531
8395240
873624
55700
978187
104672
787272
766432
1324
15
628676
725353
764579
415492
939077
95827
538
75056325
289
119625
786831
55150
656
80
9331079
143
86963
849356
443574
385517
276417
284234
804
59829079
941675
40
60972989
8835
747780
9050
468840
885150
975719
200510
55873
729352
5565
313335
20
61475
212449
332403
77430
222126
613340
787105
206105
178240
4402876
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesData))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		input := fmt.Sprintf("1\n%s\n", line)
		want := solveCase(line)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != want {
			fmt.Printf("test %d failed\ninput: %s\nexpected: %s\ngot: %s\n", idx, line, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
