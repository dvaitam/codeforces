package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   x := make([]uint32, n)
   y := make([]uint32, n)
   for i := 0; i < n; i++ {
       var v uint32
       fmt.Fscan(reader, &v)
       x[i] = v
   }
   for i := 0; i < n; i++ {
       var v uint32
       fmt.Fscan(reader, &v)
       y[i] = v
   }

   const W = 32
   baseid := make([]int, W)
   basex := make([]uint32, W)
   base := make([]uint32, W)
   for i := 0; i < W; i++ {
       baseid[i] = -1
   }
   yy := make([]uint32, n)
   nosol := false

   // build basis for x
   for i := 0; i < n; i++ {
       p := x[i]
       var tmp uint32
       for j := 0; j < W; j++ {
           if (p>>j)&1 == 0 {
               continue
           }
           if baseid[j] != -1 {
               p ^= basex[j]
               tmp ^= base[j]
           } else {
               baseid[j] = i
               base[j] = tmp | (1 << uint(j))
               basex[j] = p
               break
           }
       }
   }
   // represent y in basis
   for i := 0; i < n; i++ {
       p := y[i]
       for j := 0; j < W; j++ {
           if (p>>j)&1 == 0 {
               continue
           }
           if baseid[j] != -1 {
               p ^= basex[j]
               yy[i] ^= base[j]
           } else {
               nosol = true
               break
           }
       }
   }
   if nosol {
       writer.WriteString("-1\n")
       return
   }

   type op struct{ u, v int }
   var ans, ans2 []op
   // non-basis vectors
   for i := 0; i < n; i++ {
       fl := true
       for j := 0; j < W; j++ {
           if baseid[j] == i {
               fl = false
               break
           }
       }
       if fl {
           ans = append(ans, op{i, i})
           for j := 0; j < W; j++ {
               if (yy[i]>>j)&1 == 1 {
                   ans = append(ans, op{i, baseid[j]})
               }
           }
       }
   }
   // basis vectors
   csl := make([]int, W)
   cpos := make([]int, W)
   cc := make([]uint32, W)
   ccn := 0
   for j := 0; j < W; j++ {
       if baseid[j] != -1 {
           csl[ccn] = j
           cpos[ccn] = baseid[j]
           cc[ccn] = yy[baseid[j]]
           ccn++
       }
   }
   for i := 0; i < ccn; i++ {
       // find pivot
       r := -1
       for j := i; j < ccn; j++ {
           if (cc[j]>>csl[i])&1 == 1 {
               r = j
               break
           }
       }
       if r == -1 {
           continue
       }
       if r != i {
           ans2 = append(ans2, op{cpos[r], cpos[i]})
           ans2 = append(ans2, op{cpos[i], cpos[r]})
           ans2 = append(ans2, op{cpos[r], cpos[i]})
           cc[i], cc[r] = cc[r], cc[i]
       }
       for j := 0; j < ccn; j++ {
           if j != i && (cc[j]>>csl[i])&1 == 1 {
               ans2 = append(ans2, op{cpos[j], cpos[i]})
               cc[j] ^= cc[i]
           }
       }
   }
   // fix missing pivots
   for i := 0; i < ccn; i++ {
       if (cc[i]>>csl[i])&1 == 0 {
           for j := 0; j < i; j++ {
               if (cc[j]>>csl[i])&1 == 1 {
                   ans = append(ans, op{cpos[j], cpos[i]})
               }
           }
           ans = append(ans, op{cpos[i], cpos[i]})
       }
   }
   total := len(ans) + len(ans2)
   fmt.Fprintln(writer, total)
   for _, p := range ans {
       fmt.Fprintf(writer, "%d %d\n", p.u+1, p.v+1)
   }
   for i := len(ans2) - 1; i >= 0; i-- {
       p := ans2[i]
       fmt.Fprintf(writer, "%d %d\n", p.u+1, p.v+1)
   }
}
