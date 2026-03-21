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

// Embedded reference solver for 856E

type refSatellite struct {
	id   int
	x, y int64
}

type refNode struct {
	top [3]int
	cnt int
}

type refSolver struct {
	R          int64
	sat        []refSatellite
	sortedSats []int
	leafIdx    []int
	tree       []refNode
}

func (s *refSolver) uGreater(id1, id2 int) bool {
	p1 := (s.sat[id1].x + s.R) * s.sat[id2].y
	p2 := (s.sat[id2].x + s.R) * s.sat[id1].y
	return p1 > p2
}

func (s *refSolver) vGreater(id1, id2 int) bool {
	p1 := (s.sat[id1].x - s.R) * s.sat[id2].y
	p2 := (s.sat[id2].x - s.R) * s.sat[id1].y
	if p1 == p2 {
		return id1 > id2
	}
	return p1 > p2
}

func (s *refSolver) vGreaterEq(id1, id2 int) bool {
	p1 := (s.sat[id1].x - s.R) * s.sat[id2].y
	p2 := (s.sat[id2].x - s.R) * s.sat[id1].y
	return p1 >= p2
}

func (s *refSolver) validGeo(u_ref, v_ref int) bool {
	xu, yu := s.sat[u_ref].x, s.sat[u_ref].y
	xv, yv := s.sat[v_ref].x, s.sat[v_ref].y
	return (xu+s.R)*(xv-s.R) >= -yu*yv
}

func (s *refSolver) merge(a, b refNode) refNode {
	var res refNode
	i, j := 0, 0
	for i < a.cnt && j < b.cnt && res.cnt < 3 {
		if s.vGreater(a.top[i], b.top[j]) {
			res.top[res.cnt] = a.top[i]
			res.cnt++
			i++
		} else {
			res.top[res.cnt] = b.top[j]
			res.cnt++
			j++
		}
	}
	for i < a.cnt && res.cnt < 3 {
		res.top[res.cnt] = a.top[i]
		res.cnt++
		i++
	}
	for j < b.cnt && res.cnt < 3 {
		res.top[res.cnt] = b.top[j]
		res.cnt++
		j++
	}
	return res
}

func (s *refSolver) update(node, l, r, idx, id, op int) {
	if l == r {
		if op == 1 {
			s.tree[node].top[0] = id
			s.tree[node].cnt = 1
		} else {
			s.tree[node].cnt = 0
		}
		return
	}
	mid := (l + r) / 2
	if idx <= mid {
		s.update(2*node, l, mid, idx, id, op)
	} else {
		s.update(2*node+1, mid+1, r, idx, id, op)
	}
	s.tree[node] = s.merge(s.tree[2*node], s.tree[2*node+1])
}

func (s *refSolver) queryTree(node, l, r, ql, qr int) refNode {
	if ql <= l && r <= qr {
		return s.tree[node]
	}
	mid := (l + r) / 2
	if qr <= mid {
		return s.queryTree(2*node, l, mid, ql, qr)
	}
	if ql > mid {
		return s.queryTree(2*node+1, mid+1, r, ql, qr)
	}
	return s.merge(
		s.queryTree(2*node, l, mid, ql, qr),
		s.queryTree(2*node+1, mid+1, r, ql, qr),
	)
}

func (s *refSolver) findRightmostU(u_ref, K int) int {
	l, r_idx := 0, K-1
	ans := -1
	for l <= r_idx {
		mid := (l + r_idx) / 2
		id := s.sortedSats[mid]
		p1 := (s.sat[id].x + s.R) * s.sat[u_ref].y
		p2 := (s.sat[u_ref].x + s.R) * s.sat[id].y
		if p1 <= p2 {
			ans = mid
			l = mid + 1
		} else {
			r_idx = mid - 1
		}
	}
	return ans
}

