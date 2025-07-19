package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   rdr = bufio.NewReader(os.Stdin)
   wrt = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   sign, x := 1, 0
   var c byte
   c, _ = rdr.ReadByte()
   for (c < '0' || c > '9') && c != '-' {
       c, _ = rdr.ReadByte()
   }
   if c == '-' {
       sign = -1
       c, _ = rdr.ReadByte()
   }
   for c >= '0' && c <= '9' {
       x = x*10 + int(c-'0')
       c, _ = rdr.ReadByte()
   }
   return x * sign
}

func main() {
   defer wrt.Flush()
   n := readInt()
   m := readInt()
   p := readInt()
   lnU := make([]int, m+1)
   lnV := make([]int, m+1)
   adj := make([][]int, n+1)
   deg := make([]int, n+1)
   activeEdge := make([]bool, m+1)
   for i := 1; i <= m; i++ {
       u := readInt()
       v := readInt()
       lnU[i], lnV[i] = u, v
       adj[u] = append(adj[u], i)
       adj[v] = append(adj[v], i)
       deg[u]++
       deg[v]++
       activeEdge[i] = true
   }
   removed := make([]bool, n+1)
   ans := n
   ansArr := make([]int, m+2)
   // queue for nodes
   queue := make([]int, 0, n)

   enqueue := func(v int) {
       queue = append(queue, v)
   }

   // delete node v and update neighbors
   deleteNode := func(v int) {
       if removed[v] {
           return
       }
       removed[v] = true
       ans--
       for _, ei := range adj[v] {
           if !activeEdge[ei] {
               continue
           }
           activeEdge[ei] = false
           u := lnU[ei]
           w := lnV[ei]
           other := u
           if other == v {
               other = w
           }
           if removed[other] {
               continue
           }
           deg[other]--
           if deg[other] < p {
               enqueue(other)
           }
       }
   }

   // initial prune
   for i := 1; i <= n; i++ {
       if deg[i] < p {
           enqueue(i)
       }
   }
   for head := 0; head < len(queue); head++ {
       deleteNode(queue[head])
   }
   // clear queue before edge removal phase
   queue = queue[:0]

   // process removals in reverse edge order
   for i := m; i >= 1; i-- {
       ansArr[i] = ans
       // remove edge i
       if activeEdge[i] {
           u := lnU[i]
           v := lnV[i]
           // decrement degrees if both endpoints alive
           activeEdge[i] = false
           if !removed[u] {
               deg[u]--
               if deg[u] < p {
                   enqueue(u)
               }
           }
           if !removed[v] {
               deg[v]--
               if deg[v] < p {
                   enqueue(v)
               }
           }
           for head := 0; head < len(queue); head++ {
               deleteNode(queue[head])
           }
           // reset queue
           queue = queue[:0]
       }
   }
   // output
   for i := 1; i <= m; i++ {
       fmt.Fprintln(wrt, ansArr[i])
   }
}
