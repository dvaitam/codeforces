package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 1000000007

type Rect struct {
   xl, xr, yl, yr, z, id int
}
type Shot struct {
   x, y, yprime, xpos int
}

// Segment tree over x, each node has a y-segment-tree for min t
type Node struct {
   l, r     int
   ys       []int64
   tree     []int
   left, right *Node
}

func buildNode(l, r int, leafY [][]int64) *Node {
   node := &Node{l: l, r: r}
   if l == r {
       ys := leafY[l]
       sort.Slice(ys, func(i, j int) bool { return ys[i] < ys[j] })
       node.ys = ys
   } else {
       m := (l + r) >> 1
       node.left = buildNode(l, m, leafY)
       node.right = buildNode(m+1, r, leafY)
       a, b := node.left.ys, node.right.ys
       node.ys = make([]int64, len(a)+len(b))
       i, j := 0, 0
       for k := 0; k < len(node.ys); k++ {
           if j>=len(b) || (i<len(a) && a[i] < b[j]) {
               node.ys[k] = a[i]; i++
           } else {
               node.ys[k] = b[j]; j++
           }
       }
   }
   // build segtree array
   node.tree = make([]int, 4*len(node.ys))
   for i := range node.tree {
       node.tree[i] = INF
   }
   return node
}

// update segtree at pos to val
func segUpdate(tree []int, idx, l, r, pos, val int) {
   if l == r {
       tree[idx] = val
       return
   }
   m := (l + r) >> 1
   if pos <= m {
       segUpdate(tree, idx*2, l, m, pos, val)
   } else {
       segUpdate(tree, idx*2+1, m+1, r, pos, val)
   }
   // pull up
   a, b := tree[idx*2], tree[idx*2+1]
   if a < b {
       tree[idx] = a
   } else {
       tree[idx] = b
   }
}
// query segtree min in [ql, qr]
func segQuery(tree []int, idx, l, r, ql, qr int) int {
   if qr < l || r < ql {
       return INF
   }
   if ql <= l && r <= qr {
       return tree[idx]
   }
   m := (l + r) >> 1
   a := segQuery(tree, idx*2, l, m, ql, qr)
   b := segQuery(tree, idx*2+1, m+1, r, ql, qr)
   if a < b { return a } else { return b }
}

// point update at x=xpos, yprime
func updateNode(node *Node, xpos int, yprime int64, val int) {
   // update this node
   // find y index
   ys := node.ys
   i := sort.Search(len(ys), func(i int) bool { return ys[i] >= yprime })
   if i < len(ys) && ys[i] == yprime {
       segUpdate(node.tree, 1, 0, len(ys)-1, i, val)
   }
   if node.l == node.r {
       return
   }
   if xpos <= node.left.r {
       updateNode(node.left, xpos, yprime, val)
   } else {
       updateNode(node.right, xpos, yprime, val)
   }
}

// query rectangle x in [xl, xr] (xpos), y' in [yl, yr]
func queryNode(node *Node, xl, xr int, yL, yR int64) int {
   if node == nil || xr < node.l || node.r < xl {
       return INF
   }
   if xl <= node.l && node.r <= xr {
       // full cover
       ys := node.ys
       L := sort.Search(len(ys), func(i int) bool { return ys[i] >= yL })
       R := sort.Search(len(ys), func(i int) bool { return ys[i] > yR })
       if L >= len(ys) || L >= R {
           return INF
       }
       return segQuery(node.tree, 1, 0, len(ys)-1, L, R-1)
   }
   a := queryNode(node.left, xl, xr, yL, yR)
   b := queryNode(node.right, xl, xr, yL, yR)
   if a < b { return a } else { return b }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   rects := make([]Rect, n)
   for i := 0; i < n; i++ {
       r := &rects[i]
       fmt.Fscan(in, &r.xl, &r.xr, &r.yl, &r.yr, &r.z)
       r.id = i + 1
   }
   var m int
   fmt.Fscan(in, &m)
   shots := make([]Shot, m+1)
   xs := make([]int, 0, m)
   for i := 1; i <= m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       shots[i].x = x; shots[i].y = y
       xs = append(xs, x)
   }
   // compress xs
   sort.Ints(xs)
   xs = uniqueInts(xs)
   Nx := len(xs)
   leafY := make([][]int64, Nx)
   for i := 1; i <= m; i++ {
       // compute xpos
       xi := sort.Search(Nx, func(j int) bool { return xs[j] >= shots[i].x })
       shots[i].xpos = xi
       // y'
       yprime := int64(shots[i].y)*(int64(m)+1) + int64(i)
       shots[i].yprime = yprime
       leafY[xi] = append(leafY[xi], yprime)
   }
   // build tree
   root := buildNode(0, Nx-1, leafY)
   // insert shots
   for i := 1; i <= m; i++ {
       updateNode(root, shots[i].xpos, shots[i].yprime, i)
   }
   // sort rects by z
   sort.Slice(rects, func(i, j int) bool { return rects[i].z < rects[j].z })
   ans := make([]int, m+1)
   // process
   for _, r := range rects {
       // x range
       L := sort.Search(Nx, func(i int) bool { return xs[i] >= r.xl })
       R := sort.Search(Nx, func(i int) bool { return xs[i] > r.xr }) - 1
       if L > R { continue }
       yL := int64(r.yl)*(int64(m)+1) + 1
       yR := int64(r.yr+1)*(int64(m)+1) - 1
       t := queryNode(root, L, R, yL, yR)
       if t < INF {
           ans[t] = r.id
           updateNode(root, shots[t].xpos, shots[t].yprime, INF)
       }
   }
   // output per shot
   for i := 1; i <= m; i++ {
       fmt.Fprintln(out, ans[i])
   }
}

func uniqueInts(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if i == 0 || a[i] != a[i-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
