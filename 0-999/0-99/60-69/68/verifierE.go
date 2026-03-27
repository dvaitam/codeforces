package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// --- embedded correct solver ---

type sTriangle struct {
	d   [3][3]int
	key [3]int
}

type sPoint struct {
	x float64
	y float64
}

var sTris [4]sTriangle
var sEdgeLen [12][12]int
var sBest int
var sMemo [5]map[string]bool

func sDist2(a, b sPoint) float64 {
	dx := a.x - b.x
	dy := a.y - b.y
	return dx*dx + dy*dy
}

func sCircleIntersections(a, b sPoint, ra2, rb2 float64) []sPoint {
	dx := b.x - a.x
	dy := b.y - a.y
	d2 := dx*dx + dy*dy
	if d2 < 1e-12 {
		return nil
	}
	d := math.Sqrt(d2)
	sum := math.Sqrt(ra2) + math.Sqrt(rb2)
	diff := math.Abs(math.Sqrt(ra2) - math.Sqrt(rb2))
	if d > sum+1e-8 || d < diff-1e-8 {
		return nil
	}
	x := (ra2 - rb2 + d2) / (2 * d)
	h2 := ra2 - x*x
	if h2 < -1e-8 {
		return nil
	}
	if h2 < 0 {
		h2 = 0
	}
	mx := a.x + x*dx/d
	my := a.y + x*dy/d
	if h2 <= 1e-9 {
		return []sPoint{{mx, my}}
	}
	h := math.Sqrt(h2)
	rx := -dy * h / d
	ry := dx * h / d
	p1 := sPoint{mx + rx, my + ry}
	p2 := sPoint{mx - rx, my - ry}
	if sDist2(p1, p2) < 1e-10 {
		return []sPoint{p1}
	}
	return []sPoint{p1, p2}
}

func sPlaceRest(edge [][]int, coords []sPoint, placed []bool, cnt int) bool {
	m := len(edge)
	if cnt == m {
		return true
	}

	sel := -1
	bestCnt := -1
	for v := 0; v < m; v++ {
		if placed[v] {
			continue
		}
		c := 0
		for u := 0; u < m; u++ {
			if placed[u] && edge[v][u] != -1 {
				c++
			}
		}
		if c > bestCnt {
			bestCnt = c
			sel = v
		}
	}
	if sel == -1 {
		return true
	}
	if bestCnt < 2 {
		return false
	}

	neigh := make([]int, 0)
	for u := 0; u < m; u++ {
		if placed[u] && edge[sel][u] != -1 {
			neigh = append(neigh, u)
		}
	}

	u0, u1 := -1, -1
	bestD := -1.0
	for i := 0; i < len(neigh); i++ {
		for j := i + 1; j < len(neigh); j++ {
			d := sDist2(coords[neigh[i]], coords[neigh[j]])
			if d > bestD {
				bestD = d
				u0 = neigh[i]
				u1 = neigh[j]
			}
		}
	}
	if u0 == -1 {
		return false
	}

	cands := sCircleIntersections(coords[u0], coords[u1], float64(edge[sel][u0]), float64(edge[sel][u1]))
	for _, p := range cands {
		ok := true
		for u := 0; u < m; u++ {
			if !placed[u] {
				continue
			}
			d2 := sDist2(p, coords[u])
			if d2 < 1e-8 {
				ok = false
				break
			}
			if edge[sel][u] != -1 && math.Abs(d2-float64(edge[sel][u])) > 1e-5 {
				ok = false
				break
			}
		}
		if ok {
			coords[sel] = p
			placed[sel] = true
			if sPlaceRest(edge, coords, placed, cnt+1) {
				return true
			}
			placed[sel] = false
		}
	}
	return false
}

