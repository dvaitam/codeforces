package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Implicit treap node
type Node struct {
   imp, id, prio, sz, maxImp int
   left, right *Node
}

func sz(n *Node) int {
   if n == nil {
       return 0
   }
   return n.sz
}

func upd(n *Node) {
   if n == nil {
       return
   }
   n.sz = 1 + sz(n.left) + sz(n.right)
   n.maxImp = n.imp
   if n.left != nil && n.left.maxImp > n.maxImp {
       n.maxImp = n.left.maxImp
   }
   if n.right != nil && n.right.maxImp > n.maxImp {
       n.maxImp = n.right.maxImp
   }
}

// merge two treaps l and r, all keys in l come before r
func merge(l, r *Node) *Node {
   if l == nil {
       return r
   }
   if r == nil {
       return l
   }
   if l.prio > r.prio {
       l.right = merge(l.right, r)
       upd(l)
       return l
   }
   r.left = merge(l, r.left)
   upd(r)
   return r
}

// split treap t into [0..k-1] and [k..]
func split(t *Node, k int) (l, r *Node) {
   if t == nil {
       return nil, nil
   }
   if sz(t.left) >= k {
       // split left
       l0, r0 := split(t.left, k)
       t.left = r0
       upd(t)
       return l0, t
   }
   // left size < k
   l0, r0 := split(t.right, k - sz(t.left) - 1)
   t.right = l0
   upd(t)
   return t, r0
}

// find last position where imp >= value, return 1-based index in this treap or 0 if none
func findLastGE(n *Node, value int) int {
   if n == nil || n.maxImp < value {
       return 0
   }
   // try right
   if n.right != nil && n.right.maxImp >= value {
       idx := findLastGE(n.right, value)
       if idx > 0 {
           return sz(n.left) + 1 + idx
       }
   }
   if n.imp >= value {
       return sz(n.left) + 1
   }
   return findLastGE(n.left, value)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   c := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &c[i])
   }
   rand.Seed(time.Now().UnixNano())
   var root *Node
   for i := 0; i < n; i++ {
       imp := a[i]
       cnt := c[i]
       // current size
       cur := sz(root)
       // lo = max(1, cur+1-cnt)
       lo := 1
       if cur+1-cnt > 1 {
           lo = cur + 1 - cnt
       }
       // split [1..lo-1], [lo..]
       left, right := split(root, lo-1)
       // find last >= imp in right
       posInRight := findLastGE(right, imp)
       // compute global pos
       var pos int
       if posInRight > 0 {
           pos = (lo - 1) + posInRight
       } else {
           pos = 0
       }
       // target insert pos p (1-based) = max(pos+1, lo)
       p := lo
       if pos+1 > lo {
           p = pos + 1
       }
       // merge back left and right before insertion
       root = merge(left, right)
       // now split at p-1
       l2, r2 := split(root, p-1)
       // create node
       node := &Node{imp: imp, id: i + 1, prio: rand.Int(), sz: 1, maxImp: imp}
       // merge l2 + node + r2
       root = merge(merge(l2, node), r2)
   }
   // output ids in order
   var dfs func(n *Node)
   dfs = func(n *Node) {
       if n == nil {
           return
       }
       dfs(n.left)
       fmt.Fprint(out, n.id, ' ')
       dfs(n.right)
   }
   dfs(root)
}
