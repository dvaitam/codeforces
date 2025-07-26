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

func runCandidate(bin string, input string) (string, error) {
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

const mod int = 1e9 + 7

type Basis struct {
	b    [20]int
	size int
}

func (bs *Basis) Add(x int) {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] != 0 {
			x ^= bs.b[i]
		} else {
			bs.b[i] = x
			bs.size++
			return
		}
	}
}

func (bs *Basis) Contains(x int) bool {
	for i := 19; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] == 0 {
			return false
		}
		x ^= bs.b[i]
	}
	return true
}

func solveF(data string) string {
	reader := bufio.NewReader(strings.NewReader(data))
	var n, q int
	fmt.Fscan(reader, &n, &q)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	type Query struct {
		x   int
		idx int
	}
	queries := make([][]Query, n+1)
	for i := 0; i < q; i++ {
		var l, x int
		fmt.Fscan(reader, &l, &x)
		queries[l] = append(queries[l], Query{x, i})
	}
	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}
	ans := make([]int, q)
	var bs Basis
	for i := 1; i <= n; i++ {
		bs.Add(arr[i])
		for _, qu := range queries[i] {
			if bs.Contains(qu.x) {
				ans[qu.idx] = pow2[i-bs.size]
			} else {
				ans[qu.idx] = 0
			}
		}
	}
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func genCaseF(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	q := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1 << 20)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		l := rng.Intn(n) + 1
		x := rng.Intn(1 << 20)
		sb.WriteString(fmt.Sprintf("%d %d\n", l, x))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseF(rng)
		expect := solveF(in)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, expect, strings.TrimSpace(got), in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
