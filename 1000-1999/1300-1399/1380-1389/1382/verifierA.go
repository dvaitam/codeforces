package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1382ASource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n, m int
       fmt.Fscan(in, &n, &m)
       seen := make(map[int]bool, n)
       for i := 0; i < n; i++ {
           var v int
           fmt.Fscan(in, &v)
           seen[v] = true
       }
       res := -1
       for i := 0; i < m; i++ {
           var v int
           fmt.Fscan(in, &v)
           if res == -1 && seen[v] {
               res = v
           }
       }
       if res != -1 {
           fmt.Fprintln(out, "YES")
           fmt.Fprintln(out, 1, res)
       } else {
           fmt.Fprintln(out, "NO")
       }
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution1382ASource

type testCase struct {
	a []int
	b []int
}

const testcasesRaw = `100
4 4
1 3 5 4
4 3 4 3
5 2
5 2 3 2 1
5 3
5 5
2 3 1 1 3
4 5 1 3 4
3 5
2 5 4
4 5 3 1 5
1 1
4
1
5 4
3 2 3 1 2
5 2 2 2
5 4
1 1 3 5 4
1 3 5 3
1 5
3
5 2 5 5 5
3 4
1 5 4
3 5 2 3
2 2
2 1
5 3
4 1
1 2 2 1
1
5 4
5 3 5 2 2
5 4 5 3
4 4
3 1 3 5
1 4 5 3
2 5
1 2
2 1 1 1 1
2 3
3 1 1 2
1 1
3
2 2
2 1
2 1
5 4
4 2 4 3 1
1 1 2 1
3 3
3 3 2
1 1 2
2 1
2 2
2
4 5
1 2 5 3
3 1 1 3 5
5 2
1 2 4 5 4
2 5
5 3
3 4 4 5 1
3 2 2
2 2
2 2
5 1
4 1 1 1 1
5
5 2
1 3 4 2 5
3 3
1 3 3
3 1 3
5 4
3 4 5 5 2
1 4 2 3
1 3
3
2
5 2
1 4 2 1 2
5 4
5 2 4 5 2
2 1 1 1
4 5
3 4 4 3
1 2 4 1 5
2 2
2 2
1 2
4 5
2 3 1 4
3 1 3 5 3
5 3
2 5 4 3 4
1 3 5
1 2
1
3
3 1
2 2 1
1 2
5
1 5
2 4 3 2 5
1 3 4 3
4 2
1 1 3 5
5 2
4 3
4 4
2 5 3 2
2 1 1 1
3 1
1 2 2
5 5
1 1 3 2 2
1 4 2 2 1
1 1
1
5 2
5 4 4 3 1
3 3
1 4 3
1 1 1
5 2
5 4 4 5 1
5 4
3 4 1 1 3
4 2 2 1
4 1
1 4 2 3
2 4
4 4
1 4 2 2
1 4 1 4
1 4
3
5 5
1 4 2 1 5
5 4 5 2 4
5 2
3 5 1 5 4
3 4
2 3 3 1
5 3
3 1 4 5 3
5 1
4 5 5 3 1
2 5
5 2 1 2 5
3 3
2 5 4
4 2 2
2 3
3 5 5
5
5 5
2 1 5 4 4
1 1 2 3 1
3 4
3 4 2
3 2 4 2
2 3
4 2 3
1 1 2
4 5
3 5 1 1
5 3 2 2 1
1 1
2
2
4 4
2 3 3 3
1 2 1 1
2 2
1 2
3 3
5 4 3
2 2 3
2 5
5 5
5 4
3 2 1 5 4
4 3 1 2
1 2
3
3 1
2 1 3
5 1
5 3 3 3 2
3 5
2 5 2
5 2 2 1 3
4 2
1 2 1 1
2 1
1 3
4 3
3 1 2 1
1 4 5
2 2
1 1
2 3
3 5 4
4 1
1 3 2 4
2 3
1 1 1
3 2
4 3 4
3 3
1 4 1
2 5
1 3
1 4
5 2
4 3 2 3 2
2 3
2 2 1
5 4
5 5 5 4 4
5 2 1 5
5 1
5 2 3 1 1
1
5 4
3 1 4 5 2
3 3
4 2 3
4 3
2 2 2 1
2 2
4 2 2 3
2 2
1 1
5 2
5 3 3 2 4
4 2
2 5 2 5
4 2
4 1 2 4
2 4
2 1 1 5
4 5
4 3 2 3
1 3 2 4 2
4 5
2 4 3 2
3 3 2 2 3
4 2
2 3 1 1
2 1
1 4
4 3
4 3 2 2
5 3
2 1 1 2 3
5 5
3 5 4 2 2
1 4 2 5 3
5 1
3 2 3 1 1
1 1
1
1 3
3
1 5
5
2 3
2 1
1 1 2
1 4
2 5 3 4
2 4
5 5 2 1
2 2
4 1
4 5
2 3 1 3 2
5 1
3 4 2 3 5
3 5
1 3 3 3
1 2
3
4 1
1 1 1 1
3 1
5 1 4
1 3
5
4 4
2 3 3 3
4 4
2 3 4 1
1 3
3
2 2
4 1
4 1
2 1
1 2
3 2
2 5 1
4 3
4 1 1 4
2 4
2 4 4 3
1 5
2 4 2 3 5
1 3
5
4 3
4 2 5 2
5 3
2 1 4 3 3
3 4
1 3 4
5 1
1 4 5 4 4
5 3
1 2 4 3 3
4 1
2 5 2 4
3 2
2 5 5
4 4
4 5 1 2
2 5
1 3 1 1 2
1 3
4
5 5
2 1 3 5 4
5 3
3 4 2 1 2
3 1
1 5 1
2 3
3 2 5
2 3
1 5 2
3 3
4 4 3
2 4
5 1 4 5
4 5
4 1 1 3 2
2 5
4 2 2 5 2
3 5
3 2 5 4 2
2 5
5 3 3 3 1
4 2
5 3 3 2
4 2
2 1 1 2
2 2
3 5
1 5 2
3 4
3 2 1 5
1 4
5
2 4
4 3 3 3
5 2
2 4 2 4 5
4 2
5 4 5 5
4 2
4 2 3 1
4 1
1 4 3 2
2 3
3 4 1 2
1 2
5
`

