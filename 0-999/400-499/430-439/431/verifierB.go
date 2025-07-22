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

func solveB(g [5][5]int64) string {
	perm := []int{0, 1, 2, 3, 4}
	var maxH int64
	var dfs func(int)
	dfs = func(idx int) {
		if idx == 5 {
			p := perm
			var sum int64
			sum += g[p[0]][p[1]] + g[p[1]][p[0]]
			sum += g[p[2]][p[3]] + g[p[3]][p[2]]
			sum += g[p[1]][p[2]] + g[p[2]][p[1]]
			sum += g[p[3]][p[4]] + g[p[4]][p[3]]
			sum += g[p[2]][p[3]] + g[p[3]][p[2]]
			sum += g[p[3]][p[4]] + g[p[4]][p[3]]
			if sum > maxH {
				maxH = sum
			}
			return
		}
		for i := idx; i < 5; i++ {
			perm[idx], perm[i] = perm[i], perm[idx]
			dfs(idx + 1)
			perm[idx], perm[i] = perm[i], perm[idx]
		}
	}
	dfs(0)
	return fmt.Sprintf("%d", maxH)
}

func generateCase(rng *rand.Rand) (string, string) {
	var g [5][5]int64
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			g[i][j] = int64(rng.Intn(10))
		}
	}
	var sb strings.Builder
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", g[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String(), solveB(g)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
