package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

const INF64 = 1<<60

type Pair struct { prevJ, prevS int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       N := 1 << n
       sum := make([]int64, N)
       for i := 0; i < n; i++ {
           sum[1<<i] = a[i]
       }
       for S := 1; S < N; S++ {
           low := S & -S
           if low != S {
               sum[S] = sum[low] + sum[S^low]
           }
       }
       // DP array
       f := make([][][]int64, n+1)
       pre := make([][][]Pair, n+1)
       for i := 0; i <= n; i++ {
           f[i] = make([][]int64, n+1)
           pre[i] = make([][]Pair, n+1)
           for j := 0; j <= n; j++ {
               f[i][j] = make([]int64, N)
               pre[i][j] = make([]Pair, N)
               for S := 0; S < N; S++ {
                   f[i][j][S] = INF64
               }
           }
       }
       f[0][0][0] = 0
       fullMask := N - 1
       // transition
       for i := 0; i < n; i++ {
           for j := i; j < n; j++ {
               for S := 0; S < N; S++ {
                   cur := f[i][j][S]
                   if cur >= INF64 {
                       continue
                   }
                   rem := fullMask ^ S
                   for Tmask := rem; Tmask > 0; Tmask = (Tmask - 1) & rem {
                       if (Tmask>>j) == 0 || cur >= sum[Tmask] {
                           continue
                       }
                       shifted := Tmask >> j
                       p := bits.TrailingZeros(uint(shifted)) + j + 1
                       newS := S | Tmask
                       if f[i+1][p][newS] <= sum[Tmask] {
                           continue
                       }
                       f[i+1][p][newS] = sum[Tmask]
                       pre[i+1][p][newS] = Pair{prevJ: j, prevS: S}
                   }
               }
           }
       }
       // prepare index array
       idx := make([]int, n)
       for i := 0; i < n; i++ {
           idx[i] = i + 1
       }
       // trace function
       var trace func(i, j, S int)
       trace = func(i, j, S int) {
           if i == 0 {
               return
           }
           pr := pre[i][j][S]
           Tmask := S ^ pr.prevS
           for k := 0; k < n; k++ {
               if ((Tmask>>k)&1) == 1 && k != j-1 {
                   fmt.Fprintf(writer, "%d %d\n", idx[k], idx[j-1])
                   idx[k] = -1
                   for t := k + 1; t < n; t++ {
                       idx[t]--
                   }
               }
           }
           trace(i-1, pr.prevJ, pr.prevS)
       }
       // find answer
       for i := n; i >= 1; i-- {
           for j := i; j <= n; j++ {
               if f[i][j][fullMask] < INF64 {
                   fmt.Fprintln(writer, n-i)
                   trace(i, j, fullMask)
                   i = 0 // break outer
                   break
               }
           }
           if i == 0 {
               break
           }
       }
   }
}
