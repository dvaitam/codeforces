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

type restr struct{ l, r, x, c int }

func expected(n, h int, res []restr) int {
	heights := make([]int, n)
	best := -1 << 31
	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			profit := 0
			for i := 0; i < n; i++ {
				profit += heights[i] * heights[i]
			}
			for _, R := range res {
				mx := 0
				for i := R.l - 1; i < R.r; i++ {
					if heights[i] > mx {
						mx = heights[i]
					}
				}
				if mx > R.x {
					profit -= R.c
				}
			}
			if profit > best {
				best = profit
			}
			return
		}
		for k := 0; k <= h; k++ {
			heights[pos] = k
			dfs(pos + 1)
		}
	}
	dfs(0)
	return best
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 1
	h := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	res := make([]restr, m)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, h, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		x := rng.Intn(h + 1)
		c := rng.Intn(5) + 1
		res[i] = restr{l, r, x, c}
		fmt.Fprintf(&sb, "%d %d %d %d\n", l, r, x, c)
	}
	exp := expected(n, h, res)
	return sb.String(), exp
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != fmt.Sprint(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
