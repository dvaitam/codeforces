package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type nodeKey struct {
	f int
	c int
}

type edgeE struct {
	to int
	w  int64
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var T int
	fmt.Fscan(in, &T)
	var out strings.Builder
	for ; T > 0; T-- {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		x := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &x[i])
		}
		type ladder struct {
			a, b, c, d int
			h          int64
			s, e       int
		}
		ladders := make([]ladder, k)
		idxMap := make(map[nodeKey]int)
		floors := []int{}
		cols := []int{}
		nodesByFloor := make(map[int][]int)
		var getIndex func(int, int) int
		getIndex = func(f, c int) int {
			key := nodeKey{f, c}
			if id, ok := idxMap[key]; ok {
				return id
			}
			id := len(floors)
			idxMap[key] = id
			floors = append(floors, f)
			cols = append(cols, c)
			nodesByFloor[f] = append(nodesByFloor[f], id)
			return id
		}
		startIdx := getIndex(1, 1)
		endIdx := getIndex(n, m)
		for i := 0; i < k; i++ {
			var a, b, c, d int
			var h int64
			fmt.Fscan(in, &a, &b, &c, &d, &h)
			s := getIndex(a, b)
			e := getIndex(c, d)
			ladders[i] = ladder{a, b, c, d, h, s, e}
		}
		for f, ids := range nodesByFloor {
			sort.Slice(ids, func(i, j int) bool { return cols[ids[i]] < cols[ids[j]] })
			nodesByFloor[f] = ids
		}
		edges := make([][]edgeE, len(floors))
		for _, ld := range ladders {
			edges[ld.s] = append(edges[ld.s], edgeE{ld.e, -ld.h})
		}
		const INF int64 = 1 << 60
		dist := make([]int64, len(floors))
		for i := range dist {
			dist[i] = INF
		}
		dist[startIdx] = 0
		for f := 1; f <= n; f++ {
			ids := nodesByFloor[f]
			if len(ids) == 0 {
				continue
			}
			for i := 1; i < len(ids); i++ {
				prev := ids[i-1]
				cur := ids[i]
				cost := int64(cols[cur]-cols[prev]) * x[f]
				if dist[prev]+cost < dist[cur] {
					dist[cur] = dist[prev] + cost
				}
			}
			for i := len(ids) - 2; i >= 0; i-- {
				next := ids[i+1]
				cur := ids[i]
				cost := int64(cols[next]-cols[cur]) * x[f]
				if dist[next]+cost < dist[cur] {
					dist[cur] = dist[next] + cost
				}
			}
			for _, id := range ids {
				if dist[id] == INF {
					continue
				}
				for _, e := range edges[id] {
					if dist[id]+e.w < dist[e.to] {
						dist[e.to] = dist[id] + e.w
					}
				}
			}
		}
		if dist[endIdx] >= INF/2 {
			out.WriteString("NO ESCAPE\n")
		} else {
			out.WriteString(fmt.Sprintf("%d\n", dist[endIdx]))
		}
	}
	return strings.TrimSpace(out.String())
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(5))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		m := rng.Intn(4) + 2
		k := rng.Intn(4)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for f := 1; f <= n; f++ {
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(10)+1))
			if f < n {
				sb.WriteByte(' ')
			} else {
				sb.WriteByte('\n')
			}
		}
		for j := 0; j < k; j++ {
			a := rng.Intn(n-1) + 1
			c := rng.Intn(n-a) + a + 1
			b := rng.Intn(m) + 1
			d := rng.Intn(m) + 1
			h := rng.Int63n(10) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", a, b, c, d, h))
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, t := range tests {
		expect := solveE(t)
		got, err := runProg(bin, t)
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, t, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
