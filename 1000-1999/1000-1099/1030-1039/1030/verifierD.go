package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type testCase struct {
	n, m, k  int64
	input    string
	possible bool
}

func solve(n, m, k int64) (bool, [6]int64) {
	tn, tm, tk := n, m, k
	tg := gcd(tn, tk)
	tn /= tg
	tk /= tg
	tg = gcd(tm, tk)
	tm /= tg
	tk /= tg
	var pts [6]int64
	if tk != 1 && tk != 2 {
		return false, pts
	}
	if tk == 2 {
		pts = [6]int64{0, 0, tn, 0, 0, tm}
		return true, pts
	}
	if tn*2 <= n {
		tn *= 2
	} else if tm*2 <= m {
		tm *= 2
	} else {
		return false, pts
	}
	pts = [6]int64{0, 0, tn, 0, 0, tm}
	return true, pts
}

func buildCase(n, m, k int64) testCase {
	ok, _ := solve(n, m, k)
	input := fmt.Sprintf("%d %d %d\n", n, m, k)
	return testCase{n: n, m: m, k: k, input: input, possible: ok}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Int63n(1000000000) + 1
	m := rng.Int63n(1000000000) + 1
	k := rng.Int63n(1000000000-1) + 2
	return buildCase(n, m, k)
}

func checkOutput(tc testCase, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := strings.ToUpper(scanner.Text())
	if first == "NO" {
		if tc.possible {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("invalid first token %s", first)
	}
	if !tc.possible {
		return fmt.Errorf("expected NO but got YES")
	}
	vals := make([]int64, 0, 6)
	for scanner.Scan() {
		v, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return fmt.Errorf("invalid integer %s", scanner.Text())
		}
		vals = append(vals, v)
	}
	if len(vals) != 6 {
		return fmt.Errorf("expected 6 integers, got %d", len(vals))
	}
	x1, y1, x2, y2, x3, y3 := vals[0], vals[1], vals[2], vals[3], vals[4], vals[5]
	if x1 < 0 || x1 > tc.n || x2 < 0 || x2 > tc.n || x3 < 0 || x3 > tc.n ||
		y1 < 0 || y1 > tc.m || y2 < 0 || y2 > tc.m || y3 < 0 || y3 > tc.m {
		return fmt.Errorf("points out of range")
	}
	area2 := abs((x2-x1)*(y3-y1) - (y2-y1)*(x3-x1))
	target2 := 2 * tc.n * tc.m / tc.k
	if area2 != target2 {
		return fmt.Errorf("wrong area")
	}
	return nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutput(tc, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		buildCase(2, 3, 2),
		buildCase(1, 1, 2),
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
