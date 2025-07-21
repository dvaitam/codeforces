package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap node keyed by original index, maintains subtree size, sum and weighted sum
type Node struct {
   key      int
   val       int64
   prio     int
   left, right *Node
   size     int
   sum      int64 // sum of val in subtree
   wsum     int64 // weighted sum: sum(val * pos) in subtree
}

func update(n *Node) {
   if n == nil {
       return
   }
   n.size = 1
   n.sum = n.val
   n.wsum = n.val // will adjust below
   if n.left != nil {
       n.size += n.left.size
       n.sum += n.left.sum
       // left.wsum positions unchanged
       n.wsum += n.left.wsum + 0
   }
   // position of this node in subtree = size(left) + 1
   var leftSize int
   if n.left != nil {
       leftSize = n.left.size
   }
   // reset wsum: wsum = left.wsum + val*(leftSize+1) + right.wsum + right.sum*(leftSize+1)
   base := int64(leftSize + 1)
   w := int64(0)
   if n.left != nil {
       w += n.left.wsum
   }
   w += n.val * base
   if n.right != nil {
       w += n.right.wsum + n.right.sum*base
       n.size += n.right.size
       n.sum += n.right.sum
   }
   n.wsum = w
}

func merge(a, b *Node) *Node {
   if a == nil {
       return b
   }
   if b == nil {
       return a
   }
   if a.prio < b.prio {
       a.right = merge(a.right, b)
       update(a)
       return a
   }
   b.left = merge(a, b.left)
   update(b)
   return b
}

// split by key: left <= key, right > key
func split(n *Node, key int) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if n.key <= key {
       // go right
       ll, rr := split(n.right, key)
       n.right = ll
       update(n)
       return n, rr
   }
   // go left
   ll, rr := split(n.left, key)
   n.left = rr
   update(n)
   return ll, n
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   var k int64
   fmt.Fscan(in, &n, &k)
   a := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &a[i])
   }
   // build treap
   rand.Seed(time.Now().UnixNano())
   nodes := make([]*Node, n+1)
   for i := 1; i <= n; i++ {
       nodes[i] = &Node{key: i, val: a[i], prio: rand.Int()}
   }
   // build Cartesian tree on keys (increasing) by priority
   stack := make([]*Node, 0, n)
   var root *Node
   for i := 1; i <= n; i++ {
       curr := nodes[i]
       var last *Node
       for len(stack) > 0 && stack[len(stack)-1].prio > curr.prio {
           last = stack[len(stack)-1]
           stack = stack[:len(stack)-1]
       }
       curr.left = last
       if len(stack) > 0 {
           stack[len(stack)-1].right = curr
       }
       stack = append(stack, curr)
   }
   if len(stack) > 0 {
       root = stack[0]
       // find ultimate root
       for root != nil && root.left != nil {
           root = root.left
       }
       // but this left path may not cover entire, better get stack[0]'s ancestor chain? easier do nothing: stack[0] is root
       root = stack[0]
   }
   // update all nodes post-order
   var dfs func(n *Node)
   dfs = func(n *Node) {
       if n == nil {
           return
       }
       dfs(n.left)
       dfs(n.right)
       update(n)
   }
   dfs(root)
   // process removals
   m := n
   res := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       // split out left, mid, right
       var left, midNode, right *Node
       // left: keys <= i-1, mid+right: keys >= i
       left, right = split(root, i-1)
       // midNode: key == i, rest in right2
       var mid, right2 *Node
       mid, right2 = split(right, i)
       midNode = mid
       if midNode == nil {
           // already removed
           root = merge(left, right2)
           continue
       }
       // compute d_i
       leftSize := int64(0)
       leftSum := int64(0)
       leftWsum := int64(0)
       if left != nil {
           leftSize = int64(left.size)
           leftSum = left.sum
           leftWsum = left.wsum
       }
       pos := leftSize + 1
       m64 := int64(m)
       // sum b[j]*(j-1) = leftWsum - leftSum
       part := leftWsum - leftSum
       di := part - (pos-1)*(m64-pos)*midNode.val
       if di < k {
           // remove
           res = append(res, i)
           root = merge(left, right2)
           m--
       } else {
           // keep
           root = merge(merge(left, midNode), right2)
       }
   }
   // output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i, v := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
