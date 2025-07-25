package main

import (
	"bytes"
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Segment struct {
	start, end float64
	from, to   int
	speed      float64
}

type NodeEvent struct {
	time float64
	node int
}

func path(n int, edges [][]int, start, end int) []int {
	if start == end {
		return []int{start}
	}
	parent := make([]int, n+1)
	for i := range parent {
		parent[i] = -1
	}
	q := list.New()
	q.PushBack(start)
	parent[start] = 0
	for q.Len() > 0 {
		v := q.Remove(q.Front()).(int)
		if v == end {
			break
		}
		for _, to := range edges[v] {
			if parent[to] == -1 {
				parent[to] = v
				q.PushBack(to)
			}
		}
	}
	var res []int
	cur := end
	for cur != 0 {
		res = append(res, cur)
		if cur == start {
			break
		}
		cur = parent[cur]
	}
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}

func solveE(n, m int, edges [][]int, suits [][4]int) string {
	segments := make([][]Segment, m)
	nodeEvents := make([][]NodeEvent, m)
	for i, s := range suits {
		t0 := float64(s[0])
		c := float64(s[1])
		v := s[2]
		u := s[3]
		p := path(n, edges, v, u)
		nodeEvents[i] = append(nodeEvents[i], NodeEvent{t0, v})
		for j := 0; j < len(p)-1; j++ {
			seg := Segment{start: t0, end: t0 + 1.0/c, from: p[j], to: p[j+1], speed: c}
			segments[i] = append(segments[i], seg)
			t0 += 1.0 / c
			nodeEvents[i] = append(nodeEvents[i], NodeEvent{t0, p[j+1]})
		}
	}
	ans := math.Inf(1)
	eps := 1e-9
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			ei := nodeEvents[i]
			ej := nodeEvents[j]
			for _, a := range ei {
				for _, b := range ej {
					if a.node == b.node && math.Abs(a.time-b.time) < eps && a.time < ans {
						ans = a.time
					}
				}
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			for _, s1 := range segments[i] {
				for _, s2 := range segments[j] {
					t, ok := collideSeg(s1, s2)
					if ok && t < ans {
						ans = t
					}
				}
			}
		}
	}
	if math.IsInf(ans, 1) {
		return "-1"
	}
	return fmt.Sprintf("%.6f", ans)
}

func collideSeg(a, b Segment) (float64, bool) {
	eps := 1e-9
	if a.from == b.from && a.to == b.to {
		if math.Abs(a.speed-b.speed) < eps {
			if math.Abs(a.start-b.start) < eps {
				t := math.Max(a.start, b.start)
				if t <= a.end+eps && t <= b.end+eps {
					return t, true
				}
			}
			return 0, false
		}
		t := (a.speed*a.start - b.speed*b.start) / (a.speed - b.speed)
		if t+eps < math.Max(a.start, b.start) || t-eps > math.Min(a.end, b.end) {
			return 0, false
		}
		return t, true
	}
	if a.from == b.to && a.to == b.from {
		t := (1 + a.speed*a.start + b.speed*b.start) / (a.speed + b.speed)
		if t+eps < math.Max(a.start, b.start) || t-eps > math.Min(a.end, b.end) {
			return 0, false
		}
		return t, true
	}
	return 0, false
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	edges := make([][]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i] = append(edges[i], p)
		edges[p] = append(edges[p], i)
	}
	m := rng.Intn(3) + 2
	suits := make([][4]int, m)
	for i := 0; i < m; i++ {
		ti := rng.Intn(6)
		ci := rng.Intn(5) + 1
		vi := rng.Intn(n) + 1
		ui := rng.Intn(n) + 1
		suits[i] = [4]int{ti, ci, vi, ui}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 2; i <= n; i++ {
		for _, v := range edges[i] {
			if v < i {
				sb.WriteString(fmt.Sprintf("%d %d\n", i, v))
			}
		}
	}
	for i := 0; i < m; i++ {
		s := suits[i]
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", s[0], s[1], s[2], s[3]))
	}
	expected := solveE(n, m, edges, suits)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if expected == "-1" {
		if outStr != "-1" {
			return fmt.Errorf("expected -1 got %s", outStr)
		}
		return nil
	}
	valOut, err1 := strconv.ParseFloat(outStr, 64)
	valExp, err2 := strconv.ParseFloat(expected, 64)
	if err1 != nil || err2 != nil {
		return fmt.Errorf("invalid float output")
	}
	if math.Abs(valOut-valExp) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", valExp, valOut)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
