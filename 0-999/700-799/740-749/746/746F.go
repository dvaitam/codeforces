package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
)

type Node struct {
   delta int
   id    int
   pri   int
   left, right *Node
   size  int
   sum   int64
}

func upd(n *Node) {
   n.size = 1
   n.sum = int64(n.delta)
   if n.left != nil {
       n.size += n.left.size
       n.sum += n.left.sum
   }
   if n.right != nil {
       n.size += n.right.size
       n.sum += n.right.sum
   }
}

// merge t1 (<) and t2 (>=)
func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.pri > b.pri {
       a.right = merge(a.right, b)
       upd(a)
       return a
   }
   b.left = merge(a, b.left)
   upd(b)
   return b
}

// split by key < (delta,id)
func split(root *Node, delta, id int) (l, r *Node) {
   if root == nil {
       return nil, nil
   }
   if root.delta < delta || (root.delta == delta && root.id < id) {
       ll, rr := split(root.right, delta, id)
       root.right = ll
       upd(root)
       return root, rr
   }
   ll, rr := split(root.left, delta, id)
   root.left = rr
   upd(root)
   return ll, root
}

// insert node
func insert(root *Node, node *Node) *Node {
   if root == nil {
       return node
   }
   if node.pri > root.pri {
       l, r := split(root, node.delta, node.id)
       node.left = l
       node.right = r
       upd(node)
       return node
   }
   if node.delta < root.delta || (node.delta == root.delta && node.id < root.id) {
       root.left = insert(root.left, node)
   } else {
       root.right = insert(root.right, node)
   }
   upd(root)
   return root
}

// erase key
func erase(root *Node, delta, id int) *Node {
   if root == nil {
       return nil
   }
   if root.delta == delta && root.id == id {
       return merge(root.left, root.right)
   }
   if delta < root.delta || (delta == root.delta && id < root.id) {
       root.left = erase(root.left, delta, id)
   } else {
       root.right = erase(root.right, delta, id)
   }
   upd(root)
   return root
}

// sum of k largest deltas
func sumTopK(root *Node, k int) int64 {
   if root == nil || k <= 0 {
       return 0
   }
   // size of right subtree
   rs := 0
   if root.right != nil {
       rs = root.right.size
   }
   if rs >= k {
       return sumTopK(root.right, k)
   }
   // include right + this node
   res := int64(0)
   if root.right != nil {
       res += root.right.sum
   }
   if rs+1 <= k {
       res += int64(root.delta)
       // remaining from left
       res += sumTopK(root.left, k-rs-1)
   }
   return res
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, w int
   var kTime int64
   fmt.Fscan(in, &n, &w, &kTime)
   a := make([]int64, n)
   t := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &t[i])
   }
   del := make([]int, n)
   for i := 0; i < n; i++ {
       del[i] = int(t[i] / 2)
   }
   var root *Node
   var sumT int64
   var sumA int64
   var ans int64
   l := 0
   windowSize := 0
   for r := 0; r < n; r++ {
       // add r
       sumT += t[r]
       sumA += a[r]
       node := &Node{delta: del[r], id: r, pri: rand.Int(), size: 1, sum: int64(del[r])}
       root = insert(root, node)
       windowSize++
       // shrink while too long
       for {
           k2 := w
           if windowSize < k2 {
               k2 = windowSize
           }
           sumD := sumTopK(root, k2)
           if sumT - sumD <= kTime {
               break
           }
           // remove l
           sumT -= t[l]
           sumA -= a[l]
           root = erase(root, del[l], l)
           windowSize--
           l++
       }
       if sumA > ans {
           ans = sumA
       }
   }
   fmt.Println(ans)
}
