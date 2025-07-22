package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type TestCase struct {
	input  string
	expect string
}

func runCandidate(bin, input string) (string, error) {
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

func expectedA(n int, l, t int64, a []int64) string {
	t *= 2
	full := t / l
	left := t % l
	var res int64
	j := 0
	phase := false
	for i := 0; i < n; i++ {
		if left+a[i] >= l {
			phase = true
			j = 0
			left -= l
		}
		for j < n && a[j] <= left+a[i] {
			j++
		}
		if !phase {
			res += (full+1)*int64(j-i-1) + full*int64(n-j+i)
		} else {
			res += (full+1)*int64(n+j-i-1) + full*int64(i-j)
		}
	}
	ans := float64(res) / 4.0
	return fmt.Sprintf("%.6f", ans)
}

func genCase(rng *rand.Rand) TestCase {
	n := rng.Intn(5) + 2
	l := int64(rng.Intn(50) + int(n) + 5)
	t := int64(rng.Intn(50) + 1)
	vals := rng.Perm(int(l))[:n]
	sort.Ints(vals)
	a := make([]int64, n)
	for i, v := range vals {
		a[i] = int64(v)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, l, t)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", a[i])
	}
	sb.WriteByte('\n')
	expect := expectedA(n, l, t, a)
	return TestCase{sb.String(), expect}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if out != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.expect, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
