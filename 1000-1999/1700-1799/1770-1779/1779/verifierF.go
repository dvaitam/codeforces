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

const maxA = 32

func solveF(n int, a []int, parent []int) (string, string) {
	children := make([][]int, n+2)
	deg := make([]int, n+2)
	for i := 2; i <= n; i++ {
		p := parent[i]
		children[p] = append(children[p], i)
		deg[p]++
	}
	// euler tour order pre-order
	et := make([]int, 0, n)
	stack := []int{1}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		et = append(et, v)
		for i := len(children[v]) - 1; i >= 0; i-- {
			stack = append(stack, children[v][i])
		}
	}
	sz := make([]int, n+2)
	xr := make([]int, n+2)
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		sz[i] = 1
		xr[i] = a[i]
		if deg[i] == 0 {
			q = append(q, i)
		}
	}
	deg[1]++
	head := 0
	for head < len(q) {
		v := q[head]
		head++
		p := parent[v]
		sz[p] += sz[v]
		xr[p] ^= xr[v]
		deg[p]--
		if deg[p] == 0 {
			q = append(q, p)
		}
	}
	dp := make([][maxA]bool, n+2)
	if n+1 < len(dp) {
		dp[n+1][0] = true
	}
	for i := n; i >= 1; i-- {
		v := et[i-1]
		szi := sz[v]
		xri := xr[v]
		for x := 0; x < maxA; x++ {
			ok := dp[i+1][x]
			if !ok && szi%2 == 0 {
				nx := x ^ xri
				if nx < maxA && dp[i+szi][nx] {
					ok = true
				}
			}
			dp[i][x] = ok
		}
	}
	startX := xr[1]
	if startX >= maxA || !dp[1][startX] {
		return "-1", ""
	}
	ans := make([]int, 0, n)
	i := 1
	x := startX
	for i <= n {
		if dp[i+1][x] {
			i++
		} else {
			v := et[i-1]
			ans = append(ans, v)
			x ^= xr[v]
			i += sz[v]
		}
	}
	ans = append(ans, 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for idx, v := range ans {
		if idx > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return strings.TrimSpace(sb.String()), ""
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 2
	a := make([]int, n+2)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(32)
	}
	parent := make([]int, n+2)
	parent[1] = 1
	for i := 2; i <= n; i++ {
		parent[i] = rng.Intn(i-1) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteString("\n")
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", parent[i]))
	}
	sb.WriteString("\n")
	input := sb.String()
	expect, _ := solveF(n, a, parent)
	return input, expect
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
