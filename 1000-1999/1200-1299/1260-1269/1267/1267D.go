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

   var n int
   if _, err := fmt.Fscan(rdr, &n); err != nil {
       return
   }
   ne := make([]int, n+1)
   ct := make([]int, n+1)
   CT := make([]int, n+1)
   var Cnt [8]int
   var To [8]int
   var From [8][2]int
   type edge struct{u, v int}
   var edges []edge

   // read initial masks
   for i := 1; i <= n; i++ {
       var a, b, c int
       fmt.Fscan(rdr, &a, &b, &c)
       ne[i] = (a << 2) | (b << 1) | c
       Cnt[ne[i]]++
       if ne[i] != 1 && ne[i] != 2 && ne[i] != 4 {
           CT[i] = 1
       } else {
           CT[i] = 0
       }
   }
   // read target masks
   for i := 1; i <= n; i++ {
       var a, b, c int
       fmt.Fscan(rdr, &a, &b, &c)
       ct[i] = (a << 2) | (b << 1) | c
       To[ne[i]] |= 1 << (ct[i] & ne[i])
   }

   // compute reachable states
   sta := 1 << 7
   for {
       nwS := sta
       // single transitions
       for i := 0; i < 8; i++ {
           if (sta>>i)&1 == 1 {
               for j := 0; j < 8; j++ {
                   if (To[i]>>j)&1 == 1 && (nwS>>j)&1 == 0 {
                       From[j][0], From[j][1] = i, -1
                       nwS |= 1 << j
                   }
               }
           }
       }
       // combine two dimensions
       for d1 := 0; d1 < 3; d1++ {
           for s1 := 0; s1 < 8; s1++ {
               if (sta>>s1)&1 == 1 && (To[s1]>>(1<<d1))&1 == 1 {
                   for d2 := d1 + 1; d2 < 3; d2++ {
                       for s2 := 0; s2 < 8; s2++ {
                           if (sta>>s2)&1 == 1 && (To[s2]>>(1<<d2))&1 == 1 {
                               nb := (1<<d1)|(1<<d2)
                               if (nwS>>nb)&1 == 0 {
                                   From[nb][0], From[nb][1] = s1, s2
                                   nwS |= 1 << nb
                               }
                           }
                       }
                   }
               }
           }
       }
       if nwS == sta {
           break
       }
       sta = nwS
   }
   // check feasibility
   for i := 1; i < 8; i++ {
       if Cnt[i] > 0 && (sta>>i)&1 == 0 {
           fmt.Fprintln(wtr, "Impossible")
           return
       }
   }
   // build edges
   add := func(u, v int) {
       edges = append(edges, edge{u, v})
   }
   for i := 1; i < 7; i++ {
       f0, f1 := From[i][0], From[i][1]
       if f1 == -1 {
           // single source
           var id0 int
           for j := 1; j <= n; j++ {
               if ne[j] == f0 && (ne[j]&ct[j]) == i {
                   id0 = j
                   break
               }
           }
           for j := 1; j <= n; j++ {
               if ne[j] == i {
                   add(id0, j)
               }
           }
       } else {
           // two sources
           x := i & -i
           y := i - x
           var id1, id2, pp int
           for j := 1; j <= n; j++ {
               if ne[j] == f0 && (ne[j]&ct[j]) == x {
                   id1 = j
                   break
               }
           }
           for j := 1; j <= n; j++ {
               if ne[j] == f1 && (ne[j]&ct[j]) == y {
                   id2 = j
                   break
               }
           }
           for j := 1; j <= n; j++ {
               if ne[j] == i {
                   pp = j
                   break
               }
           }
           if pp != 0 {
               CT[pp] = 0
               add(id1, pp)
               add(id2, pp)
           }
           for j := 1; j <= n; j++ {
               if ne[j] == i && j != pp {
                   add(pp, j)
               }
           }
       }
   }
   // edges for full mask
   for i := 2; i <= n; i++ {
       if ne[i] == 7 {
           add(1, i)
       }
   }
   // output
   fmt.Fprintln(wtr, "Possible")
   for i := 1; i <= n; i++ {
       if i > 1 {
           wtr.WriteByte(' ')
       }
       fmt.Fprint(wtr, CT[i])
   }
   fmt.Fprintln(wtr)
   fmt.Fprintln(wtr, len(edges))
   for _, e := range edges {
       fmt.Fprintln(wtr, e.u, e.v)
   }
}
