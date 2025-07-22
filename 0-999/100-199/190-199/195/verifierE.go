package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod int64 = 1000000007

type DSU struct {
	parent []int
	weight []int64
}

func NewDSU(n int) *DSU {
	parent := make([]int, n+1)
	weight := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		parent[i] = i
		weight[i] = 0
	}
	return &DSU{parent, weight}
}

func (d *DSU) find(u int) (int, int64) {
	if d.parent[u] == u {
		return u, 0
	}
	root, dist := d.find(d.parent[u])
	d.parent[u] = root
	d.weight[u] += dist
	return root, d.weight[u]
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveEFromInput(input string) (int64, error) {
	r := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return 0, err
	}
	dsu := NewDSU(n)
	var ans int64
	for i := 1; i <= n; i++ {
		var k int
		fmt.Fscan(r, &k)
		for j := 0; j < k; j++ {
			var v int
			var x int64
			fmt.Fscan(r, &v, &x)
			root, depth := dsu.find(v)
			w := depth + x
			dsu.parent[root] = i
			dsu.weight[root] = w
			wmod := w % mod
			if wmod < 0 {
				wmod += mod
			}
			ans = (ans + wmod) % mod
		}
	}
	if ans < 0 {
		ans += mod
	}
	return ans % mod, nil
}

func generateCaseE(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		maxk := i - 1
		if maxk > 2 {
			maxk = 2
		}
		k := rng.Intn(maxk + 1)
		fmt.Fprintf(&sb, "%d", k)
		for j := 0; j < k; j++ {
			v := rng.Intn(i-1) + 1
			x := rng.Intn(21) - 10
			fmt.Fprintf(&sb, " %d %d", v, x)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expectInt, err := solveEFromInput(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
			os.Exit(1)
		}
		expect := fmt.Sprint(expectInt)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