func sSolveBlock(verts []int) bool {
	m := len(verts)
	if m <= 2 {
		return true
	}

	local := make([][]int, m)
	for i := 0; i < m; i++ {
		local[i] = make([]int, m)
		for j := 0; j < m; j++ {
			local[i][j] = -1
		}
	}
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			local[i][j] = sEdgeLen[verts[i]][verts[j]]
			local[j][i] = local[i][j]
		}
	}

	type Triple struct{ a, b, c int }
	triples := make([]Triple, 0)
	for i := 0; i < m; i++ {
		for j := i + 1; j < m; j++ {
			if local[i][j] == -1 {
				continue
			}
			for k := j + 1; k < m; k++ {
				if local[i][k] != -1 && local[j][k] != -1 {
					triples = append(triples, Triple{i, j, k})
				}
			}
		}
	}
	if len(triples) == 0 {
		return false
	}

	for _, tr := range triples {
		ab2 := float64(local[tr.a][tr.b])
		ac2 := float64(local[tr.a][tr.c])
		bc2 := float64(local[tr.b][tr.c])
		ab := math.Sqrt(ab2)
		x := (ac2 - bc2 + ab2) / (2 * ab)
		y2 := ac2 - x*x
		if y2 <= 1e-9 {
			continue
		}
		y := math.Sqrt(y2)

		coords := make([]sPoint, m)
		placed := make([]bool, m)
		coords[tr.a] = sPoint{0, 0}
		coords[tr.b] = sPoint{ab, 0}
		coords[tr.c] = sPoint{x, y}
		placed[tr.a] = true
		placed[tr.b] = true
		placed[tr.c] = true

		if sPlaceRest(local, coords, placed, 3) {
			return true
		}
	}
	return false
}

func sIsRealizable(n int) bool {
	adj := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if sEdgeLen[i][j] != -1 {
				adj[i] = append(adj[i], j)
				adj[j] = append(adj[j], i)
			}
		}
	}

	disc := make([]int, n)
	low := make([]int, n)
	timer := 0
	type Edge struct{ u, v int }
	stack := make([]Edge, 0)
	blocks := make([][]int, 0)

	var dfs func(int, int)
	dfs = func(u, parent int) {
		timer++
		disc[u] = timer
		low[u] = timer
		for _, v := range adj[u] {
			if disc[v] == 0 {
				stack = append(stack, Edge{u, v})
				dfs(v, u)
				if low[v] < low[u] {
					low[u] = low[v]
				}
				if low[v] >= disc[u] {
					used := make([]bool, n)
					comp := make([]int, 0)
					for len(stack) > 0 {
						e := stack[len(stack)-1]
						stack = stack[:len(stack)-1]
						if !used[e.u] {
							used[e.u] = true
							comp = append(comp, e.u)
						}
						if !used[e.v] {
							used[e.v] = true
							comp = append(comp, e.v)
						}
						if (e.u == u && e.v == v) || (e.u == v && e.v == u) {
							break
						}
					}
					blocks = append(blocks, comp)
				}
			} else if v != parent && disc[v] < disc[u] {
				stack = append(stack, Edge{u, v})
				if disc[v] < low[u] {
					low[u] = disc[v]
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		if disc[i] == 0 {
			dfs(i, -1)
			if len(stack) > 0 {
				used := make([]bool, n)
				comp := make([]int, 0)
				for len(stack) > 0 {
					e := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					if !used[e.u] {
						used[e.u] = true
						comp = append(comp, e.u)
					}
					if !used[e.v] {
						used[e.v] = true
						comp = append(comp, e.v)
					}
				}
				blocks = append(blocks, comp)
			}
		}
	}

	for _, b := range blocks {
		if !sSolveBlock(b) {
			return false
		}
	}
	return true
}

func sMakeSig(idx, n int) string {
	buf := make([]byte, 0, 300)
	buf = strconv.AppendInt(buf, int64(idx), 10)
	buf = append(buf, '|')
	buf = strconv.AppendInt(buf, int64(n), 10)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			buf = append(buf, ',')
			if sEdgeLen[i][j] == -1 {
				buf = append(buf, '_')
			} else {
				buf = strconv.AppendInt(buf, int64(sEdgeLen[i][j]), 10)
			}
		}
	}
	return string(buf)
}

func sSearch(idx, n int) {
	if n >= sBest || sBest == 3 {
		return
	}
	sig := sMakeSig(idx, n)
	if sMemo[idx] == nil {
		sMemo[idx] = make(map[string]bool)
	}
	if sMemo[idx][sig] {
		return
	}
	sMemo[idx][sig] = true

	if idx == 4 {
		if sIsRealizable(n) {
			sBest = n
		}
		return
	}

	var mp [3]int
	sAssignTriangle(idx, 0, n, 0, &mp)
}

