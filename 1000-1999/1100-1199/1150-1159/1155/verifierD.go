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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func expected(n int, x int64, arr []int64) int64 {
	const inf = int64(4e18)
	var dp0, dp1, dp2, ans int64
	dp1 = -inf
	dp2 = -inf
	for i := 0; i < n; i++ {
		a := arr[i]
		noPrev := dp0
		mulPrev := dp1
		aftPrev := dp2
		dp0 = maxInt64(noPrev+a, 0)
		dp1 = maxInt64(mulPrev+a*x, noPrev+a*x)
		dp2 = maxInt64(aftPrev+a, mulPrev+a)
		ans = maxInt64(ans, dp0)
		ans = maxInt64(ans, dp1)
		ans = maxInt64(ans, dp2)
	}
	return ans
}

func checkCase(bin string, n int, x int64, arr []int64) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	input := sb.String()
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	want := fmt.Sprintf("%d", expected(n, x, arr))
	if strings.TrimSpace(out) != want {
		return fmt.Errorf("expected %s got %s", want, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	type caseD struct {
		n   int
		x   int64
		arr []int64
	}
	tests := []caseD{
		{n: 1, x: 1, arr: []int64{5}},
		{n: 3, x: -1, arr: []int64{1, 2, 3}},
	}
	for len(tests) < 100 {
		n := rng.Intn(7) + 1
		x := int64(rng.Intn(7) - 3)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(rng.Intn(11) - 5)
		}
		tests = append(tests, caseD{n: n, x: x, arr: arr})
	}
	for i, tc := range tests {
		if err := checkCase(bin, tc.n, tc.x, tc.arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
