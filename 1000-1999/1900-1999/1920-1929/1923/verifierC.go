package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type query struct{ l, r int }

type test struct {
	n  int
	q  int
	c  []int64
	qs []query
}

func solve(n int, c []int64, qs []query) []string {
	prefSum := make([]int64, n+1)
	prefOnes := make([]int, n+1)
	for i := 1; i <= n; i++ {
		prefSum[i] = prefSum[i-1] + c[i-1]
		prefOnes[i] = prefOnes[i-1]
		if c[i-1] == 1 {
			prefOnes[i]++
		}
	}
	ans := make([]string, len(qs))
	for i, q := range qs {
		length := q.r - q.l + 1
		sum := prefSum[q.r] - prefSum[q.l-1]
		ones := prefOnes[q.r] - prefOnes[q.l-1]
		if length <= 1 || sum < int64(length+ones) {
			ans[i] = "NO"
		} else {
			ans[i] = "YES"
		}
	}
	return ans
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	tests := make([]test, 0, 100)
	for len(tests) < 100 {
		n := rng.Intn(8) + 2
		q := rng.Intn(5) + 1
		c := make([]int64, n)
		for i := 0; i < n; i++ {
			c[i] = int64(rng.Intn(5) + 1)
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			qs[i] = query{l, r}
		}
		tests = append(tests, test{n, q, c, qs})
	}
	return tests
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.q))
		for j, v := range t.c {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteString("\n")
		for _, qq := range t.qs {
			sb.WriteString(fmt.Sprintf("%d %d\n", qq.l, qq.r))
		}
		expLines := solve(t.n, t.c, t.qs)
		expected := strings.Join(expLines, "\n")
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, sb.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
