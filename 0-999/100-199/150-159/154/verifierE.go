package main

import (
	"bytes"
	"container/heap"
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

type Point struct{ x, y float64 }

type Item struct {
	r   float64
	idx int
	ver int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].r > pq[j].r }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func sub(a, b Point) Point     { return Point{a.x - b.x, a.y - b.y} }
func cross(a, b Point) float64 { return a.x*b.y - a.y*b.x }
func dot(a, b Point) float64   { return a.x*b.x + a.y*b.y }
func absPt(a Point) float64    { return math.Hypot(a.x, a.y) }

func solveE(n int, R float64, pts []Point) float64 {
	if n == 1 {
		return 0.0
	}
	sort.Slice(pts, func(i, j int) bool {
		if math.Abs(pts[i].x-pts[j].x) > 1e-9 {
			return pts[i].x < pts[j].x
		}
		return pts[i].y < pts[j].y
	})
	q := make([]Point, 0, 2*n)
	for i := 0; i < n; i++ {
		for len(q) >= 2 && cross(sub(q[len(q)-1], q[len(q)-2]), sub(pts[i], q[len(q)-2])) <= 0 {
			q = q[:len(q)-1]
		}
		q = append(q, pts[i])
	}
	t := len(q)
	for i := n - 1; i >= 0; i-- {
		for len(q) > t && cross(sub(q[len(q)-1], q[len(q)-2]), sub(pts[i], q[len(q)-2])) <= 0 {
			q = q[:len(q)-1]
		}
		q = append(q, pts[i])
	}
	q = q[:len(q)-1]
	m := len(q)
	l := make([]int, m)
	r := make([]int, m)
	del := make([]bool, m)
	ver := make([]int, m)
	for i := 0; i < m; i++ {
		l[i] = (i - 1 + m) % m
		r[i] = (i + 1) % m
	}
	clk := 0
	pqh := &PriorityQueue{}
	heap.Init(pqh)
	delFunc := func(x int) { del[x] = true; r[l[x]] = r[x]; l[r[x]] = l[x] }
	update := func(x int) {
		if dot(sub(q[l[x]], q[x]), sub(q[r[x]], q[x])) > 0 {
			ver[x] = -1
			return
		}
		a, b, c := q[x], q[l[x]], q[r[x]]
		ab := absPt(sub(a, b))
		ac := absPt(sub(a, c))
		bc := absPt(sub(b, c))
		area2 := math.Abs(cross(sub(a, b), sub(c, b)))
		radius := ab * ac * bc / (2 * area2)
		clk++
		ver[x] = clk
		heap.Push(pqh, &Item{r: radius, idx: x, ver: clk})
	}
	for i := 0; i < m; i++ {
		if dot(sub(q[l[i]], q[i]), sub(q[r[i]], q[i])) < 0 {
			update(i)
		}
	}
	for pqh.Len() > 0 {
		top := (*pqh)[0]
		if top.r <= R {
			break
		}
		item := heap.Pop(pqh).(*Item)
		if item.ver != ver[item.idx] {
			continue
		}
		delFunc(item.idx)
		update(l[item.idx])
		update(r[item.idx])
	}
	s := 0
	for del[s] {
		s++
	}
	ans := 0.0
	i := s
	for {
		ni := r[i]
		ans += cross(q[i], q[ni]) / 2
		d := absPt(sub(q[i], q[ni]))
		v := d / (2 * R)
		if v > 1 {
			v = 1
		} else if v < -1 {
			v = -1
		}
		angle := 2 * math.Asin(v)
		ans += R * R * (angle - math.Sin(angle)) / 2
		i = ni
		if i == s {
			break
		}
	}
	return ans
}

func genCaseE(rng *rand.Rand) (int, float64, []Point) {
	n := rng.Intn(8) + 2
	R := rng.Float64()*5 + 1
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		pts[i] = Point{x: rng.Float64()*20 - 10, y: rng.Float64()*20 - 10}
	}
	return n, R, pts
}

func runCaseE(bin string, n int, R float64, pts []Point) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %.6f\n", n, R)
	for _, p := range pts {
		fmt.Fprintf(&input, "%.6f %.6f\n", p.x, p.y)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	exp := solveE(n, R, pts)
	got, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64)
	if err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	if math.Abs(got-exp) > 1e-6*math.Max(1, math.Abs(exp)) {
		return fmt.Errorf("expected %.6f got %.6f", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, R, pts := genCaseE(rng)
		if err := runCaseE(bin, n, R, pts); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
