package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// DSU
var parent []int

func find(x int) int {
   if parent[x] < 0 {
       return x
   }
   parent[x] = find(parent[x])
   return parent[x]
}

func unite(a, b int) {
   a = find(a)
   b = find(b)
   if a == b {
       return
   }
   // union by size
   if parent[a] > parent[b] {
       a, b = b, a
   }
   parent[a] += parent[b]
   parent[b] = a
}

// Treap for keys (x,id)
type Key struct{ x, id int }
type Node struct {
   key      Key
   pri      int
   left, right *Node
}

func split(root *Node, key Key) (l, r *Node) {
   if root == nil {
       return nil, nil
   }
   if less(root.key, key) {
       // root.key < key => keep root in l
       ll, rr := split(root.right, key)
       root.right = ll
       return root, rr
   } else {
       ll, rr := split(root.left, key)
       root.left = rr
       return ll, root
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
   }
   b.left = merge(a, b.left)
   return b
}

func less(a, b Key) bool {
   if a.x != b.x {
       return a.x < b.x
   }
   return a.id < b.id
}

func insert(root *Node, node *Node) *Node {
   if root == nil {
       return node
   }
   if node.pri < root.pri {
       l, r := split(root, node.key)
       node.left, node.right = l, r
       return node
   }
   if less(node.key, root.key) {
       root.left = insert(root.left, node)
   } else {
       root.right = insert(root.right, node)
   }
   return root
}

func erase(root *Node, key Key) *Node {
   if root == nil {
       return nil
   }
   if root.key == key {
       return merge(root.left, root.right)
   }
   if less(key, root.key) {
       root.left = erase(root.left, key)
   } else {
       root.right = erase(root.right, key)
   }
   return root
}

// traverse and collect ids
func collect(root *Node, ids *[]int) {
   if root == nil {
       return
   }
   collect(root.left, ids)
   *ids = append(*ids, root.key.id)
   collect(root.right, ids)
}

func main() {
   rand.Seed(time.Now().UnixNano())
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   parent = make([]int, n+5)
   for i := range parent {
       parent[i] = -1
   }
   // treaps
   var rootL, rootR *Node
   // store intervals
   lval := make([]int, n+5)
   rval := make([]int, n+5)
   cur := 0
   for i := 0; i < n; i++ {
       var t int
       fmt.Fscan(in, &t)
       if t == 1 {
           cur++
           x, y := 0, 0
           fmt.Fscan(in, &x, &y)
           lval[cur] = x
           rval[cur] = y
           // Step1: find in rootL keys in [x,y)
           // split <(x,0) and >=
           a, b := split(rootL, Key{x, 0})
           // split b into c <(y,0) and d
           c, d := split(b, Key{y, 0})
           // c are nodes with l in [x,y)
           ids := []int{}
           collect(c, &ids)
           // remove those from rootR and union
           for _, id := range ids {
               unite(cur, id)
               rootR = erase(rootR, Key{rval[id], id})
           }
           // merge back rootL = merge(a,d)
           rootL = merge(a, d)
           // Step2: in rootR keys with r in (x,y)
           // split <=(x,inf) and >
           a2, b2 := split(rootR, Key{x+1, 0})
           // split b2 into c2 <(y,0) and d2
           c2, d2 := split(b2, Key{y, 0})
           ids = ids[:0]
           collect(c2, &ids)
           for _, id := range ids {
               unite(cur, id)
               rootL = erase(rootL, Key{lval[id], id})
           }
           rootR = merge(a2, d2)
           // finally insert new into both
           nodeL := &Node{key: Key{x, cur}, pri: rand.Int()}
           rootL = insert(rootL, nodeL)
           nodeR := &Node{key: Key{y, cur}, pri: rand.Int()}
           rootR = insert(rootR, nodeR)
       } else {
           a, b := 0, 0
           fmt.Fscan(in, &a, &b)
           if find(a) == find(b) {
               fmt.Fprintln(out, "YES")
           } else {
               fmt.Fprintln(out, "NO")
           }
       }
   }
}
