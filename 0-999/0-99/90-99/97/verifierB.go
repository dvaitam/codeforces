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
)

const (
	maxCoord  = 1_000_000_000
	maxPoints = 200000
	infCoord  = 2_000_000_010
)

type point struct {
	x, y int
}

type testCase struct {
	input string
	base  []point
}

type node struct {
	y           int
	idx         int
	x           int
	nextGreater int
	nextLower   int
	left, right *node
	priority    int
}

type treap struct {
	root *node
	rnd  *rand.Rand
	idx  int
}

func newTreap() *treap {
	return &treap{
		rnd: rand.New(rand.NewSource(7)),
	}
}

func lessNode(a, b *node) bool {
	if a.y != b.y {
		return a.y < b.y
	}
	return a.idx < b.idx
}

func rotateLeft(root *node) *node {
	r := root.right
	root.right = r.left
	r.left = root
	return r
}

func rotateRight(root *node) *node {
	l := root.left
	root.left = l.right
	l.right = root
	return l
}

func insertNode(root, nd *node) *node {
	if root == nil {
		return nd
	}
	if lessNode(nd, root) {
		root.left = insertNode(root.left, nd)
		if root.left.priority < root.priority {
			root = rotateRight(root)
		}
	} else {
		root.right = insertNode(root.right, nd)
		if root.right != nil && root.right.priority < root.priority {
			root = rotateLeft(root)
		}
	}
	return root
}

func (t *treap) insert(p point) {
	t.idx++
	nd := &node{
		y:           p.y,
		idx:         t.idx,
		x:           p.x,
		nextGreater: infCoord,
		nextLower:   -infCoord,
		priority:    t.rnd.Int(),
	}
	t.root = insertNode(t.root, nd)
}

func (t *treap) findPred(y int) *node {
	cur := t.root
	var res *node
	for cur != nil {
		if cur.y < y {
			if res == nil || lessNode(res, cur) {
				res = cur
			}
			cur = cur.right
		} else {
			cur = cur.left
		}
	}
	return res
}

