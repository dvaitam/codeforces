package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   wtr := bufio.NewWriter(os.Stdout)
   defer wtr.Flush()
   var t int
   fmt.Fscan(rdr, &t)
   for t > 0 {
       t--
       solve(rdr, wtr)
   }
}

func solve(rdr *bufio.Reader, wtr *bufio.Writer) {
   var n, m int
   fmt.Fscan(rdr, &n, &m)
   k := (m + 63) / 64
   // bitsets per row
   inSet := make([][]uint64, n)
   notInSet := make([][]uint64, n)
   for i := 0; i < n; i++ {
       inSet[i] = make([]uint64, k)
       notInSet[i] = make([]uint64, k)
   }
   needCheck := make([]bool, n)
   for i := range needCheck {
       needCheck[i] = true
   }
   degForSet := make([]int, m)
   G := make([][]bool, n)
   for i := 0; i < n; i++ {
       G[i] = make([]bool, m)
   }
   var s string
   for j := 0; j < m; j++ {
       fmt.Fscan(rdr, &s)
       for i := 0; i < n && i < len(s); i++ {
           if s[i] == '1' {
               G[i][j] = true
               degForSet[j]++
           }
       }
       if degForSet[j] <= 1 {
           continue
       }
       x := j >> 6
       z := uint64(1) << (uint(j) & 63)
       for i := 0; i < n; i++ {
           if G[i][j] {
               inSet[i][x] |= z
           } else {
               notInSet[i][x] |= z
           }
       }
   }
   alive := make([]int, n)
   for i := 0; i < n; i++ {
       alive[i] = i
   }
   type pair struct{ v, u int }
   ans := make([]pair, 0, n)
   for len(alive) > 2 {
       fnd := false
       for idx := 0; idx < len(alive); idx++ {
           v := alive[idx]
           if !needCheck[v] {
               continue
           }
           needCheck[v] = false
           p := -1
           // find u
           for _, u := range alive {
               if u == v {
                   continue
               }
               ok := true
               for b := 0; b < k; b++ {
                   if inSet[v][b]&notInSet[u][b] != 0 {
                       ok = false
                       break
                   }
               }
               if ok {
                   p = u
                   break
               }
           }
           if p < 0 {
               continue
           }
           fnd = true
           ans = append(ans, pair{v, p})
           // kill sets
           for j := 0; j < m; j++ {
               if !G[v][j] {
                   continue
               }
               degForSet[j]--
               if degForSet[j] == 1 {
                   id := j
                   x := id >> 6
                   z := uint64(1) << (uint(id) & 63)
                   for i := 0; i < n; i++ {
                       if G[i][id] {
                           inSet[i][x] ^= z
                           needCheck[i] = true
                       }
                   }
               }
           }
           // remove v
           alive[idx] = alive[len(alive)-1]
           alive = alive[:len(alive)-1]
           break
       }
       if !fnd {
           fmt.Fprintln(wtr, "NO")
           return
       }
   }
   if len(alive) == 2 {
       ans = append(ans, pair{alive[0], alive[1]})
   }
   fmt.Fprintln(wtr, "YES")
   for _, p := range ans {
       fmt.Fprintln(wtr, p.v+1, p.u+1)
   }
}
