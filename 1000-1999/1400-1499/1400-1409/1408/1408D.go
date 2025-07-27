package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   c := make([]int, m)
   d := make([]int, m)
   for j := 0; j < m; j++ {
       fmt.Fscan(reader, &c[j], &d[j])
   }
   const MAX = 1000005
   mx := make([]int, MAX+5)
   maxRC := 0
   // Build constraints: for each robber i and searchlight j that can see it,
   // compute required moves right (rc) and up (rd).
   for i := 0; i < n; i++ {
       ai, bi := a[i], b[i]
       for j := 0; j < m; j++ {
           cj, dj := c[j], d[j]
           if cj >= ai && dj >= bi {
               rc := cj - ai + 1
               rd := dj - bi + 1
               if rd > mx[rc] {
                   mx[rc] = rd
               }
               if rc > maxRC {
                   maxRC = rc
               }
           }
       }
   }
   // Sweep x from maxRC down to 0, tracking maximum rd for rc > x
   ans := MAX * 2
   currMaxRD := 0
   for x := maxRC; x >= 0; x-- {
       if x+1 <= MAX {
           if mx[x+1] > currMaxRD {
               currMaxRD = mx[x+1]
           }
       }
       if x + currMaxRD < ans {
           ans = x + currMaxRD
       }
   }
   fmt.Println(ans)
}
