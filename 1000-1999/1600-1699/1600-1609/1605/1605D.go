package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       adj := make([][]int, n+1)
       for i := 1; i < n; i++ {
           var u, v int
           fmt.Fscan(in, &u, &v)
           adj[u] = append(adj[u], v)
           adj[v] = append(adj[v], u)
       }
       // group values by highest bit
       maxBit := bits.Len(uint(n))
       vBits := make([][]int, maxBit)
       for i := 1; i <= n; i++ {
           b := bits.Len(uint(i)) - 1
           vBits[b] = append(vBits[b], i)
       }
       // BFS for depths
       depth := make([]int, n+1)
       parent := make([]int, n+1)
       parGroups := [2][]int{}
       queue := make([]int, n)
       qi, qj := 0, 0
       queue[qj] = 1
       qj++
       parent[1] = 0
       depth[1] = 0
       for qi < qj {
           x := queue[qi]
           qi++
           parGroups[depth[x]&1] = append(parGroups[depth[x]&1], x)
           for _, y := range adj[x] {
               if y == parent[x] {
                   continue
               }
               parent[y] = x
               depth[y] = depth[x] + 1
               queue[qj] = y
               qj++
           }
       }
       // choose smaller parity group as group0
       if len(parGroups[0]) > len(parGroups[1]) {
           parGroups[0], parGroups[1] = parGroups[1], parGroups[0]
       }
       cnt0 := len(parGroups[0])
       val := make([]int, n+1)
       cnts := [2]int{}
       // assign values
       for b := 0; b < maxBit; b++ {
           var tmp int
           if (cnt0>>b)&1 == 0 {
               tmp = 1
           } else {
               tmp = 0
           }
           for _, x := range vBits[b] {
               idx := cnts[tmp]
               val[parGroups[tmp][idx]] = x
               cnts[tmp]++
           }
       }
       // output
       for i := 1; i <= n; i++ {
           if i > 1 {
               out.WriteByte(' ')
           }
           fmt.Fprint(out, val[i])
       }
       out.WriteByte('\n')
   }
}
