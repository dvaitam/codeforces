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

// testcasesRaw embeds contents of testcasesA.txt.
const testcasesRaw = `100
6312
6891
664
4243
8377
7962
6635
4970
7809
5867
9559
3579
8269
2282
4618
2290
1554
4105
8726
9862
2408
5082
1619
1209
5410
7736
9172
1650
5797
7114
5181
3351
9053
7816
7254
8542
4268
1021
8990
231
1529
6535
19
8087
5459
3997
5329
1032
3131
9299
3633
3910
2335
8897
7340
1495
1319
5244
8323
8017
1787
4939
9032
4770
2045
8970
5452
8853
3330
9883
8966
9628
4713
7291
1502
9770
6307
5195
9432
3967
4757
3013
3103
3060
541
4261
7808
1132
1472
2134
2451
634
1315
8858
6411
8595
4516
8550
3859
3526`

type testCase struct {
	n int
}

// referenceSolution embeds 125A.go logic so no external oracle is needed.
func referenceSolution(n int) (int, int) {
	inches := n / 3
	if n%3 == 2 {
		inches++
	}
	feet := inches / 12
	inches %= 12
	return feet, inches
}

func parseTestcases() []testCase {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		panic("no testcase count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		panic("invalid testcase count")
	}
	tests := make([]testCase, 0, t)
	for scan.Scan() {
		val, err := strconv.Atoi(scan.Text())
		if err != nil {
			panic("invalid testcase value")
		}
		tests = append(tests, testCase{n: val})
	}
	if len(tests) != t {
		panic("testcase count mismatch")
	}
	return tests
}

func solve(n int) (int, int) {
	inches := n / 3
	if n%3 == 2 {
		inches++
	}
	feet := inches / 12
	inches %= 12
	return feet, inches
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
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
	for i, tc := range tests {
		expectFeet, expectInches := referenceSolution(tc.n)
		expected := fmt.Sprintf("%d %d", expectFeet, expectInches)
		input := fmt.Sprintf("%d\n", tc.n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("case %d failed: expected %q got %q\n", i+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