func (t *treap) findSucc(y int) *node {
	cur := t.root
	var res *node
	for cur != nil {
		if cur.y > y {
			if res == nil || lessNode(cur, res) {
				res = cur
			}
			cur = cur.left
		} else {
			cur = cur.right
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()

	for i, tc := range tests {
		if err := ensureReferencePasses(refBin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "reference solution failed on test %d: %v\nInput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\nInput:\n%sOutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
		if err := validateOutput(tc, out); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%sOutput:\n%s\n", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func ensureReferencePasses(ref string, tc testCase) error {
	out, err := runProgram(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, out)
	}
	if err := validateOutput(tc, out); err != nil {
		return fmt.Errorf("reference produced invalid output: %v\n%s", err, out)
	}
	return nil
}

func buildReference() (string, error) {
	path := "./ref97B.bin"
	cmd := exec.Command("go", "build", "-o", path, "97B.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func validateOutput(tc testCase, out string) error {
	points, err := parseOutput(out, len(tc.base))
	if err != nil {
		return err
	}
	if err := ensureSuperset(points, tc.base); err != nil {
		return err
	}
	if err := checkGood(points); err != nil {
		return err
	}
	return nil
}

func parseOutput(out string, minPoints int) ([]point, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return nil, fmt.Errorf("failed to read number of points: %v", err)
	}
	if m < minPoints {
		return nil, fmt.Errorf("reported %d points but need at least %d", m, minPoints)
	}
	if m > maxPoints {
		return nil, fmt.Errorf("reported %d points which exceeds limit %d", m, maxPoints)
	}
	points := make([]point, 0, m)
	for i := 0; i < m; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return nil, fmt.Errorf("failed to read point %d: %v", i+1, err)
		}
		if abs(x) > maxCoord || abs(y) > maxCoord {
			return nil, fmt.Errorf("point %d has coordinates out of range: (%d,%d)", i+1, x, y)
		}
		points = append(points, point{x, y})
	}
	return points, nil
}

func ensureSuperset(points, base []point) error {
	seen := make(map[[2]int]struct{}, len(points))
	for _, p := range points {
		key := [2]int{p.x, p.y}
		if _, ok := seen[key]; ok {
			return fmt.Errorf("duplicate point (%d,%d) in output", p.x, p.y)
		}
		seen[key] = struct{}{}
	}
	for _, p := range base {
		key := [2]int{p.x, p.y}
		if _, ok := seen[key]; !ok {
			return fmt.Errorf("missing original point (%d,%d)", p.x, p.y)
		}
	}
	return nil
}

func checkGood(points []point) error {
	if len(points) == 0 {
		return fmt.Errorf("no points provided")
	}
	order := append([]point(nil), points...)
	columns := make(map[int][]int)
	for _, p := range points {
		columns[p.x] = append(columns[p.x], p.y)
	}
	for x := range columns {
		sort.Ints(columns[x])
	}
	sort.Slice(order, func(i, j int) bool {
		if order[i].x != order[j].x {
			return order[i].x < order[j].x
		}
		if order[i].y != order[j].y {
			return order[i].y < order[j].y
		}
		return i < j
	})
	t := newTreap()
	for _, p := range order {
		if pred := t.findPred(p.y); pred != nil {
			if pred.x != p.x && pred.nextGreater > p.y {
				if !hasVerticalBridge(columns[pred.x], pred.y, p.y, pred.y) &&
					!hasVerticalBridge(columns[p.x], pred.y, p.y, p.y) {
					return fmt.Errorf("points (%d,%d) and (%d,%d) need an intermediate point", pred.x, pred.y, p.x, p.y)
				}
			}
			if p.y < pred.nextGreater {
				pred.nextGreater = p.y
			}
		}
		if succ := t.findSucc(p.y); succ != nil {
			if succ.x != p.x && succ.nextLower < p.y {
				if !hasVerticalBridge(columns[succ.x], p.y, succ.y, succ.y) &&
					!hasVerticalBridge(columns[p.x], p.y, succ.y, p.y) {
					return fmt.Errorf("points (%d,%d) and (%d,%d) need an intermediate point", succ.x, succ.y, p.x, p.y)
				}
			}
			if p.y > succ.nextLower {
				succ.nextLower = p.y
			}
		}
		t.insert(p)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	add := func(pts []point) {
		cp := append([]point(nil), pts...)
		tests = append(tests, testCase{
			input: formatInput(cp),
			base:  cp,
		})
	}
	add([]point{{0, 0}})
	add([]point{{0, 0}, {0, 5}})
	add([]point{{0, 0}, {5, 0}})
	add([]point{{0, 0}, {5, 5}})
	add([]point{{0, 0}, {5, 5}, {10, 10}})
	add([]point{{-3, -4}, {7, 8}, {7, -4}})
	add([]point{{-10, 0}, {0, 10}, {10, 0}, {0, -10}})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 90 {
		n := rng.Intn(25) + 2
		limit := 50 + rng.Intn(50)
		pts := randomPoints(rng, n, limit)
		add(pts)
	}

	// build some grid-based tests
	for size := 2; size <= 4; size++ {
		var pts []point
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				pts = append(pts, point{i * 3, j * 3})
			}
		}
		add(pts)
	}

	return tests
}

func randomPoints(rng *rand.Rand, n, limit int) []point {
	pts := make([]point, 0, n)
	seen := make(map[[2]int]struct{}, n)
	for len(pts) < n {
		x := rng.Intn(2*limit+1) - limit
		y := rng.Intn(2*limit+1) - limit
		key := [2]int{x, y}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		pts = append(pts, point{x, y})
	}
	return pts
}

func formatInput(points []point) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(points))
	for _, p := range points {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String()
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func hasVerticalBridge(ys []int, low, high, skip int) bool {
	if len(ys) == 0 || low > high {
		return false
	}
	i := sort.SearchInts(ys, low)
	for i < len(ys) && ys[i] <= high {
		if ys[i] != skip {
			return true
		}
		i++
	}
	return false
}
