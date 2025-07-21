package main

import (
   "bufio"
   "fmt"
   "math/rand"
   "os"
   "io"
   "time"
)

type Node struct {
   left, right *Node
   sz          int
   label       int
   prio        int
}

func upd(n *Node) {
   n.sz = 1
   if n.left != nil {
       n.sz += n.left.sz
   }
   if n.right != nil {
       n.sz += n.right.sz
   }
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
       upd(a)
       return a
   }
   b.left = merge(a, b.left)
   upd(b)
   return b
}

// split first k nodes into a, rest into b
func split(n *Node, k int) (a, b *Node) {
   if n == nil {
       return nil, nil
   }
   var lsz int
   if n.left != nil {
       lsz = n.left.sz
   }
   if k <= lsz {
       a, n.left = split(n.left, k)
       if n.left != nil {
           upd(n.left)
       }
       upd(n)
       return a, n
   }
   // k > lsz
   n.right, b = split(n.right, k-lsz-1)
   if n.right != nil {
       upd(n.right)
   }
   upd(n)
   return n, b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   data, _ := io.ReadAll(reader)
   var pos int
   nextInt := func() int {
       for pos < len(data) && (data[pos] < '0' || data[pos] > '9') && data[pos] != '-' {
           pos++
       }
       sign := 1
       if pos < len(data) && data[pos] == '-' {
           sign = -1
           pos++
       }
       x := 0
       for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
           x = x*10 + int(data[pos]-'0')
           pos++
       }
       return x * sign
   }
   n := nextInt()
   m := nextInt()
   // build initial treap of n unnamed nodes
   rand.Seed(time.Now().UnixNano())
   var root *Node
   for i := 0; i < n; i++ {
       nd := &Node{label: 0, prio: rand.Int()}
       nd.sz = 1
       root = merge(root, nd)
   }
   assigned := make([]bool, n+1)
   for i := 0; i < m; i++ {
       x := nextInt()
       y := nextInt()
       if y < 1 || y > root.sz {
           fmt.Fprintln(writer, -1)
           return
       }
       // split out y-th node
       t1, t2 := split(root, y-1)
       tmid, t3 := split(t2, 1)
       if tmid == nil {
           fmt.Fprintln(writer, -1)
           return
       }
       // check assignment
       if tmid.label != 0 && tmid.label != x {
           fmt.Fprintln(writer, -1)
           return
       }
       if tmid.label == 0 {
           if assigned[x] {
               fmt.Fprintln(writer, -1)
               return
           }
           tmid.label = x
           assigned[x] = true
       }
       // move to front: merge tmid + (t1+t3)
       tmp := merge(t1, t3)
       root = merge(tmid, tmp)
   }
   // inorder traverse to collect initial permutation
   res := make([]int, 0, n)
   var stack []*Node
   cur := root
   curLabel := 1
   // traverse in-order
   for cur != nil || len(stack) > 0 {
       for cur != nil {
           stack = append(stack, cur)
           cur = cur.left
       }
       node := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if node.label == 0 {
           for curLabel <= n && assigned[curLabel] {
               curLabel++
           }
           node.label = curLabel
           assigned[curLabel] = true
       }
       res = append(res, node.label)
       cur = node.right
   }
   // output
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