type refEvent struct {
	typ  int
	a, b int
	id   int
}

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))

	readInt := func() int {
		var n int
		var sign int = 1
		for {
			b, err := reader.ReadByte()
			if err != nil {
				return 0
			}
			if b == '-' {
				sign = -1
				continue
			}
			if b >= '0' && b <= '9' {
				n = int(b - '0')
				break
			}
		}
		for {
			b, err := reader.ReadByte()
			if err != nil || b < '0' || b > '9' {
				return n * sign
			}
			n = n*10 + int(b-'0')
		}
	}

	r_planet := readInt()
	nEvents := readInt()
	s := &refSolver{R: int64(r_planet)}

	events := make([]refEvent, nEvents)
	K := 0
	curId := 0
	for i := 0; i < nEvents; i++ {
		typ := readInt()
		events[i].typ = typ
		if typ == 1 {
			events[i].a = readInt()
			events[i].b = readInt()
			curId++
			events[i].id = curId
			K++
		} else if typ == 2 {
			events[i].a = readInt()
		} else {
			events[i].a = readInt()
			events[i].b = readInt()
		}
	}

	if K == 0 {
		return ""
	}

	s.sat = make([]refSatellite, K+1)
	for i := 0; i < nEvents; i++ {
		if events[i].typ == 1 {
			id := events[i].id
			s.sat[id] = refSatellite{id: id, x: int64(events[i].a), y: int64(events[i].b)}
		}
	}

	s.sortedSats = make([]int, K)
	for i := 0; i < K; i++ {
		s.sortedSats[i] = i + 1
	}
	sort.Slice(s.sortedSats, func(i, j int) bool {
		id1, id2 := s.sortedSats[i], s.sortedSats[j]
		p1 := (s.sat[id1].x + s.R) * s.sat[id2].y
		p2 := (s.sat[id2].x + s.R) * s.sat[id1].y
		if p1 == p2 {
			return id1 < id2
		}
		return p1 < p2
	})

	s.leafIdx = make([]int, K+1)
	for i, id := range s.sortedSats {
		s.leafIdx[id] = i
	}

	s.tree = make([]refNode, 4*K)

	var result strings.Builder
	for i := 0; i < nEvents; i++ {
		ev := events[i]
		if ev.typ == 1 {
			id := ev.id
			idx := s.leafIdx[id]
			s.update(1, 0, K-1, idx, id, 1)
		} else if ev.typ == 2 {
			id := ev.a
			idx := s.leafIdx[id]
			s.update(1, 0, K-1, idx, id, 0)
		} else if ev.typ == 3 {
			i_id, j_id := ev.a, ev.b

			u_ref := i_id
			if s.uGreater(j_id, i_id) {
				u_ref = j_id
			}

			v_ref := i_id
			if s.vGreater(i_id, j_id) {
				v_ref = j_id
			}

			if !s.validGeo(u_ref, v_ref) {
				result.WriteString("NO\n")
				continue
			}

			R_idx := s.findRightmostU(u_ref, K)
			if R_idx == -1 {
				result.WriteString("NO\n")
				continue
			}

			res := s.queryTree(1, 0, K-1, 0, R_idx)

			interference := false
			for idx := 0; idx < res.cnt; idx++ {
				k := res.top[idx]
				if k == i_id || k == j_id {
					continue
				}
				if s.vGreaterEq(k, v_ref) {
					interference = true
					break
				}
			}

			if interference {
				result.WriteString("NO\n")
			} else {
				result.WriteString("YES\n")
			}
		}
	}
	return strings.TrimSpace(result.String())
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func generateCaseE(rng *rand.Rand) string {
	r := rng.Intn(5) + 1
	n := rng.Intn(10) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", r, n)
	nextID := 1
	active := make([]int, 0)
	coords := map[[2]int]bool{}
	for i := 0; i < n; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 || len(active) < 2 {
			typ = 1
		}
		if typ == 1 {
			x := rng.Intn(10) + r + 1
			y := rng.Intn(10) + 1
			for coords[[2]int{x, y}] {
				x = rng.Intn(10) + r + 1
				y = rng.Intn(10) + 1
			}
			coords[[2]int{x, y}] = true
			fmt.Fprintf(&sb, "1 %d %d\n", x, y)
			active = append(active, nextID)
			nextID++
		} else if typ == 2 {
			idx := rng.Intn(len(active))
			id := active[idx]
			active = append(active[:idx], active[idx+1:]...)
			fmt.Fprintf(&sb, "2 %d\n", id)
		} else {
			i1 := rng.Intn(len(active))
			j1 := rng.Intn(len(active) - 1)
			if j1 >= i1 {
				j1++
			}
			fmt.Fprintf(&sb, "3 %d %d\n", active[i1], active[j1])
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		expected := solveE(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\ninput:\n%s", i+1, err, got, tc)
			os.Exit(1)
		}
		gotTrimmed := strings.TrimSpace(got)
		if gotTrimmed != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expected, gotTrimmed, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
