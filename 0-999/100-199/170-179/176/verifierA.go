package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m, k int, a, b, c [][]int) int {
	best := 0
	for i := 0; i < n; i++ {
		for h := 0; h < n; h++ {
			if i == h {
				continue
			}
			profits := make([]int, 0)
			for j := 0; j < m; j++ {
				p := b[h][j] - a[i][j]
				if p > 0 {
					for cnt := 0; cnt < c[i][j]; cnt++ {
						profits = append(profits, p)
					}
				}
			}
			sort.Sort(sort.Reverse(sort.IntSlice(profits)))
			sum := 0
			limit := k
			if len(profits) < limit {
				limit = len(profits)
			}
			for t := 0; t < limit; t++ {
				sum += profits[t]
			}
			if sum > best {
				best = sum
			}
		}
	}
	return best
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rand.Seed(0)
	for tcase := 0; tcase < 100; tcase++ {
		n := rand.Intn(3) + 2 // 2..4
		m := rand.Intn(3) + 1 // 1..3
		k := rand.Intn(5) + 1 // 1..5
		names := make([]string, n)
		a := make([][]int, n)
		bvals := make([][]int, n)
		cvals := make([][]int, n)
		for i := 0; i < n; i++ {
			names[i] = fmt.Sprintf("p%d", rand.Intn(100))
			a[i] = make([]int, m)
			bvals[i] = make([]int, m)
			cvals[i] = make([]int, m)
			for j := 0; j < m; j++ {
				a[i][j] = rand.Intn(10) + 1
				bvals[i][j] = rand.Intn(10) + 1
				cvals[i][j] = rand.Intn(5)
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			fmt.Fprint(&sb, names[i])
			for j := 0; j < m; j++ {
				fmt.Fprintf(&sb, " %d %d %d", a[i][j], bvals[i][j], cvals[i][j])
			}
			fmt.Fprintln(&sb)
		}
		input := sb.String()
		exp := expected(n, m, k, a, bvals, cvals)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tcase+1, err)
			return
		}
		var got int
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Printf("test %d failed: expected %d got %d\ninput:\n%s", tcase+1, exp, got, input)
			return
		}
	}
	fmt.Println("All tests passed.")
}
