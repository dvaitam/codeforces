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

type testcase struct {
	a int
}

const testcasesRaw = `577
550
828
884
965
1429
1919
1899
120
337
1453
1817
1750
1074
732
415
511
310
454
1637
1761
1163
127
1215
1941
1754
406
871
431
1537
461
1903
2
1986
383
516
1743
180
265
912
1312
820
1271
1272
1921
55
984
1149
1613
655
774
1758
271
1199
102
1342
1607
120
1028
1976
1821
730
1324
1708
400
1254
1532
1
1690
1155
266
1872
939
1097
984
650
285
492
1688
1086
563
1672
101
1318
840
1161
1441
963
1705
627
1005
1218
1887
961
1866
1357
1233
1382
1698
613
901
1598
1693
239`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []testcase {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(raw)))
	cases := []testcase{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			panic(fmt.Sprintf("invalid integer %q: %v", line, err))
		}
		cases = append(cases, testcase{a: v})
	}
	if err := scanner.Err(); err != nil {
		panic(fmt.Sprintf("scanner error: %v", err))
	}
	if len(cases) == 0 {
		panic("no testcases parsed")
	}
	return cases
}

// generate builds the string output for a given input following 1202D.go logic.
func generate(a int) string {
	var buf bytes.Buffer
	mil := a / 1122457
	a %= 1122457
	thou := a / 1000
	one := a % 1000

	buf.WriteString("133")
	for j := 0; j < one; j++ {
		buf.WriteByte('7')
	}
	if thou > 0 {
		for j := 0; j < 994; j++ {
			buf.WriteByte('1')
		}
		buf.WriteString("33")
		for j := 0; j < thou; j++ {
			buf.WriteByte('7')
		}
	}
	if mil > 0 {
		for j := 0; j < 46; j++ {
			buf.WriteByte('3')
		}
		for j := 0; j < mil; j++ {
			buf.WriteByte('7')
		}
	}
	return buf.String()
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
	return out.String(), nil
}

func checkCase(bin string, idx int, tc testcase) error {
	input := fmt.Sprintf("1\n%d\n", tc.a)
	expected := generate(tc.a)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got := strings.TrimSpace(out)
	if got != expected {
		return fmt.Errorf("case %d: expected %q got %q", idx+1, expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
