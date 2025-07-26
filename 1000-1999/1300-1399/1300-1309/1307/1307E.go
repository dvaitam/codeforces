package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

type Cow struct {
   posL, posR int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   s := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // occurrences of sweetness
   occ := make([][]int, n+1)
   for i, v := range s {
       occ[v] = append(occ[v], i+1)
   }
   // cows grouped by flavor
   cowsByF := make([][]Cow, n+1)
   for i := 0; i < m; i++ {
       var f, h int
       fmt.Fscan(reader, &f, &h)
       list := occ[f]
       sz := len(list)
       var posL, posR int
       if h <= sz {
           posL = list[h-1]
           posR = list[sz-h]
       } else {
           posL = -1
           posR = -1
       }
       if posL < 0 && posR < 0 {
           continue
       }
       cowsByF[f] = append(cowsByF[f], Cow{posL, posR})
   }
   // precompute powers of 2
   pow2 := make([]int, m+1)
   pow2[0] = 1
   for i := 1; i <= m; i++ {
       pow2[i] = pow2[i-1] * 2 % MOD
   }
   bestCnt := 0
   bestWays := 0
   // try boundary k from 0 to n
   for k := 0; k <= n; k++ {
       total := 0
       ways := 1
       for f := 1; f <= n; f++ {
           list := cowsByF[f]
           if len(list) == 0 {
               continue
           }
           ca, cb, inter := 0, 0, 0
           for _, c := range list {
               inA := c.posL > 0 && c.posL <= k
               inB := c.posR > 0 && c.posR > k
               if inA {
                   ca++
               }
               if inB {
                   cb++
               }
               if inA && inB {
                   inter++
               }
           }
           uni := ca + cb - inter
           if uni == 0 {
               continue
           }
           total += uni
           if inter > 0 {
               ways = ways * pow2[inter] % MOD
           }
       }
       if total > bestCnt {
           bestCnt = total
           bestWays = ways
       } else if total == bestCnt {
           bestWays = (bestWays + ways) % MOD
       }
   }
   fmt.Fprintln(writer, bestCnt, bestWays)
}