func sAssignTriangle(ti, pos, n, newUsed int, mp *[3]int) {
	if n+newUsed >= sBest {
		return
	}
	if pos == 3 {
		sSearch(ti+1, n+newUsed)
		return
	}

	maxID := n + newUsed
	for g := 0; g <= maxID; g++ {
		ok := true
		for i := 0; i < pos; i++ {
			if mp[i] == g {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}

		added := make([][2]int, 0, 3)
		for i := 0; i < pos; i++ {
			a, b := g, mp[i]
			need := sTris[ti].d[pos][i]
			if sEdgeLen[a][b] == -1 {
				sEdgeLen[a][b] = need
				sEdgeLen[b][a] = need
				added = append(added, [2]int{a, b})
			} else if sEdgeLen[a][b] != need {
				ok = false
				break
			}
		}
		if ok {
			mp[pos] = g
			nu := newUsed
			if g == maxID {
				nu++
			}
			sAssignTriangle(ti, pos+1, n, nu, mp)
		}
		for _, p := range added {
			sEdgeLen[p[0]][p[1]] = -1
			sEdgeLen[p[1]][p[0]] = -1
		}
	}
}

func solveCase(input string) string {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	lineIdx := 0
	readLine := func() string {
		s := lines[lineIdx]
		lineIdx++
		return s
	}

	for i := 0; i < 4; i++ {
		fields := strings.Fields(readLine())
		var x [3]int
		var y [3]int
		x[0], _ = strconv.Atoi(fields[0])
		y[0], _ = strconv.Atoi(fields[1])
		x[1], _ = strconv.Atoi(fields[2])
		y[1], _ = strconv.Atoi(fields[3])
		x[2], _ = strconv.Atoi(fields[4])
		y[2], _ = strconv.Atoi(fields[5])
		var t sTriangle
		for a := 0; a < 3; a++ {
			for b := a + 1; b < 3; b++ {
				dx := x[a] - x[b]
				dy := y[a] - y[b]
				v := dx*dx + dy*dy
				t.d[a][b] = v
				t.d[b][a] = v
			}
		}
		k := []int{t.d[0][1], t.d[0][2], t.d[1][2]}
		sort.Ints(k)
		t.key = [3]int{k[0], k[1], k[2]}
		sTris[i] = t
	}
	sort.Slice(sTris[:], func(i, j int) bool {
		for k := 0; k < 3; k++ {
			if sTris[i].key[k] != sTris[j].key[k] {
				return sTris[i].key[k] < sTris[j].key[k]
			}
		}
		return false
	})

	for i := 0; i < 12; i++ {
		for j := 0; j < 12; j++ {
			sEdgeLen[i][j] = -1
		}
	}

	sEdgeLen[0][1] = sTris[0].d[0][1]
	sEdgeLen[1][0] = sTris[0].d[0][1]
	sEdgeLen[0][2] = sTris[0].d[0][2]
	sEdgeLen[2][0] = sTris[0].d[0][2]
	sEdgeLen[1][2] = sTris[0].d[1][2]
	sEdgeLen[2][1] = sTris[0].d[1][2]

	sBest = 12
	for i := range sMemo {
		sMemo[i] = nil
	}
	sSearch(1, 3)
	return fmt.Sprintf("%d", sBest)
}

// --- verifier infrastructure ---

func nonDegenerate(x1, y1, x2, y2, x3, y3 int) bool {
	return (x2-x1)*(y3-y1)-(y2-y1)*(x3-x1) != 0
}

func generateCase(rng *rand.Rand) string {
	var sb strings.Builder
	for i := 0; i < 4; i++ {
		for {
			x1 := rng.Intn(21)
			y1 := rng.Intn(21)
			x2 := rng.Intn(21)
			y2 := rng.Intn(21)
			x3 := rng.Intn(21)
			y3 := rng.Intn(21)
			if nonDegenerate(x1, y1, x2, y2, x3, y3) {
				sb.WriteString(fmt.Sprintf("%d %d %d %d %d %d\n", x1, y1, x2, y2, x3, y3))
				break
			}
		}
	}
	return sb.String()
}

func runCase(bin, input string) error {
	expected := solveCase(input)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
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
		in := generateCase(rng)
		if err := runCase(bin, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
