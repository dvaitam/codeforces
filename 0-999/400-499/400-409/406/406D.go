package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node for linked list representing convex hull
type Node struct {
   idx  int
   next *Node
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   x := make([]int64, n+1)
   y := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(in, &x[i], &y[i])
   }
   // next rightmost visible hill
   next := make([]int, n+1)
   // build upper convex hull of points to the right using linked list
   var head *Node
   for i := n; i >= 1; i-- {
       // insert at front
       node := &Node{idx: i, next: head}
       head = node
       // maintain convexity: while head, head.next, head.next.next exist
       for head.next != nil && head.next.next != nil {
           p := head.idx
           q := head.next.idx
           r := head.next.next.idx
           // if q is below line p->r, remove q (strict)
           if (y[q]-y[p])*(x[r]-x[q]) < (y[r]-y[q])*(x[q]-x[p]) {
               // remove q
               head.next = head.next.next
           } else {
               break
           }
       }
       if head.next != nil {
           next[i] = head.next.idx
       } else {
           next[i] = 0
       }
   }
   // compute depth to root (next==0)
   depth := make([]int, n+1)
   for i := n; i >= 1; i-- {
       if next[i] == 0 {
           depth[i] = 0
       } else {
           depth[i] = depth[next[i]] + 1
       }
   }
   // binary lifting
   const maxLog = 18
   up := make([][maxLog]int, n+1)
   for i := 1; i <= n; i++ {
       up[i][0] = next[i]
   }
   for k := 1; k < maxLog; k++ {
       for i := 1; i <= n; i++ {
           prev := up[i][k-1]
           if prev > 0 {
               up[i][k] = up[prev][k-1]
           } else {
               up[i][k] = 0
           }
       }
   }
   // process queries
   var m int
   fmt.Fscan(in, &m)
   res := make([]int, m)
   for qi := 0; qi < m; qi++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       u, v := a, b
       // align depths
       if depth[u] > depth[v] {
           diff := depth[u] - depth[v]
           for k := 0; k < maxLog; k++ {
               if diff&(1<<k) != 0 {
                   u = up[u][k]
               }
           }
       } else if depth[v] > depth[u] {
           diff := depth[v] - depth[u]
           for k := 0; k < maxLog; k++ {
               if diff&(1<<k) != 0 {
                   v = up[v][k]
               }
           }
       }
       // if meet or one is ancestor
       if u == v {
           res[qi] = u
           continue
       }
       // lift both
       for k := maxLog - 1; k >= 0; k-- {
           if up[u][k] != 0 && up[u][k] != up[v][k] {
               u = up[u][k]
               v = up[v][k]
           }
       }
       // now parents are same
       res[qi] = next[u]
   }
   // output
   for i, ans := range res {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, ans)
   }
   out.WriteByte('\n')
}
