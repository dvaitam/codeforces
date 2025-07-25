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

const modE = 1000000007

type testCaseE struct {
	n int
	k int
	s string
}

func generateCase(rng *rand.Rand) (string, testCaseE) {
	n := rng.Intn(7) + 1
	k := rng.Intn(n)
	digits := make([]byte, n)
	for i := 0; i < n; i++ {
		digits[i] = byte(rng.Intn(10)) + '0'
	}
	s := string(digits)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%s\n", n, k, s)
	return sb.String(), testCaseE{n: n, k: k, s: s}
}

func evalExpression(parts []string) int64 {
	var sum int64
	for _, p := range parts {
		var v int64
		for i := 0; i < len(p); i++ {
			v = v*10 + int64(p[i]-'0')
		}
		sum += v
	}
	return sum
}

func expected(tc testCaseE) int64 {
	n, k := tc.n, tc.k
	indices := make([]int, 0, k)
	var best int64
	var count int64
	var dfs func(pos int)
	dfs = func(pos int) {
		if len(indices) == k {
			parts := make([]string, 0, k+1)
			last := 0
			for _, idx := range indices {
				parts = append(parts, tc.s[last:idx])
				last = idx
			}
			parts = append(parts, tc.s[last:])
			val := evalExpression(parts)
			if val > best {
				best = val
				count = 1
			} else if val == best {
				count++
			}
			return
		}
		if pos >= n-1 {
			return
		}
		// choose no plus at pos
		dfs(pos + 1)
		// add plus after pos
		indices = append(indices, pos+1)
		dfs(pos + 1)
		indices = indices[:len(indices)-1]
	}
	dfs(0)
	return best % modE
}

func runCase(bin string, input string, tc testCaseE) error {
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
	gotStr := strings.TrimSpace(out.String())
	var got int64
	if _, err := fmt.Sscan(gotStr, &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	want := expected(tc)
	if got%modE != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
