package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type Test struct {
	n     int
	l     int
	r     int
	a     []int
	b     []int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(8) + 1
	sum := 0
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
		sum += a[i]
	}
	l := rng.Intn(sum + 1)
	r := l + rng.Intn(sum-l+1)
	b := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, r))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 0
		} else {
			b[i] = 1
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	return Test{n: n, l: l, r: r, a: a, b: b, input: sb.String()}
}

func solve(t Test) string {
	sum := 0
	for _, v := range t.a {
		sum += v
	}
	const negInf = -1 << 30
	dp := make([]int, sum+1)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0
	for i := 0; i < t.n; i++ {
		h := t.a[i]
		imp := t.b[i] == 1
		for s := sum - h; s >= 0; s-- {
			if dp[s] == negInf {
				continue
			}
			gain := 0
			if imp && s >= t.l && s <= t.r {
				gain = 1
			}
			if v := dp[s] + gain; v > dp[s+h] {
				dp[s+h] = v
			}
		}
	}
	return fmt.Sprintf("%d", dp[sum])
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
