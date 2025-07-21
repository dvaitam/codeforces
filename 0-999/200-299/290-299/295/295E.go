package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap node storing key, subtree size, sum of keys, and sum of pairwise distances
type Node struct {
   key   int64
   pri   int
   left  *Node
   right *Node
   sz    int
   sumX  int64
   sumP  int64
}

// update recalculates sz, sumX, sumP for node n
func update(n *Node) {
   if n == nil {
       return
   }
   n.sz = 1
   n.sumX = n.key
   n.sumP = 0
   l, r := n.left, n.right
   if l != nil {
       n.sz += l.sz
       n.sumX += l.sumX
       n.sumP += l.sumP
   }
   if r != nil {
       n.sz += r.sz
       n.sumX += r.sumX
       n.sumP += r.sumP
   }
   // pairs between L and root
   if l != nil {
       n.sumP += int64(l.sz)*n.key - l.sumX
   }
   // pairs between root and R
   if r != nil {
       n.sumP += r.sumX - int64(r.sz)*n.key
   }
   // pairs between L and R
   if l != nil && r != nil {
       n.sumP += int64(l.sz)*r.sumX - int64(r.sz)*l.sumX
   }
}

// merge merges treaps a and b, all keys in a <= keys in b
func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.pri > b.pri {
       a.right = merge(a.right, b)
       update(a)
       return a
   }
   b.left = merge(a, b.left)
   update(b)
   return b
}

// split splits treap n into l and r, where keys in l <= key, keys in r > key
func split(n *Node, key int64) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if n.key <= key {
       var rr *Node
       n.right, rr = split(n.right, key)
       update(n)
       return n, rr
   }
   var ll *Node
   ll, n.left = split(n.left, key)
   update(n)
   return ll, n
}

// insert key into treap
func insert(root *Node, key int64) *Node {
   newNode := &Node{key: key, pri: rand.Int(), sz: 1, sumX: key, sumP: 0}
   var a, b *Node
   a, b = split(root, key)
   return merge(merge(a, newNode), b)
}

// erase key from treap (assumes key exists)
func erase(root *Node, key int64) *Node {
   var a, b, c *Node
   a, b = split(root, key-1)
   c, b = split(b, key)
   // drop c
   return merge(a, b)
}

// rangeSumP returns sum of pairwise distances for keys in [l, r]
func rangeSumP(root *Node, lkey, rkey int64) (ans int64, rootOut *Node) {
   var a, b, c, d *Node
   a, b = split(root, lkey-1)
   c, d = split(b, rkey)
   if c != nil {
       ans = c.sumP
   }
   // merge back
   rootOut = merge(a, merge(c, d))
   return
}

func main() {
   rand.Seed(time.Now().UnixNano())
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   pos := make([]int64, n+1)
   var root *Node
   for i := 1; i <= n; i++ {
       var x int64
       fmt.Fscan(in, &x)
       pos[i] = x
       root = insert(root, x)
   }
   var m int
   fmt.Fscan(in, &m)
   for j := 0; j < m; j++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           var p int
           var d int64
           fmt.Fscan(in, &p, &d)
           old := pos[p]
           new := old + d
           root = erase(root, old)
           root = insert(root, new)
           pos[p] = new
       } else if t == 2 {
           var l, r int64
           fmt.Fscan(in, &l, &r)
           ans, newRoot := rangeSumP(root, l, r)
           root = newRoot
           fmt.Fprintln(out, ans)
       }
   }
}
