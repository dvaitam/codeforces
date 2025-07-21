package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
)

const MAXLOG = 18

var (
   n int
   adj [][]edge
   up [][]int
   depth []int
   dist []int64
   tin []int
   timer int
)

type edge struct{ to, w int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &n)
   adj = make([][]edge, n+1)
   for i := 0; i < n-1; i++ {
       var a, b, c int
       fmt.Fscan(in, &a, &b, &c)
       adj[a] = append(adj[a], edge{b, c})
       adj[b] = append(adj[b], edge{a, c})
   }
   up = make([][]int, n+1)
   for i := range up {
       up[i] = make([]int, MAXLOG)
   }
   depth = make([]int, n+1)
   dist = make([]int64, n+1)
   tin = make([]int, n+1)
   // iterative DFS
   timer = 0
   type item struct{ v, p, it int }
   stack := []item{{1, 1, 0}}
   up[1][0] = 1
   depth[1] = 0
   dist[1] = 0
   for len(stack) > 0 {
       top := &stack[len(stack)-1]
       v, p, it := top.v, top.p, top.it
       if it == 0 {
           tin[v] = timer; timer++
       }
       if it < len(adj[v]) {
           e := adj[v][it]
           top.it++
           if e.to == p {
               continue
           }
           u := e.to
           depth[u] = depth[v] + 1
           dist[u] = dist[v] + int64(e.w)
           up[u][0] = v
           stack = append(stack, item{u, v, 0})
       } else {
           stack = stack[:len(stack)-1]
       }
   }
   // build LCA
   for j := 1; j < MAXLOG; j++ {
       for v := 1; v <= n; v++ {
           up[v][j] = up[ up[v][j-1] ][j-1]
       }
   }
   // active set treap
   var root *Node
   activeCount := 0
   var totalPerim int64 = 0
   // process queries
   var q int
   fmt.Fscan(in, &q)
   rand.Seed(1)
   for i := 0; i < q; i++ {
       var ops string
       fmt.Fscan(in, &ops)
       op := ops[0]
       if op == '?' {
           if activeCount < 2 {
               fmt.Fprintln(out, 0)
           } else {
               fmt.Fprintln(out, totalPerim/2)
           }
       } else {
           var x int
           fmt.Fscan(in, &x)
           if op == '+' {
               // insert x
               if activeCount == 0 {
                   root = insertNode(root, tin[x], x)
                   activeCount = 1
               } else {
                   // find pred and succ
                   pnode := predecessor(root, tin[x])
                   if pnode == nil {
                       pnode = maxNode(root)
                   }
                   snode := successor(root, tin[x])
                   if snode == nil {
                       snode = minNode(root)
                   }
                   dxp := distNodes(pnode.val, x)
                   dxs := distNodes(x, snode.val)
                   dps := distNodes(pnode.val, snode.val)
                   totalPerim += dxp + dxs - dps
                   root = insertNode(root, tin[x], x)
                   activeCount++
               }
           } else if op == '-' {
               // remove x
               if activeCount <= 1 {
                   root = deleteNode(root, tin[x])
                   activeCount = 0
                   totalPerim = 0
               } else {
                   // find pred and succ
                   pnode := predecessor(root, tin[x])
                   if pnode == nil {
                       pnode = maxNode(root)
                   }
                   snode := successor(root, tin[x])
                   if snode == nil {
                       snode = minNode(root)
                   }
                   dxp := distNodes(pnode.val, x)
                   dxs := distNodes(x, snode.val)
                   dps := distNodes(pnode.val, snode.val)
                   totalPerim -= dxp + dxs - dps
                   root = deleteNode(root, tin[x])
                   activeCount--
               }
           }
       }
   }
}

func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   diff := depth[u] - depth[v]
   for j := 0; j < MAXLOG; j++ {
       if diff&(1<<j) != 0 {
           u = up[u][j]
       }
   }
   if u == v {
       return u
   }
   for j := MAXLOG-1; j >= 0; j-- {
       if up[u][j] != up[v][j] {
           u = up[u][j]
           v = up[v][j]
       }
   }
   return up[u][0]
}

func distNodes(u, v int) int64 {
   w := lca(u, v)
   return dist[u] + dist[v] - 2*dist[w]
}

// Treap for (key, val)
type Node struct {
   key, val, pri int
   left, right *Node
}

func insertNode(root *Node, key, val int) *Node {
   nnode := &Node{key: key, val: val, pri: rand.Int()}
   var l, r *Node
   l, r = split(root, key)
   return merge(merge(l, nnode), r)
}

func deleteNode(root *Node, key int) *Node {
   // split out nodes < key and >= key
   l, rest := split(root, key)
   // split rest into nodes == key and > key
   _, r := split(rest, key+1)
   return merge(l, r)
}

func split(root *Node, key int) (l, r *Node) {
   if root == nil {
       return nil, nil
   }
   if key <= root.key {
       l2, r2 := split(root.left, key)
       root.left = r2
       return l2, root
   } else {
       l2, r2 := split(root.right, key)
       root.right = l2
       return root, r2
   }
}

func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.pri < b.pri {
       a.right = merge(a.right, b)
       return a
   } else {
       b.left = merge(a, b.left)
       return b
   }
}

func predecessor(root *Node, key int) *Node {
   var res *Node
   for node := root; node != nil; {
       if node.key < key {
           res = node
           node = node.right
       } else {
           node = node.left
       }
   }
   return res
}

func successor(root *Node, key int) *Node {
   var res *Node
   for node := root; node != nil; {
       if node.key > key {
           res = node
           node = node.left
       } else {
           node = node.right
       }
   }
   return res
}

func minNode(root *Node) *Node {
   node := root
   if node == nil {
       return nil
   }
   for node.left != nil {
       node = node.left
   }
   return node
}

func maxNode(root *Node) *Node {
   node := root
   if node == nil {
       return nil
   }
   for node.right != nil {
       node = node.right
   }
   return node
}
