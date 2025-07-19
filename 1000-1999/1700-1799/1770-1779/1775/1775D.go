package main

import (
   "bufio"
   "fmt"
   "os"
)

// Node holds a BFS state: node index, distance, and previous node
type Node struct {
   w, d, prev int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   mx := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i] > mx {
           mx = a[i]
       }
   }
   var s, t int
   fmt.Fscan(reader, &s, &t)
   s--
   t--

   // Sieve for minimum prime factors
   minp := make([]int, mx+1)
   primes := make([]int, 0)
   for i := 2; i <= mx; i++ {
       if minp[i] == 0 {
           minp[i] = i
           primes = append(primes, i)
       }
       for _, p := range primes {
           prod := p * i
           if prod > mx {
               break
           }
           minp[prod] = p
           if p == minp[i] {
               break
           }
       }
   }

   // Build adjacency: prime nodes to indices
   edg := make([][]int, mx+1)
   for i := 0; i < n; i++ {
       tmp := a[i]
       for tmp > 1 {
           p := minp[tmp]
           edg[p] = append(edg[p], i)
           for tmp%p == 0 {
               tmp /= p
           }
       }
   }

   total := n + mx + 1
   dist := make([]int, total)
   nxt := make([]int, total)
   for i := range dist {
       dist[i] = -1
       nxt[i] = -1
   }

   // BFS from t
   queue := make([]Node, 0, total)
   queue = append(queue, Node{t, 0, -1})
   head := 0
   for head < len(queue) {
       cur := queue[head]
       head++
       w, d, pr := cur.w, cur.d, cur.prev
       if dist[w] != -1 {
           continue
       }
       dist[w] = d
       nxt[w] = pr
       if w < n {
           tmp := a[w]
           for tmp > 1 {
               p := minp[tmp]
               nb := n + p
               if dist[nb] == -1 {
                   queue = append(queue, Node{nb, d + 1, w})
               }
               for tmp%p == 0 {
                   tmp /= p
               }
           }
       } else {
           p := w - n
           for _, idx := range edg[p] {
               if dist[idx] == -1 {
                   queue = append(queue, Node{idx, d + 1, w})
               }
           }
       }
   }

   // Output path
   if dist[s] == -1 {
       fmt.Fprintln(writer, -1)
       return
   }
   length := dist[s]/2 + 1
   path := make([]int, 0, length)
   path = append(path, s)
   for i := nxt[s]; i != -1; i = nxt[i] {
       if i < n {
           path = append(path, i)
       }
   }
   fmt.Fprintln(writer, length)
   for i, v := range path {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, v+1)
   }
   fmt.Fprintln(writer)
}