func parseTestcases() []testCase {
	reader := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil
	}
	res := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[j])
		}
		b := make([]int, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(reader, &b[j])
		}
		res[i] = testCase{a: a, b: b}
	}
	return res
}

func solveExpected(tc testCase) (bool, int) {
	seen := make(map[int]bool, len(tc.a))
	for _, v := range tc.a {
		seen[v] = true
	}
	for _, v := range tc.b {
		if seen[v] {
			return true, v
		}
	}
	return false, 0
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", len(tc.a), len(tc.b)))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	ok, val := solveExpected(tc)
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	tokens := strings.Fields(strings.TrimSpace(string(out)))
	if len(tokens) == 0 {
		return fmt.Errorf("case %d failed: empty output", idx)
	}
	if strings.ToUpper(tokens[0]) == "NO" {
		if ok {
			return fmt.Errorf("case %d failed: expected YES got NO\ninput:\n%s", idx, input)
		}
		return nil
	}
	if strings.ToUpper(tokens[0]) != "YES" {
		return fmt.Errorf("case %d failed: first token must be YES or NO\ninput:\n%s", idx, input)
	}
	if !ok {
		return fmt.Errorf("case %d failed: expected NO got YES\ninput:\n%s", idx, input)
	}
	if len(tokens) < 3 {
		return fmt.Errorf("case %d failed: missing count/value\ninput:\n%s", idx, input)
	}
	cnt, err := strconv.Atoi(tokens[1])
	if err != nil || cnt != 1 {
		return fmt.Errorf("case %d failed: expected count 1 got %s\ninput:\n%s", idx, tokens[1], input)
	}
	v, err := strconv.Atoi(tokens[2])
	if err != nil || v != val {
		return fmt.Errorf("case %d failed: expected value %d got %s\ninput:\n%s", idx, val, tokens[2], input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
