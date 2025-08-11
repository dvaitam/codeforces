package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	dir, _ := os.Getwd()
	oracle := filepath.Join(dir, "oracleE")
	cmd := exec.Command("go", "build", "-o", oracle, "732E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n, m int
	s    []int
	v    []int
}

func parseInput(in string) testCase {
	reader := strings.NewReader(in)
	var tc testCase
	fmt.Fscan(reader, &tc.n, &tc.m)
	tc.s = make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		fmt.Fscan(reader, &tc.s[i])
	}
	tc.v = make([]int, tc.m)
	for i := 0; i < tc.m; i++ {
		fmt.Fscan(reader, &tc.v[i])
	}
	return tc
}

func parseOutput(out string, n, m int) (num, cost int, t, ans []int, err error) {
	fields := strings.Fields(out)
	if len(fields) != 2+m+n {
		err = fmt.Errorf("expected %d integers, got %d", 2+m+n, len(fields))
		return
	}
	if num, err = strconv.Atoi(fields[0]); err != nil {
		return
	}
	if cost, err = strconv.Atoi(fields[1]); err != nil {
		return
	}
	t = make([]int, m)
	idx := 2
	for i := 0; i < m; i++ {
		if t[i], err = strconv.Atoi(fields[idx]); err != nil {
			return
		}
		idx++
	}
	ans = make([]int, n)
	for i := 0; i < n; i++ {
		if ans[i], err = strconv.Atoi(fields[idx]); err != nil {
			return
		}
		idx++
	}
	return
}

func verify(tc testCase, num, cost int, t, ans []int) error {
	if len(t) != tc.m || len(ans) != tc.n {
		return fmt.Errorf("wrong output length")
	}
	sum := 0
	assigned := make([]bool, tc.m)
	assignedCnt := 0
	for j := 0; j < tc.m; j++ {
		if t[j] < 0 {
			return fmt.Errorf("negative t for query %d", j+1)
		}
		sum += t[j]
	}
	if sum != cost {
		return fmt.Errorf("reported cost %d but sum of t is %d", cost, sum)
	}
	for i := 0; i < tc.n; i++ {
		j := ans[i]
		if j == 0 {
			continue
		}
		if j < 1 || j > tc.m {
			return fmt.Errorf("invalid assignment %d at position %d", j, i+1)
		}
		if assigned[j-1] {
			return fmt.Errorf("query %d assigned multiple times", j)
		}
		val := tc.v[j-1]
		for k := 0; k < t[j-1]; k++ {
			val = (val + 1) / 2
		}
		if val != tc.s[i] {
			return fmt.Errorf("query %d becomes %d after %d steps, expected %d", j, val, t[j-1], tc.s[i])
		}
		assigned[j-1] = true
		assignedCnt++
	}
	if assignedCnt != num {
		return fmt.Errorf("claimed %d assignments but have %d", num, assignedCnt)
	}
	for j := 0; j < tc.m; j++ {
		if !assigned[j] && t[j] != 0 {
			return fmt.Errorf("unassigned query %d has non-zero t", j+1)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(rng.Intn(20) + 1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []string
	for i := 0; i < 105; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, in := range cases {
		tc := parseInput(in)
		expOut, err := run(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expNum, expCost, _, _, err := parseOutput(expOut, tc.n, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle parse error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		gotOut, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		num, cost, t, ans, err := parseOutput(gotOut, tc.n, tc.m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if num != expNum || cost != expCost {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %d %d\ninput:%s", i+1, expNum, expCost, num, cost, in)
			os.Exit(1)
		}
		if err := verify(tc, num, cost, t, ans); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
