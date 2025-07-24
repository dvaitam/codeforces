package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const MAXK = 10

// Point in k-dimensional space with weight.
type Point struct {
	coord [MAXK]float64
	w     int64
}

type Node struct {
	point       *Point
	left, right *Node
	min, max    [MAXK]float64
	sum         int64
}

var k int

func build(points []*Point, depth int) *Node {
	if len(points) == 0 {
		return nil
	}
	axis := depth % k
	sort.Slice(points, func(i, j int) bool {
		return points[i].coord[axis] < points[j].coord[axis]
	})
	mid := len(points) / 2
	node := &Node{point: points[mid]}
	node.left = build(points[:mid], depth+1)
	node.right = build(points[mid+1:], depth+1)
	for i := 0; i < k; i++ {
		v := node.point.coord[i]
		mn, mx := v, v
		if node.left != nil {
			if node.left.min[i] < mn {
				mn = node.left.min[i]
			}
			if node.left.max[i] > mx {
				mx = node.left.max[i]
			}
		}
		if node.right != nil {
			if node.right.min[i] < mn {
				mn = node.right.min[i]
			}
			if node.right.max[i] > mx {
				mx = node.right.max[i]
			}
		}
		node.min[i] = mn
		node.max[i] = mx
	}
	node.sum = node.point.w
	if node.left != nil {
		node.sum += node.left.sum
	}
	if node.right != nil {
		node.sum += node.right.sum
	}
	return node
}

func minDist(node *Node, center [MAXK]float64) float64 {
	dist := 0.0
	for i := 0; i < k; i++ {
		if center[i] < node.min[i] {
			d := node.min[i] - center[i]
			if d > dist {
				dist = d
			}
		} else if center[i] > node.max[i] {
			d := center[i] - node.max[i]
			if d > dist {
				dist = d
			}
		}
	}
	return dist
}

func maxDist(node *Node, center [MAXK]float64) float64 {
	dist := 0.0
	for i := 0; i < k; i++ {
		d := math.Max(math.Abs(center[i]-node.min[i]), math.Abs(center[i]-node.max[i]))
		if d > dist {
			dist = d
		}
	}
	return dist
}

func query(node *Node, center [MAXK]float64, t float64) int64 {
	if node == nil {
		return 0
	}
	if minDist(node, center) > t {
		return 0
	}
	if maxDist(node, center) <= t {
		return node.sum
	}
	// check the point stored at this node
	dist := 0.0
	for i := 0; i < k; i++ {
		d := math.Abs(node.point.coord[i] - center[i])
		if d > dist {
			dist = d
		}
	}
	res := int64(0)
	if dist <= t {
		res += node.point.w
	}
	res += query(node.left, center, t)
	res += query(node.right, center, t)
	return res
}

func cross(ax, ay, bx, by int64) int64 {
	return ax*by - ay*bx
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &k, &n, &q)

	vx := make([]int64, k)
	vy := make([]int64, k)
	for i := 0; i < k; i++ {
		fmt.Fscan(in, &vx[i], &vy[i])
	}

	S := make([]float64, k)
	for j := 0; j < k; j++ {
		sum := int64(0)
		for i := 0; i < k; i++ {
			if i == j {
				continue
			}
			sum += abs64(cross(vx[j], vy[j], vx[i], vy[i]))
		}
		S[j] = float64(sum)
	}

	factories := make([]*Point, n)
	for i := 0; i < n; i++ {
		var fx, fy, a int64
		fmt.Fscan(in, &fx, &fy, &a)
		p := &Point{w: a}
		for j := 0; j < k; j++ {
			val := float64(cross(vx[j], vy[j], fx, fy)) / S[j]
			p.coord[j] = val
		}
		factories[i] = p
	}

	root := build(factories, 0)

	for ; q > 0; q-- {
		var px, py, t int64
		fmt.Fscan(in, &px, &py, &t)
		var center [MAXK]float64
		for j := 0; j < k; j++ {
			center[j] = float64(cross(vx[j], vy[j], px, py)) / S[j]
		}
		ans := query(root, center, float64(t))
		fmt.Fprintln(out, ans)
	}
}

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
