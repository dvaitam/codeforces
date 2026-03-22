package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const solveMOD = 1000000007

type solvePair struct {
	g int64
	c int
}

func solveGCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	scanInt := func() int {
		scanner.Scan()
		var res int
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	scanInt64 := func() int64 {
		scanner.Scan()
		var res int64
		for _, b := range scanner.Bytes() {
			res = res*10 + int64(b-'0')
		}
		return res
	}

	scanner.Scan()
	n := 0
	for _, b := range scanner.Bytes() {
		n = n*10 + int(b-'0')
	}

	x := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		x[i] = scanInt64()
	}

	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := scanInt()
		v := scanInt()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	var ans int64

	var dfs func(u, p int, parentDP []solvePair)
	dfs = func(u, p int, parentDP []solvePair) {
		currDP := make([]solvePair, 0, len(parentDP)+1)
		currDP = append(currDP, solvePair{x[u], 1})

		for _, pair := range parentDP {
			g := solveGCD(pair.g, x[u])
			found := false
			for i := range currDP {
				if currDP[i].g == g {
					currDP[i].c = (currDP[i].c + pair.c) % solveMOD
					found = true
					break
				}
			}
			if !found {
				currDP = append(currDP, solvePair{g, pair.c})
			}
		}

		for _, pair := range currDP {
			term := (pair.g % solveMOD) * int64(pair.c) % solveMOD
			ans = (ans + term) % solveMOD
		}

		for _, v := range adj[u] {
			if v != p {
				dfs(v, u, currDP)
			}
		}
	}

	dfs(1, 0, nil)

	return fmt.Sprintf("%d", ans)
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < n-1; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		for a == b {
			b = rng.Intn(n) + 1
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		want := solve(input)
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
