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

type edge struct {
	to int
	w  int
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(n, k int, adj [][]edge) string {
	Np := n + 2
	encode := func(L, R, pos int) int { return ((L*Np+R)*n + (pos - 1)) }
	decode := func(code int) (int, int, int) {
		pos := code%n + 1
		code /= n
		R := code % Np
		L := code / Np
		return L, R, pos
	}
	const INF int64 = 1 << 60
	dpPrev := map[int]int64{}
	for i := 1; i <= n; i++ {
		idx := encode(0, n+1, i)
		dpPrev[idx] = 0
	}
	for step := 1; step < k; step++ {
		dpNext := map[int]int64{}
		for idx, cost := range dpPrev {
			L, R, pos := decode(idx)
			for _, e := range adj[pos] {
				nxt := e.to
				if nxt <= L || nxt >= R || nxt == pos {
					continue
				}
				var newL, newR int
				if nxt > pos {
					newL = pos
					newR = R
				} else {
					newL = L
					newR = pos
				}
				nidx := encode(newL, newR, nxt)
				nc := cost + int64(e.w)
				if old, ok := dpNext[nidx]; !ok || nc < old {
					dpNext[nidx] = nc
				}
			}
		}
		dpPrev = dpNext
		if len(dpPrev) == 0 {
			break
		}
	}
	ans := INF
	for _, c := range dpPrev {
		if c < ans {
			ans = c
		}
	}
	if ans == INF {
		return "-1"
	}
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) (int, int, [][]edge) {
	n := rng.Intn(5) + 2
	k := rng.Intn(n) + 1
	m := rng.Intn(10)
	adj := make([][]edge, n+1)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		c := rng.Intn(10) + 1
		adj[u] = append(adj[u], edge{v, c})
	}
	return n, k, adj
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, k, adj := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		m := 0
		for _, edges := range adj {
			m += len(edges)
		}
		fmt.Fprintf(&sb, "%d\n", m)
		for u := 1; u <= n; u++ {
			for _, e := range adj[u] {
				fmt.Fprintf(&sb, "%d %d %d\n", u, e.to, e.w)
			}
		}
		expect := solveD(n, k, adj)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
