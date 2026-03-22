package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

// Embedded reference solver from the accepted solution.

type refCar struct {
	x, y   int64
	dx, dy int64
	s      int64
	xf, yf float64
	slope  float64
	interc float64
	reach  float64
	vx, vy float64
}

type refTempSeg struct {
	lim          float64
	leftX, rightX float64
}

type refEvent struct {
	x, y float64
	typ  int
	id   int
}

type refNode struct {
	id          int
	pr          uint32
	left, right *refNode
}

func refSolve(input string) string {
	rd := strings.NewReader(input)
	var n int
	fmt.Fscan(rd, &n)
	if n < 2 {
		return "No show :("
	}

	cars := make([]refCar, n)
	segs := make([]refTempSeg, n)
	events := make([]refEvent, 2*n)
	nodes := make([]refNode, n)

	seed := uint32(2463534242)
	rnd := func() uint32 {
		seed ^= seed << 13
		seed ^= seed >> 17
		seed ^= seed << 5
		return seed
	}

	for i := 0; i < n; i++ {
		var x, y, dx, dy, s int64
		fmt.Fscan(rd, &x, &y, &dx, &dy, &s)
		lenv := math.Hypot(float64(dx), float64(dy))
		slope := float64(dy) / float64(dx)
		xf := float64(x)
		yf := float64(y)
		reach := float64(s) / lenv
		vx := float64(dx) * reach
		vy := float64(dy) * reach
		cars[i] = refCar{
			x: x, y: y, dx: dx, dy: dy, s: s,
			xf: xf, yf: yf, slope: slope,
			interc: yf - slope*xf, reach: reach,
			vx: vx, vy: vy,
		}
		nodes[i].id = i
		nodes[i].pr = rnd()
	}

	var curX float64

	less := func(i, j int) bool {
		if i == j {
			return false
		}
		ci := cars[i]
		cj := cars[j]
		yi := ci.slope*curX + ci.interc
		yj := cj.slope*curX + cj.interc
		if yi < yj {
			return true
		}
		if yi > yj {
			return false
		}
		if ci.slope < cj.slope {
			return true
		}
		if ci.slope > cj.slope {
			return false
		}
		if ci.interc < cj.interc {
			return true
		}
		if ci.interc > cj.interc {
			return false
		}
		return i < j
	}

	var rotateRight func(*refNode) *refNode
	rotateRight = func(p *refNode) *refNode {
		q := p.left
		p.left = q.right
		q.right = p
		return q
	}
	var rotateLeft func(*refNode) *refNode
	rotateLeft = func(p *refNode) *refNode {
		q := p.right
		p.right = q.left
		q.left = p
		return q
	}

	var insertNode func(*refNode, *refNode) *refNode
	insertNode = func(root *refNode, node *refNode) *refNode {
		if root == nil {
			return node
		}
		if less(node.id, root.id) {
			root.left = insertNode(root.left, node)
			if root.left.pr < root.pr {
				root = rotateRight(root)
			}
		} else {
			root.right = insertNode(root.right, node)
			if root.right.pr < root.pr {
				root = rotateLeft(root)
			}
		}
		return root
	}

	var mergeNodes func(*refNode, *refNode) *refNode
	mergeNodes = func(a, b *refNode) *refNode {
		if a == nil {
			return b
		}
		if b == nil {
			return a
		}
		if a.pr < b.pr {
			a.right = mergeNodes(a.right, b)
			return a
		}
		b.left = mergeNodes(a, b.left)
		return b
	}

	var eraseNode func(*refNode, int) *refNode
	eraseNode = func(root *refNode, id int) *refNode {
		if root == nil {
			return nil
		}
		if root.id == id {
			return mergeNodes(root.left, root.right)
		}
		if less(id, root.id) {
			root.left = eraseNode(root.left, id)
		} else {
			root.right = eraseNode(root.right, id)
		}
		return root
	}

	minNode := func(root *refNode) int {
		for root.left != nil {
			root = root.left
		}
		return root.id
	}
	maxNode := func(root *refNode) int {
		for root.right != nil {
			root = root.right
		}
		return root.id
	}

	findPredSucc := func(root *refNode, id int) (int, int) {
		pred, succ := -1, -1
		cur := root
		for cur != nil {
			if less(id, cur.id) {
				succ = cur.id
				cur = cur.left
			} else if less(cur.id, id) {
				pred = cur.id
				cur = cur.right
			} else {
				if cur.left != nil {
					pred = maxNode(cur.left)
				}
				if cur.right != nil {
					succ = minNode(cur.right)
				}
				break
			}
		}
		return pred, succ
	}

	intersects := func(i, j int) bool {
		ci := cars[i]
		cj := cars[j]
		den := ci.dx*cj.dy - ci.dy*cj.dx
		diffx := cj.x - ci.x
		diffy := cj.y - ci.y
		if den != 0 {
			numT := diffx*cj.dy - diffy*cj.dx
			numU := diffx*ci.dy - diffy*ci.dx
			if den < 0 {
				den = -den
				numT = -numT
				numU = -numU
			}
			if numT < 0 || numU < 0 {
				return false
			}
			t := float64(numT) / float64(den)
			u := float64(numU) / float64(den)
			limI := segs[i].lim
			limJ := segs[j].lim
			epsI := 1e-12 * math.Max(1.0, math.Max(t, limI))
			epsJ := 1e-12 * math.Max(1.0, math.Max(u, limJ))
			return t <= limI+epsI && u <= limJ+epsJ
		}
		if diffx*ci.dy-diffy*ci.dx != 0 {
			return false
		}
		l := segs[i].leftX
		if segs[j].leftX > l {
			l = segs[j].leftX
		}
		r := segs[i].rightX
		if segs[j].rightX < r {
			r = segs[j].rightX
		}
		return l <= r+1e-9
	}

	feasible := func(T float64) bool {
		if T <= 0 {
			return false
		}
		m := 2 * n
		for i := 0; i < n; i++ {
			c := cars[i]
			lim := c.reach * T
			qx := c.xf + c.vx*T
			qy := c.yf + c.vy*T
			if qx < c.xf {
				segs[i].lim = lim
				segs[i].leftX = qx
				segs[i].rightX = c.xf
				events[2*i] = refEvent{x: qx, y: qy, typ: 0, id: i}
				events[2*i+1] = refEvent{x: c.xf, y: c.yf, typ: 1, id: i}
			} else {
				segs[i].lim = lim
				segs[i].leftX = c.xf
				segs[i].rightX = qx
				events[2*i] = refEvent{x: c.xf, y: c.yf, typ: 0, id: i}
				events[2*i+1] = refEvent{x: qx, y: qy, typ: 1, id: i}
			}
			nodes[i].left = nil
			nodes[i].right = nil
		}

		sortEvents := events[:m]
		// sort
		for i := 1; i < len(sortEvents); i++ {
			key := sortEvents[i]
			j := i - 1
			for j >= 0 {
				e := sortEvents[j]
				swap := false
				if e.x > key.x {
					swap = true
				} else if e.x == key.x {
					if e.y > key.y {
						swap = true
					} else if e.y == key.y {
						if e.typ > key.typ {
							swap = true
						} else if e.typ == key.typ {
							if e.id > key.id {
								swap = true
							}
						}
					}
				}
				if swap {
					sortEvents[j+1] = sortEvents[j]
					j--
				} else {
					break
				}
			}
			sortEvents[j+1] = key
		}

		var root *refNode
		starts := make([]int, 0, n)
		ends := make([]int, 0, n)

		for l := 0; l < m; {
			r := l + 1
			x := events[l].x
			for r < m && events[r].x == x {
				r++
			}

			starts = starts[:0]
			ends = ends[:0]
			for k := l; k < r; k++ {
				if events[k].typ == 0 {
					starts = append(starts, events[k].id)
				} else {
					ends = append(ends, events[k].id)
				}
			}

			curX = math.Nextafter(x, math.Inf(-1))
			for _, id := range starts {
				pred, succ := findPredSucc(root, id)
				if pred != -1 && intersects(id, pred) {
					return true
				}
				if succ != -1 && intersects(id, succ) {
					return true
				}
			}

			for _, id := range ends {
				pred, succ := findPredSucc(root, id)
				if pred != -1 && intersects(id, pred) {
					return true
				}
				if succ != -1 && intersects(id, succ) {
					return true
				}
				root = eraseNode(root, id)
				if pred != -1 && succ != -1 && intersects(pred, succ) {
					return true
				}
			}

			curX = math.Nextafter(x, math.Inf(1))
			for _, id := range starts {
				pred, succ := findPredSucc(root, id)
				root = insertNode(root, &nodes[id])
				if pred != -1 && intersects(id, pred) {
					return true
				}
				if succ != -1 && intersects(id, succ) {
					return true
				}
			}

			l = r
		}

		return false
	}

	const HI = 6000000000.0
	if !feasible(HI) {
		return "No show :("
	}

	lo, hi := 0.0, HI
	for it := 0; it < 70; it++ {
		mid := (lo + hi) * 0.5
		if feasible(mid) {
			hi = mid
		} else {
			lo = mid
		}
	}

	return fmt.Sprintf("%.15f", hi)
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
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

type testCase struct {
	n    int
	vals []string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(10) + 1
	vals := make([]string, n)
	for i := 0; i < n; i++ {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		dx := rng.Intn(10) + 1
		if rng.Intn(2) == 0 {
			dx = -dx
		}
		dy := rng.Intn(10) + 1
		if rng.Intn(2) == 0 {
			dy = -dy
		}
		s := rng.Intn(10) + 1
		vals[i] = fmt.Sprintf("%d %d %d %d %d", x, y, dx, dy, s)
	}
	return testCase{n, vals}
}

func compareFloats(a, b float64) bool {
	diff := math.Abs(a - b)
	denom := math.Max(1.0, math.Abs(a))
	return diff/denom <= 1e-6
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for _, v := range tc.vals {
			sb.WriteString(v + "\n")
		}
		input := sb.String()
		expStr := refSolve(input)
		gotStr, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if expStr == "No show :(" {
			if gotStr != expStr {
				fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expStr, gotStr, input)
				os.Exit(1)
			}
			continue
		}
		expF, _ := strconv.ParseFloat(expStr, 64)
		gotF, err := strconv.ParseFloat(gotStr, 64)
		if err != nil || !compareFloats(expF, gotF) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, expStr, gotStr, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
