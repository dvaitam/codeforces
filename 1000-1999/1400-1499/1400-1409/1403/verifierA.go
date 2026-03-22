package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// Embedded correct solver for 1403A.
func solveCase(input string) string {
	r := strings.NewReader(input)
	var N, D, U, Q int
	fmt.Fscan(r, &N, &D, &U, &Q)
	H := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(r, &H[i])
	}
	upU := make([][2]int, U)
	for i := 0; i < U; i++ {
		fmt.Fscan(r, &upU[i][0], &upU[i][1])
	}
	type Query struct{ x, y, v, idx int }
	qs := make([]Query, Q)
	for i := 0; i < Q; i++ {
		fmt.Fscan(r, &qs[i].x, &qs[i].y, &qs[i].v)
		qs[i].idx = i
	}
	ord := make([]int, Q)
	for i := range ord {
		ord[i] = i
	}
	sort.Slice(ord, func(i, j int) bool {
		return qs[ord[i]].v < qs[ord[j]].v
	})
	adj := make([]map[int]struct{}, N)
	for i := 0; i < N; i++ {
		adj[i] = make(map[int]struct{})
	}
	ans := make([]int, Q)

	queryAns := func(x, y int) int {
		Nx := adj[x]
		Ny := adj[y]
		if len(Nx) == 0 || len(Ny) == 0 {
			return 1000000000
		}
		sx := make([]int, 0, len(Nx))
		for u := range Nx {
			sx = append(sx, H[u])
		}
		sy := make([]int, 0, len(Ny))
		for v := range Ny {
			sy = append(sy, H[v])
		}
		sort.Ints(sx)
		sort.Ints(sy)
		i, j, best := 0, 0, int(1e18)
		for i < len(sx) && j < len(sy) {
			a, b := sx[i], sy[j]
			d := a - b
			if d < 0 {
				d = -d
			}
			if d < best {
				best = d
			}
			if sx[i] < sy[j] {
				i++
			} else {
				j++
			}
		}
		return best
	}

	qi := 0
	for qi < Q && qs[ord[qi]].v == 0 {
		q := qs[ord[qi]]
		ans[q.idx] = queryAns(q.x, q.y)
		qi++
	}
	for day := 1; day <= U; day++ {
		a := upU[day-1][0]
		b := upU[day-1][1]
		if _, ok := adj[a][b]; ok {
			delete(adj[a], b)
			delete(adj[b], a)
		} else {
			adj[a][b] = struct{}{}
			adj[b][a] = struct{}{}
		}
		for qi < Q && qs[ord[qi]].v == day {
			q := qs[ord[qi]]
			ans[q.idx] = queryAns(q.x, q.y)
			qi++
		}
	}

	var sb strings.Builder
	for i := 0; i < Q; i++ {
		fmt.Fprintln(&sb, ans[i])
	}
	return strings.TrimSpace(sb.String())
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	N := rng.Intn(5) + 2
	D := rng.Intn(5) + 1
	U := rng.Intn(4) + 1
	Q := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", N, D, U, Q))
	for i := 0; i < N; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)))
	}
	sb.WriteByte('\n')
	for i := 0; i < U; i++ {
		a := rng.Intn(N)
		b := rng.Intn(N - 1)
		if b >= a {
			b++
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", a, b))
	}
	for i := 0; i < Q; i++ {
		x := rng.Intn(N)
		y := rng.Intn(N)
		v := rng.Intn(U + 1)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", x, y, v))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := solveCase(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed:\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
