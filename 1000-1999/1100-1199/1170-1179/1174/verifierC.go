package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `49
45
47
31
50
4
38
9
28
26
47
31
1
7
3
39
47
50
16
50
20
12
40
45
4
28
23
38
15
20
16
22
9
23
33
13
48
47
31
49
5
14
46
16
14
50
12
1
3
5
30
10
5
50
14
3
48
7
1
34
16
40
40
26
46
34
8
46
22
19
42
41
42
37
35
22
43
39
29
16
13
39
35
23
49
50
10
9
14
49
49
26
26
19
13
49
25
1
29
38`

// solve implements the logic from 1174C.go.
func solve(n int) string {
	if n < 2 {
		return ""
	}
	primes := make([]int, 0, int(math.Log(float64(n))+3))
	var sb strings.Builder
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		limit := int(math.Sqrt(float64(i)))
		first := 0
		found := false
		for j, p := range primes {
			if p > limit {
				break
			}
			if i%p == 0 {
				found = true
				first = j + 1
				break
			}
		}
		if !found {
			primes = append(primes, i)
			first = len(primes)
		}
		sb.WriteString(strconv.Itoa(first))
	}
	return sb.String()
}

type testCase struct {
	n int
}

func parseTestcases() ([]testCase, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	var tests []testCase
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" {
			continue
		}
		v, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		tests = append(tests, testCase{n: v})
	}
	if err := scan.Err(); err != nil {
		return nil, err
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("%d\n", tc.n)
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	expected := solve(tc.n)
	if strings.TrimSpace(got) != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
