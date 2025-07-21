package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "time"
)

// Treap node
type Node struct {
   key        int
   pri        int
   left, right *Node
   size       int
   ans        [5]int64
}

// update recalculates size and ans for node n
func update(n *Node) {
   if n == nil {
       return
   }
   // reset
   n.size = 1
   for i := 0; i < 5; i++ {
       n.ans[i] = 0
   }
   // left subtree
   lsz := 0
   if n.left != nil {
       lsz = n.left.size
       n.size += n.left.size
       for i := 0; i < 5; i++ {
           n.ans[i] += n.left.ans[i]
       }
   }
   // this node at position lsz (0-based)
   mod := lsz % 5
   n.ans[mod] += int64(n.key)
   // right subtree
   if n.right != nil {
       // shift right.ans by (lsz+1)%5
       shift := (lsz + 1) % 5
       for i := 0; i < 5; i++ {
           j := (i + shift) % 5
           n.ans[j] += n.right.ans[i]
       }
       n.size += n.right.size
   }
}

// merge combines two treaps l and r
func merge(l, r *Node) *Node {
   if l == nil {
       return r
   }
   if r == nil {
       return l
   }
   if l.pri < r.pri {
       l.right = merge(l.right, r)
       update(l)
       return l
   }
   r.left = merge(l, r.left)
   update(r)
   return r
}

// split splits treap n into (< key) and (>= key)
func split(n *Node, key int) (l, r *Node) {
   if n == nil {
       return nil, nil
   }
   if n.key < key {
       // go right
       lr, rr := split(n.right, key)
       n.right = lr
       update(n)
       return n, rr
   }
   // n.key >= key
   ll, rl := split(n.left, key)
   n.left = rl
   update(n)
   return ll, n
}

func main() {
   rand.Seed(time.Now().UnixNano())
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var root *Node
   for i := 0; i < n; i++ {
       var op string
       fmt.Fscan(reader, &op)
       switch op {
       case "add":
           var x int
           fmt.Fscan(reader, &x)
           // insert x
           a, b := split(root, x)
           node := &Node{key: x, pri: rand.Int()}
           // initialize node
           node.size = 1
           node.ans[0] = int64(x)
           root = merge(merge(a, node), b)
       case "del":
           var x int
           fmt.Fscan(reader, &x)
           // remove x
           a, bc := split(root, x)
           b, c := split(bc, x+1)
           // b is the node with key x
           root = merge(a, c)
       case "sum":
           if root == nil {
               fmt.Fprintln(writer, 0)
           } else {
               // positions with j mod 5 == 2 (0-based) => median sum
               fmt.Fprintln(writer, root.ans[2])
           }
       }
   }
}
