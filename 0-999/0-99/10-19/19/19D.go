package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "math/rand"
)

// Treap node for y-values
type Node struct {
   key, priority, max int
   left, right *Node
}

func newNode(key int) *Node {
   return &Node{key: key, priority: rand.Int(), max: key}
}

func update(n *Node) {
   n.max = n.key
   if n.left != nil && n.left.max > n.max {
       n.max = n.left.max
   }
   if n.right != nil && n.right.max > n.max {
       n.max = n.right.max
   }
}

// split n into l: keys <= key, r: keys > key
func split(n *Node, key int) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if n.key <= key {
       // split right
       rl, rr := split(n.right, key)
       n.right = rl
       update(n)
       return n, rr
   }
   // n.key > key, split left
   ll, lr := split(n.left, key)
   n.left = lr
   update(n)
   return ll, n
}

// merge trees where all keys in a <= keys in b
func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.priority > b.priority {
       a.right = merge(a.right, b)
       update(a)
       return a
   }
   b.left = merge(a, b.left)
   update(b)
   return b
}

func insertTreap(root *Node, key int) *Node {
   l, r := split(root, key)
   n := newNode(key)
   return merge(merge(l, n), r)
}

func removeTreap(root *Node, key int) *Node {
   // split <=key and >key
   l, r := split(root, key)
   // split <=key-1 and ==key
   ll, _ := split(l, key-1)
   return merge(ll, r)
}

// lowerBound finds minimal node with key > key
func lowerBound(n *Node, key int) *Node {
   var res *Node
   for n != nil {
       if n.key > key {
           res = n
           n = n.left
       } else {
           n = n.right
       }
   }
   return res
}

// Segment tree for max y per x-index
var st []int
var size int

func segUpdate(node, l, r, pos, val int) {
   if l == r {
       st[node] = val
       return
   }
   mid := (l + r) >> 1
   if pos <= mid {
       segUpdate(node<<1, l, mid, pos, val)
   } else {
       segUpdate(node<<1|1, mid+1, r, pos, val)
   }
   if st[node<<1] > st[node<<1|1] {
       st[node] = st[node<<1]
   } else {
       st[node] = st[node<<1|1]
   }
}

// queryFirst finds leftmost index in [ql, qr] with max>y0, or -1
func queryFirst(node, l, r, ql, qr, y0 int) int {
   if r < ql || l > qr || st[node] <= y0 {
       return -1
   }
   if l == r {
       return l
   }
   mid := (l + r) >> 1
   // search left
   if mid >= ql {
       idx := queryFirst(node<<1, l, mid, ql, qr, y0)
       if idx != -1 {
           return idx
       }
   }
   return queryFirst(node<<1|1, mid+1, r, ql, qr, y0)
}

func main() {
   rand.Seed(1)
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   ops := make([]string, n)
   xs := make([]int, 0, n)
   xarr := make([]int, n)
   yarr := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &ops[i], &xarr[i], &yarr[i])
       if ops[i] == "add" {
           xs = append(xs, xarr[i])
       }
   }
   sort.Ints(xs)
   xs = unique(xs)
   m := len(xs)
   // map x to index
   // use binary search on xs

   // init segment tree and treaps
   size = m
   st = make([]int, 4*m+4)
   roots := make([]*Node, m)

   for i := 0; i < n; i++ {
       op := ops[i]
       x := xarr[i]
       y := yarr[i]
       switch op {
       case "add":
           idx := sort.SearchInts(xs, x)
           roots[idx] = insertTreap(roots[idx], y)
           segUpdate(1, 0, m-1, idx, roots[idx].max)
       case "remove":
           idx := sort.SearchInts(xs, x)
           roots[idx] = removeTreap(roots[idx], y)
           maxv := 0
           if roots[idx] != nil {
               maxv = roots[idx].max
           }
           segUpdate(1, 0, m-1, idx, maxv)
       case "find":
           // find first x > x
           pos := sort.SearchInts(xs, x+1)
           if pos >= m {
               fmt.Fprintln(out, -1)
               continue
           }
           idx := queryFirst(1, 0, m-1, pos, m-1, y)
           if idx == -1 {
               fmt.Fprintln(out, -1)
           } else {
               node := lowerBound(roots[idx], y)
               // node must exist
               fmt.Fprintf(out, "%d %d\n", xs[idx], node.key)
           }
       }
   }
}

// unique returns sorted unique slice
func unique(a []int) []int {
   j := 0
   for i := 0; i < len(a); i++ {
       if j == 0 || a[i] != a[j-1] {
           a[j] = a[i]
           j++
       }
   }
   return a[:j]
}
