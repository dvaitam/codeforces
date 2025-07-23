package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func add(a, b int) int {
   a += b
   if a >= MOD {
       a -= MOD
   }
   return a
}
func sub(a, b int) int {
   a -= b
   if a < 0 {
       a += MOD
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n, m int
   fmt.Fscan(reader, &n, &m)
   // Precompute countA[w][h]: number of monotonic shapes (Lmin=1,Rmax=w) of height h and width w
   // countA dimensions [0..m][0..n]
   countA := make([][]int, m+1)
   for w := range countA {
       countA[w] = make([]int, n+1)
   }
   // Temporary dp arrays max size 150
   const MAX = 150
   var dpPrev, dpNext [MAX + 2][MAX + 2]int
   var rowCum [MAX + 2][MAX + 2]int
   var prefCum [MAX + 2][MAX + 2]int
   var diagCum [MAX + 2]int
   for w := 1; w <= m; w++ {
       // initialize dpPrev for h=1: all segments [l..r]
       for l := 1; l <= w; l++ {
           for r := 1; r <= w; r++ {
               if l <= r {
                   dpPrev[l][r] = 1
               } else {
                   dpPrev[l][r] = 0
               }
           }
       }
       // h=1
       countA[w][1] = w * (w + 1) / 2 % MOD
       // for h from 2..n
       for h := 2; h <= n; h++ {
           // compute rowCum[l][r] = sum_{k=l..r} dpPrev[l][k]
           for l := 1; l <= w; l++ {
               var sum int
               // rowCum[l][l-1] = 0
               rowCum[l][l-1] = 0
               sum = 0
               for r := l; r <= w; r++ {
                   sum = add(sum, dpPrev[l][r])
                   rowCum[l][r] = sum
               }
           }
           // prefCum[r][i] = sum_{l=1..i} rowCum[l][r]
           for r := 1; r <= w; r++ {
               var sum int
               for i := 1; i <= w; i++ {
                   sum = add(sum, rowCum[i][r])
                   prefCum[r][i] = sum
               }
           }
           // diagCum[i] = sum_{l=1..i} rowCum[l][l-1]
           diagCum[0] = 0
           for i := 1; i <= w; i++ {
               diagCum[i] = add(diagCum[i-1], rowCum[i][i-1])
           }
           // compute dpNext
           for l := 1; l <= w; l++ {
               for r := 1; r <= w; r++ {
                   if l > r {
                       dpNext[l][r] = 0
                       continue
                   }
                   // sum over l0 in [l..r] of (rowCum[l0][r] - rowCum[l0][l0-1])
                   // = (prefCum[r][r] - prefCum[r][l-1]) - (diagCum[r] - diagCum[l-1])
                   a := prefCum[r][r]
                   if l-1 >= 1 {
                       a = sub(a, prefCum[r][l-1])
                   }
                   b := sub(diagCum[r], diagCum[l-1])
                   dpNext[l][r] = sub(a, b)
               }
           }
           // swap dpPrev, dpNext
           for l := 1; l <= w; l++ {
               for r := l; r <= w; r++ {
                   dpPrev[l][r] = dpNext[l][r]
               }
           }
           // total for this h
           var total int
           for l := 1; l <= w; l++ {
               for r := l; r <= w; r++ {
                   total = add(total, dpPrev[l][r])
               }
           }
           countA[w][h] = total
       }
   }
   // accumulate answer
   var ans int
   for h := 1; h <= n; h++ {
       for w := 1; w <= m; w++ {
           // DP1_valid = countA[w][h] -2*countA[w-1][h] + countA[w-2][h]
           v := countA[w][h]
           if w-1 >= 1 {
               v = sub(v, (countA[w-1][h]*2)%MOD)
           }
           if w-2 >= 1 {
               v = add(v, countA[w-2][h])
           }
           // f = 2*v - 1
           f := (2*v - 1) % MOD
           if f < 0 {
               f += MOD
           }
           ways := (n - h + 1) * (m - w + 1) % MOD
           ans = (ans + int64(f)*int64(ways)%MOD) % MOD
       }
   }
   fmt.Fprintln(writer, ans)
}
