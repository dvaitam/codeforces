package main

import (
   "bufio"
   "fmt"
   "os"
)

const threshold = 600

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   fmt.Fscan(reader, &n, &m, &q)
   status := make([]bool, n+1)
   var o int
   fmt.Fscan(reader, &o)
   for i := 0; i < o; i++ {
       var x int
       fmt.Fscan(reader, &x)
       status[x] = true
   }
   adj := make([]map[int]struct{}, n+1)
   for i := 1; i <= n; i++ {
       adj[i] = make(map[int]struct{})
   }
   for i := 0; i < m; i++ {
       var u, v int
       fmt.Fscan(reader, &u, &v)
       adj[u][v] = struct{}{}
       adj[v][u] = struct{}{}
   }
   isHeavy := make([]bool, n+1)
   heavyNeighbors := make([][]int, n+1)
   countOnline := make([]int, n+1)
   // initial heavy marking
   for u := 1; u <= n; u++ {
       if len(adj[u]) > threshold {
           isHeavy[u] = true
       }
   }
   // build heavyNeighbors and countOnline for initial heavies
   for u := 1; u <= n; u++ {
       if isHeavy[u] {
           cnt := 0
           for v := range adj[u] {
               if status[v] {
                   cnt++
               }
               heavyNeighbors[v] = append(heavyNeighbors[v], u)
           }
           countOnline[u] = cnt
           for v := range adj[u] {
               if isHeavy[v] {
                   heavyNeighbors[u] = append(heavyNeighbors[u], v)
               }
           }
       }
   }
   // promote function
   promote := func(u int) {
       isHeavy[u] = true
       cnt := 0
       for v := range adj[u] {
           if status[v] {
               cnt++
           }
           heavyNeighbors[v] = append(heavyNeighbors[v], u)
       }
       countOnline[u] = cnt
       for v := range adj[u] {
           if isHeavy[v] {
               heavyNeighbors[u] = append(heavyNeighbors[u], v)
           }
       }
   }
   // process queries
   for i := 0; i < q; i++ {
       var op byte
       fmt.Fscan(reader, &op)
       switch op {
       case 'O':
           var u int
           fmt.Fscan(reader, &u)
           status[u] = true
           for _, h := range heavyNeighbors[u] {
               countOnline[h]++
           }
       case 'F':
           var u int
           fmt.Fscan(reader, &u)
           status[u] = false
           for _, h := range heavyNeighbors[u] {
               countOnline[h]--
           }
       case 'A':
           var u, v int
           fmt.Fscan(reader, &u, &v)
           adj[u][v] = struct{}{}
           adj[v][u] = struct{}{}
           if isHeavy[u] {
               if status[v] {
                   countOnline[u]++
               }
               heavyNeighbors[v] = append(heavyNeighbors[v], u)
           }
           if isHeavy[v] {
               if status[u] {
                   countOnline[v]++
               }
               heavyNeighbors[u] = append(heavyNeighbors[u], v)
           }
           if !isHeavy[u] && len(adj[u]) > threshold {
               promote(u)
           }
           if !isHeavy[v] && len(adj[v]) > threshold {
               promote(v)
           }
       case 'D':
           var u, v int
           fmt.Fscan(reader, &u, &v)
           delete(adj[u], v)
           delete(adj[v], u)
           if isHeavy[u] {
               if status[v] {
                   countOnline[u]--
               }
               arr := heavyNeighbors[v]
               for idx, h := range arr {
                   if h == u {
                       arr[idx] = arr[len(arr)-1]
                       heavyNeighbors[v] = arr[:len(arr)-1]
                       break
                   }
               }
           }
           if isHeavy[v] {
               if status[u] {
                   countOnline[v]--
               }
               arr := heavyNeighbors[u]
               for idx, h := range arr {
                   if h == v {
                       arr[idx] = arr[len(arr)-1]
                       heavyNeighbors[u] = arr[:len(arr)-1]
                       break
                   }
               }
           }
       case 'C':
           var u int
           fmt.Fscan(reader, &u)
           if isHeavy[u] {
               fmt.Fprintln(writer, countOnline[u])
           } else {
               cnt := 0
               for v := range adj[u] {
                   if status[v] {
                       cnt++
                   }
               }
               fmt.Fprintln(writer, cnt)
           }
       }
   }
}
