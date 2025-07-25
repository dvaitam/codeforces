package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Node represents a treap node storing fruits by deadline
type Node struct {
   key   int       // deadline
   weight int64    // sum of weights for this deadline
   sum   int64     // sum of subtree weights
   cnt   int       // count of nodes in subtree
   pri   int
   left, right *Node
}

// update recalculates sum and cnt
func (n *Node) update() {
   n.sum = n.weight
   n.cnt = 1
   if n.left != nil {
       n.sum += n.left.sum
       n.cnt += n.left.cnt
   }
   if n.right != nil {
       n.sum += n.right.sum
       n.cnt += n.right.cnt
   }
}

// rotate right
func rotateRight(n *Node) *Node {
   l := n.left
   n.left = l.right
   l.right = n
   n.update()
   l.update()
   return l
}

// rotate left
func rotateLeft(n *Node) *Node {
   r := n.right
   n.right = r.left
   r.left = n
   n.update()
   r.update()
   return r
}

// insert or add weight for key
func insert(n *Node, key int, w int64) *Node {
   if n == nil {
       return &Node{key: key, weight: w, sum: w, cnt: 1, pri: rand.Int()}
   }
   if key == n.key {
       n.weight += w
   } else if key < n.key {
       n.left = insert(n.left, key, w)
       if n.left.pri > n.pri {
           n = rotateRight(n)
       }
   } else {
       n.right = insert(n.right, key, w)
       if n.right.pri > n.pri {
           n = rotateLeft(n)
       }
   }
   n.update()
   return n
}

// split treap into <= key and > key
func split(n *Node, key int) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if n.key <= key {
       var rr *Node
       rr, r = split(n.right, key)
       n.right = rr
       n.update()
       l = n
   } else {
       var ll *Node
       l, ll = split(n.left, key)
       n.left = ll
       n.update()
       r = n
   }
   return
}

// merge assumes all keys in l <= keys in r
func merge(l, r *Node) *Node {
   if l == nil {
       return r
   }
   if r == nil {
       return l
   }
   if l.pri > r.pri {
       l.right = merge(l.right, r)
       l.update()
       return l
   }
   r.left = merge(l, r.left)
   r.update()
   return r
}

// Bag holds treap root
type Bag struct {
   root *Node
}

// merge another bag into this
func (b *Bag) absorb(other *Bag) {
   // merge smaller into larger
   if other.root == nil {
       return
   }
   if b.root == nil {
       b.root = other.root
       return
   }
   if b.root.cnt < other.root.cnt {
       // swap
       b.root, other.root = other.root, b.root
   }
   // insert all nodes from other into b
   other.traverse(func(key int, w int64) {
       b.root = insert(b.root, key, w)
   })
}

// traverse calls fn for each node
func (b *Bag) traverse(fn func(int, int64)) {
   var dfs func(*Node)
   dfs = func(n *Node) {
       if n == nil {
           return
       }
       dfs(n.left)
       fn(n.key, n.weight)
       dfs(n.right)
   }
   dfs(b.root)
}

func main() {
   rand.Seed(time.Now().UnixNano())
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   children := make([][]int, n+1)
   parents := make([]int, n+1)
   for i := 2; i <= n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       parents[i] = p
       children[p] = append(children[p], i)
   }
   // fruits data
   fruitD := make([]int, n+1)
   fruitW := make([]int64, n+1)
   for i := 0; i < m; i++ {
       var v, d int
       var w int64
       fmt.Fscan(reader, &v, &d, &w)
       fruitD[v] = d
       fruitW[v] = w
   }
   // post-order using BFS reversed
   order := make([]int, 0, n)
   queue := []int{1}
   for i := 0; i < len(queue); i++ {
       v := queue[i]
       order = append(order, v)
       for _, c := range children[v] {
           queue = append(queue, c)
       }
   }
   // DP bags
   bags := make([]*Bag, n+1)
   for i := n - 1; i >= 0; i-- {
       v := order[i]
       b := &Bag{}
       // absorb children
       for _, c := range children[v] {
           if bags[c] != nil {
               b.absorb(bags[c])
               bags[c] = nil
           }
       }
       // process fruit at v
       if fruitW[v] > 0 {
           total0 := int64(0)
           if b.root != nil {
               total0 = b.root.sum
           }
           // split by deadline
           l, r := split(b.root, fruitD[v])
           sumLE := int64(0)
           if l != nil {
               sumLE = l.sum
           }
           cand := sumLE + fruitW[v]
           if cand > total0 {
               // keep l as new root
               b.root = insert(l, fruitD[v], fruitW[v])
           } else {
               // restore full
               b.root = merge(l, r)
           }
       }
       bags[v] = b
   }
   // answer is sum in bag[1]
   ans := int64(0)
   if bags[1] != nil && bags[1].root != nil {
       ans = bags[1].root.sum
   }
   fmt.Println(ans)
}
