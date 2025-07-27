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

type Test struct {
	n int
	k int
	a []int
	b []int
}

func (t Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for i := 0; i < t.n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", t.a[i], t.b[i]))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expected(t Test) string {
	n, k := t.n, t.k
	sumR := 0
	sumB := 0
	for i := 0; i < n; i++ {
		sumR += t.a[i]
		sumB += t.b[i]
	}
	dp := make([]bool, k)
	dp[0] = true
	for i := 0; i < n; i++ {
		ndp := make([]bool, k)
		for r := 0; r < k; r++ {
			if !dp[r] {
				continue
			}
			maxx := min(t.a[i], k-1)
			for x := 0; x <= maxx; x++ {
				y := (k - ((r + x) % k)) % k
				if y <= t.b[i] && x+y <= t.a[i]+t.b[i] {
					ndp[(r+x)%k] = true
				}
			}
		}
		dp = ndp
	}
	total := sumR + sumB
	res := total / k
	if !dp[sumR%k] {
		res--
	}
	if res < 0 {
		res = 0
	}
	return fmt.Sprintf("%d", res)
}

func runProg(bin, input string) (string, error) {
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

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(10) + 1
	k := rng.Intn(10) + 1
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(10)
		b[i] = rng.Intn(10)
	}
	return Test{n: n, k: k, a: a, b: b}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
