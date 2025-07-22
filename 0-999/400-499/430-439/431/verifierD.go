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

var dp [70][70]int64
var x [70]int
var K int

func dfs(pos, cnt int, flag bool) int64 {
	if pos == 0 {
		if cnt == K {
			return 1
		}
		return 0
	}
	if !flag && dp[pos][cnt] != -1 {
		return dp[pos][cnt]
	}
	u := 1
	if flag {
		u = x[pos]
	}
	var ret int64
	for i := 0; i <= u; i++ {
		ret += dfs(pos-1, cnt+i, flag && i == u)
	}
	if !flag {
		dp[pos][cnt] = ret
	}
	return ret
}

func ju(mid int64) int64 {
	e := 0
	for mid > 0 {
		e++
		x[e] = int(mid & 1)
		mid >>= 1
	}
	tmp1 := dfs(e, 0, true)
	e++
	for i := e; i >= 2; i-- {
		x[i] = x[i-1]
	}
	x[1] = 0
	tmp2 := dfs(e, 0, true)
	return tmp2 - tmp1
}

func solveD(m int64, k int) string {
	K = k
	for i := range dp {
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	l, r := int64(1), int64(1e18+1)
	for l < r {
		mid := (l + r) >> 1
		if ju(mid) >= m {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return fmt.Sprintf("%d", l)
}

func generateCase(rng *rand.Rand) (string, string) {
	m := rng.Int63n(1000)
	k := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d\n", m, k)
	return input, solveD(m, k)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	expStr := strings.TrimSpace(expected)
	if outStr != expStr {
		return fmt.Errorf("expected %q got %q", expStr, outStr)
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

	cases := []struct{ in, out string }{}
	for i := 0; i < 102; i++ {
		in, out := generateCase(rng)
		cases = append(cases, struct{ in, out string }{in, out})
	}

	for i, tc := range cases {
		if err := runCase(bin, tc.in, tc.out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
