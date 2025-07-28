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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n, m int, a []int64, qs [][2]int) []int64 {
	prefixRight := make([]int64, n)
	for i := 1; i < n; i++ {
		diff := a[i-1] - a[i]
		if diff > 0 {
			prefixRight[i] = prefixRight[i-1] + diff
		} else {
			prefixRight[i] = prefixRight[i-1]
		}
	}
	prefixLeft := make([]int64, n)
	for i := n - 2; i >= 0; i-- {
		diff := a[i+1] - a[i]
		if diff > 0 {
			prefixLeft[i] = prefixLeft[i+1] + diff
		} else {
			prefixLeft[i] = prefixLeft[i+1]
		}
	}
	ans := make([]int64, m)
	for idx, q := range qs {
		s, t := q[0], q[1]
		if s < t {
			ans[idx] = prefixRight[t-1] - prefixRight[s-1]
		} else {
			ans[idx] = prefixLeft[t-1] - prefixLeft[s-1]
		}
	}
	return ans
}

func genCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(9) + 2 // 2..10
	m := rng.Intn(5) + 1 // 1..5
	a := make([]int64, n)
	for i := range a {
		a[i] = rng.Int63n(20) + 1
	}
	qs := make([][2]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		s := rng.Intn(n) + 1
		t := rng.Intn(n) + 1
		for s == t {
			t = rng.Intn(n) + 1
		}
		qs[i] = [2]int{s, t}
		sb.WriteString(fmt.Sprintf("%d %d\n", s, t))
	}
	expect := solveCase(n, m, a, qs)
	return sb.String(), expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		outFields := strings.Fields(out)
		if len(outFields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", i+1, len(expect), len(outFields), in)
			os.Exit(1)
		}
		for j, val := range expect {
			if outFields[j] != fmt.Sprintf("%d", val) {
				fmt.Fprintf(os.Stderr, "case %d failed on query %d: expected %d got %s\ninput:\n%s", i+1, j+1, val, outFields[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
