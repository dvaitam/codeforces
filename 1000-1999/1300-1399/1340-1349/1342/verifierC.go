package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type query struct {
	l, r int64
}

type testCase struct {
	a, b    int64
	queries []query
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.a, tc.b, len(tc.queries)))
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.l, q.r))
	}
	return sb.String()
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(tc testCase) string {
	a, b := tc.a, tc.b
	if a > b {
		a, b = b, a
	}
	lcm := a / gcd(a, b) * b
	period := int(lcm)
	pre := make([]int, period+1)
	for i := 0; i < period; i++ {
		if (int64(i)%a)%b != (int64(i)%b)%a {
			pre[i+1] = pre[i] + 1
		} else {
			pre[i+1] = pre[i]
		}
	}
	total := pre[period]
	calc := func(x int64) int64 {
		if x < 0 {
			return 0
		}
		q := x / lcm
		r := int(x % lcm)
		return int64(total)*q + int64(pre[r+1])
	}
	var ans []string
	for _, q := range tc.queries {
		val := calc(q.r) - calc(q.l-1)
		ans = append(ans, fmt.Sprint(val))
	}
	return strings.Join(ans, " ")
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(3))
	tests := make([]testCase, 100)
	for i := range tests {
		a := int64(rng.Intn(20) + 1)
		b := int64(rng.Intn(20) + 1)
		q := rng.Intn(5) + 1
		qs := make([]query, q)
		for j := 0; j < q; j++ {
			l := int64(rng.Intn(1000))
			r := l + int64(rng.Intn(1000))
			qs[j] = query{l: l, r: r}
		}
		tests[i] = testCase{a: a, b: b, queries: qs}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		input := tc.Input()
		want := solve(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if strings.TrimSpace(out) != want {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
