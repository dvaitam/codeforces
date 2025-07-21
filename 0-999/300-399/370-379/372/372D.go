package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap node for storing tin keys
type Node struct {
   key   int
   prio  int
   left, right *Node
   id    int // original node id
}

// split treap t into l (< key) and r (>= key)
func split(t *Node, key int) (l, r *Node) {
   if t == nil {
       return nil, nil
   }
   if key <= t.key {
       lsub, rsub := split(t.left, key)
       t.left = rsub
       return lsub, t
   }
   lsub, rsub := split(t.right, key)
   t.right = lsub
   return t, rsub
}

// merge treaps a and b, where all keys in a < keys in b
func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio < b.prio {
       a.right = merge(a.right, b)
       return a
   }
   b.left = merge(a, b.left)
   return b
}

// insert node x into treap t
func insert(t *Node, x *Node) *Node {
   if t == nil {
       return x
   }
   if x.prio < t.prio {
       l, r := split(t, x.key)
       x.left, x.right = l, r
       return x
   }
   if x.key < t.key {
       t.left = insert(t.left, x)
   } else {
       t.right = insert(t.right, x)
   }
   return t
}

// delete key from treap t
func erase(t *Node, key int) *Node {
   if t == nil {
       return nil
   }
   if t.key == key {
       return merge(t.left, t.right)
   }
   if key < t.key {
       t.left = erase(t.left, key)
   } else {
       t.right = erase(t.right, key)
   }
   return t
}

// find predecessor (max key < key) in treap
func findPred(t *Node, key int) *Node {
   var res *Node
   for t != nil {
       if t.key < key {
           res = t
           t = t.right
       } else {
           t = t.left
       }
   }
   return res
}

// find successor (min key > key) in treap
func findSucc(t *Node, key int) *Node {
   var res *Node
   for t != nil {
       if t.key > key {
           res = t
           t = t.left
       } else {
           t = t.right
       }
   }
   return res
}

// find min node in treap
func findMin(t *Node) *Node {
   if t == nil {
       return nil
   }
   for t.left != nil {
       t = t.left
   }
   return t
}

// find max node in treap
func findMax(t *Node) *Node {
   if t == nil {
       return nil
   }
   for t.right != nil {
       t = t.right
   }
   return t
}

var (
   n, k int
   adj [][]int
   up  [][]int
   depth []int
   tin   []int
   timer int
)

// LCA preprocessing
func dfs(u, p int) {
   tin[u] = timer; timer++
   up[0][u] = p
   if p >= 0 {
       depth[u] = depth[p] + 1
   }
   for _, v := range adj[u] {
       if v == p {
           continue
       }
       dfs(v, u)
   }
}

func lca(u, v int) int {
   if depth[u] < depth[v] {
       u, v = v, u
   }
   d := depth[u] - depth[v]
   for i := 0; d > 0; i++ {
       if d&1 != 0 {
           u = up[i][u]
       }
       d >>= 1
   }
   if u == v {
       return u
   }
   for i := len(up)-1; i >= 0; i-- {
       if up[i][u] != up[i][v] {
           u = up[i][u]
           v = up[i][v]
       }
   }
   return up[0][u]
}

func dist(u, v int) int {
   w := lca(u, v)
   return depth[u] + depth[v] - 2*depth[w]
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &k)
   adj = make([][]int, n+1)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(reader, &a, &b)
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   // init LCA
   logn := 1
   for (1<<logn) <= n {
       logn++
   }
   up = make([][]int, logn)
   for i := range up {
       up[i] = make([]int, n+1)
   }
   depth = make([]int, n+1)
   tin = make([]int, n+1)
   timer = 0
   dfs(1, -1)
   for i := 1; i < len(up); i++ {
       for v := 1; v <= n; v++ {
           pu := up[i-1][v]
           if pu >= 0 {
               up[i][v] = up[i-1][pu]
           } else {
               up[i][v] = -1
           }
       }
   }
   rand.Seed(time.Now().UnixNano())
   var root *Node
   var cycleLen int64
   var sz int
   l := 1
   ans := 1
   // sliding window on labels
   for r := 1; r <= n; r++ {
       // insert r
       x := &Node{key: tin[r], prio: rand.Int(), id: r}
       if sz > 0 {
           pred := findPred(root, x.key)
           if pred == nil {
               pred = findMax(root)
           }
           succ := findSucc(root, x.key)
           if succ == nil {
               succ = findMin(root)
           }
           delta := int64(dist(pred.id, x.id) + dist(x.id, succ.id) - dist(pred.id, succ.id))
           cycleLen += delta
       }
       root = insert(root, x)
       sz++
       // shrink from left while tree size > k
       for {
           // current tree size = cycleLen/2 + 1
           if sz > 0 && cycleLen/2 + 1 > int64(k) {
               // remove l
               key := tin[l]
               // find node x's neighbors
               pred := findPred(root, key)
               if pred == nil {
                   pred = findMax(root)
               }
               succ := findSucc(root, key)
               if succ == nil {
                   succ = findMin(root)
               }
               delta := int64(dist(pred.id, l) + dist(l, succ.id) - dist(pred.id, succ.id))
               cycleLen -= delta
               root = erase(root, key)
               sz--
               l++
           } else {
               break
           }
       }
       // update answer
       if r-l+1 > ans {
           ans = r - l + 1
       }
   }
   fmt.Println(ans)
}
