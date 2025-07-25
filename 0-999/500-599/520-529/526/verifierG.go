package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct{ to, w int }

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveG(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, q int
	fmt.Fscan(in, &n, &q)
	adj := make([][]edge, n)
	for i := 0; i < n-1; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		adj[u] = append(adj[u], edge{v, w})
		adj[v] = append(adj[v], edge{u, w})
	}
	down1 := make([]int, n)
	down2 := make([]int, n)
	child1 := make([]int, n)
	up := make([]int, n)
	parent := make([]int, n)
	var dfs1 func(int, int)
	dfs1 = func(u, p int) {
		parent[u] = p
		down1[u], down2[u] = 0, 0
		for _, e := range adj[u] {
			v, w := e.to, e.w
			if v == p {
				continue
			}
			dfs1(v, u)
			d := down1[v] + w
			if d > down1[u] {
				down2[u], down1[u] = down1[u], d
				child1[u] = v
			} else if d > down2[u] {
				down2[u] = d
			}
		}
	}
	var dfs2 func(int, int)
	dfs2 = func(u, p int) {
		for _, e := range adj[u] {
			v, w := e.to, e.w
			if v == p {
				continue
			}
			viaUp := up[u] + w
			viaSib := 0
			if child1[u] == v {
				viaSib = down2[u] + w
			} else {
				viaSib = down1[u] + w
			}
			up[v] = maxInt(viaUp, viaSib)
			dfs2(v, u)
		}
	}
	dfs1(0, -1)
	up[0] = 0
	dfs2(0, -1)
	D1 := make([]int, n)
	D2 := make([]int, n)
	Bs := make([][]int64, n)
	PBs := make([][]int64, n)
	for u := 0; u < n; u++ {
		var ds []int
		ds = append(ds, up[u])
		for _, e := range adj[u] {
			v, w := e.to, e.w
			if v == parent[u] {
				continue
			}
			ds = append(ds, down1[v]+w)
		}
		a, b := 0, 0
		for _, d := range ds {
			if d > a {
				b = a
				a = d
			} else if d > b {
				b = d
			}
		}
		D1[u], D2[u] = a, b
		remA, remB := 1, 1
		var bs []int64
		for _, d := range ds {
			if d == a && remA > 0 {
				remA--
				continue
			}
			if d == b && remB > 0 {
				remB--
				continue
			}
			bs = append(bs, int64(d))
		}
		if len(bs) > 1 {
			for i := 0; i < len(bs); i++ {
				for j := i + 1; j < len(bs); j++ {
					if bs[j] > bs[i] {
						bs[i], bs[j] = bs[j], bs[i]
					}
				}
			}
		}
		Bs[u] = bs
		ps := make([]int64, len(bs)+1)
		for i := 0; i < len(bs); i++ {
			ps[i+1] = ps[i] + bs[i]
		}
		PBs[u] = ps
	}
	var x, y, ansPrev int64
	var out strings.Builder
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &x, &y)
		if i > 0 {
			x = ((x + ansPrev - 1) % int64(n)) + 1
			y = ((y + ansPrev - 1) % int64(n)) + 1
		}
		u := int(x - 1)
		k := int(y) - 1
		if k < 0 {
			k = 0
		}
		if k > len(PBs[u])-1 {
			k = len(PBs[u]) - 1
		}
		ans := int64(D1[u] + D2[u])
		if k > 0 {
			ans += PBs[u][k]
		}
		if i > 0 {
			out.WriteByte('\n')
		}
		out.WriteString(fmt.Sprint(ans))
		ansPrev = ans
	}
	return out.String()
}

func genTestG(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	q := rng.Intn(5) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d %d\n", n, q)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(10) + 1
		parent[i] = p
		fmt.Fprintf(&buf, "%d %d %d\n", p, i, w)
	}
	for i := 0; i < q; i++ {
		x := rng.Intn(n) + 1
		y := rng.Intn(n) + 1
		fmt.Fprintf(&buf, "%d %d\n", x, y)
	}
	return buf.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestG(rng)
		expect := solveG(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
