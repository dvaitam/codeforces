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

const mod = 998244353

type TestCaseE struct {
	n  int
	m  int
	b0 int
}

func genCaseE(rng *rand.Rand) TestCaseE {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	b0 := rng.Intn(3)
	return TestCaseE{n: n, m: m, b0: b0}
}

func possible(a []int, b0 int) bool {
	cur := map[int]struct{}{b0: {}}
	for _, x := range a {
		nxt := map[int]struct{}{}
		for v := range cur {
			if v+1 != x {
				nxt[v+1] = struct{}{}
			}
			if v-1 >= 0 && v-1 != x {
				nxt[v-1] = struct{}{}
			}
		}
		if len(nxt) == 0 {
			return false
		}
		cur = nxt
	}
	return true
}

func expectedE(tc TestCaseE) int64 {
	arr := make([]int, tc.n)
	var ans int64
	var dfs func(pos int)
	dfs = func(pos int) {
		if pos == tc.n {
			if possible(arr, tc.b0) {
				ans++
			}
			return
		}
		for v := 1; v <= tc.m; v++ {
			arr[pos] = v
			dfs(pos + 1)
		}
	}
	dfs(0)
	return ans % mod
}

func runCaseE(bin string, tc TestCaseE, expect int64) error {
	input := fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.m, tc.b0)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got%mod != expect {
		return fmt.Errorf("expected %d got %d", expect, got%mod)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseE(rng)
		exp := expectedE(tc)
		if err := runCaseE(bin, tc, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d m=%d b0=%d\n", i+1, err, tc.n, tc.m, tc.b0)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
