package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type node struct {
	y, x        int64
	pr          int
	left, right *node
}

func split(root *node, y int64) (l, r *node) {
	if root == nil {
		return nil, nil
	}
	if y < root.y {
		l, root.left = split(root.left, y)
		return l, root
	}
	root.right, r = split(root.right, y)
	return root, r
}

func merge(a, b *node) *node {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.pr < b.pr {
		a.right = merge(a.right, b)
		return a
	}
	b.left = merge(a, b.left)
	return b
}

func insert(root, n *node) *node {
	if root == nil {
		return n
	}
	if n.pr < root.pr {
		l, r := split(root, n.y)
		n.left = l
		n.right = r
		return n
	}
	if n.y < root.y {
		root.left = insert(root.left, n)
	} else {
		root.right = insert(root.right, n)
	}
	return root
}

func erase(root *node, y int64) *node {
	if root == nil {
		return nil
	}
	if y < root.y {
		root.left = erase(root.left, y)
		return root
	}
	if y > root.y {
		root.right = erase(root.right, y)
		return root
	}
	return merge(root.left, root.right)
}

func find(root *node, y int64) *node {
	for root != nil {
		if y < root.y {
			root = root.left
		} else if y > root.y {
			root = root.right
		} else {
			return root
		}
	}
	return nil
}

func predecessor(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y <= root.y {
			root = root.left
		} else {
			res = root
			root = root.right
		}
	}
	return res
}

func successor(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y >= root.y {
			root = root.right
		} else {
			res = root
			root = root.left
		}
	}
	return res
}

func lowerBound(root *node, y int64) *node {
	var res *node
	for root != nil {
		if y <= root.y {
			res = root
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func maxNode(root *node) *node {
	if root == nil {
		return nil
	}
	for root.right != nil {
		root = root.right
	}
	return root
}

func cross(a, b, c *node) int64 {
	return (b.y-a.y)*(c.x-b.x) - (b.x-a.x)*(c.y-b.y)
}

var rootF *node

func isBad(n *node) bool {
	if n == nil || n.y == 0 {
		return false
	}
	p := predecessor(rootF, n.y)
	r := successor(rootF, n.y)
	if p == nil || r == nil {
		return false
	}
	if p.y == 0 && r == nil {
		return false
	}
	return cross(p, n, r) >= 0
}

func addSpell(x, y int64) {
	if existing := find(rootF, y); existing != nil {
		if existing.x >= x {
			return
		}
		rootF = erase(rootF, y)
	}
	n := &node{y: y, x: x, pr: rand.Int()}
	rootF = insert(rootF, n)
	if isBad(n) {
		rootF = erase(rootF, n.y)
		return
	}
	for {
		p := predecessor(rootF, n.y)
		if p == nil || p.y == 0 {
			break
		}
		if isBad(p) {
			rootF = erase(rootF, p.y)
		} else {
			break
		}
	}
	for {
		r := successor(rootF, n.y)
		if r == nil {
			break
		}
		if isBad(r) {
			rootF = erase(rootF, r.y)
		} else {
			break
		}
	}
}

func hullValue(z float64) float64 {
	if rootF == nil {
		return 0
	}
	if z <= 0 {
		return 0
	}
	mx := maxNode(rootF)
	if z >= float64(mx.y) {
		return float64(mx.x)
	}
	r := lowerBound(rootF, int64(math.Ceil(z)))
	if r == nil {
		r = mx
	}
	if float64(r.y) == z {
		return float64(r.x)
	}
	l := predecessor(rootF, r.y)
	if l == nil {
		return (float64(r.x) / float64(r.y)) * z
	}
	return float64(l.x) + (float64(r.x-l.x))*(z-float64(l.y))/float64(r.y-l.y)
}

type queryF struct {
	k    int
	a, b int64
}

func solveCaseF(q int, m int64, qs []queryF) []string {
	rootF = &node{y: 0, x: 0, pr: rand.Int()}
	lastOK := 0
	res := make([]string, 0, q)
	for i := 1; i <= q; i++ {
		k := qs[i-1].k
		a := (qs[i-1].a+int64(lastOK))%1000000 + 1
		b := (qs[i-1].b+int64(lastOK))%1000000 + 1
		if k == 1 {
			addSpell(a, b)
		} else {
			t := a
			h := b
			rate := hullValue(float64(m) / float64(t))
			if rate*float64(t) >= float64(h) {
				lastOK = i
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		}
	}
	return res
}

func genCaseF(rng *rand.Rand) (string, []string) {
	q := rng.Intn(10) + 1
	m := int64(rng.Intn(1000) + 1)
	qs := make([]queryF, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", q, m))
	for i := 0; i < q; i++ {
		k := rng.Intn(2) + 1
		a := int64(rng.Intn(100) + 1)
		b := int64(rng.Intn(100) + 1)
		qs[i] = queryF{k, a, b}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", k, a, b))
	}
	exp := solveCaseF(q, m, qs)
	return sb.String(), exp
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCaseF(bin string, in string, exp []string) error {
	out, err := runCandidate(bin, []byte(in))
	if err != nil {
		return err
	}
	tokens := strings.Fields(out)
	if len(tokens) != len(exp) {
		return fmt.Errorf("expected %d tokens, got %d", len(exp), len(tokens))
	}
	for i := range tokens {
		if tokens[i] != exp[i] {
			return fmt.Errorf("at query %d expected %s got %s", i+1, exp[i], tokens[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCaseF(rng)
		if err := runCaseF(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
