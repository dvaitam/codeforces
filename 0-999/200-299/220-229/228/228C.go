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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }
   // prefix sums of '*' counts
   sum := make([][]int, n+1)
   for i := range sum {
       sum[i] = make([]int, m+1)
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           sum[i+1][j+1] = sum[i+1][j] + sum[i][j+1] - sum[i][j]
           if grid[i][j] == '*' {
               sum[i+1][j+1]++
           }
       }
   }
   // helper to check if block all '*'
   allStars := func(r, c, h int) bool {
       total := sum[r+h][c+h] - sum[r][c+h] - sum[r+h][c] + sum[r][c]
       return total == h*h
   }
   // max level such that 2^level <= min(n,m)
   maxLevel := 0
   for (1<<maxLevel) <= n && (1<<maxLevel) <= m {
       maxLevel++
   }
   maxLevel-- // now 2^maxLevel <= n,m
   // we need levels >=2
   // counted[level][i][j]
   counted := make([][][]bool, maxLevel+1)
   for k := 0; k <= maxLevel; k++ {
       counted[k] = make([][]bool, n)
       for i := range counted[k] {
           counted[k][i] = make([]bool, m)
       }
   }
   ans := 0
   // DP arrays
   dpPrev := make([][]bool, n)
   dpCurr := make([][]bool, n)
   for i := 0; i < n; i++ {
       dpPrev[i] = make([]bool, m)
       dpCurr[i] = make([]bool, m)
   }
   // iterate over all patterns P (2x2)
   for mask := 0; mask < 16; mask++ {
       // extract pattern bits P[a][b]
       var P [2][2]bool
       P[0][0] = mask&1 != 0
       P[0][1] = mask&2 != 0
       P[1][0] = mask&4 != 0
       P[1][1] = mask&8 != 0
       // level 1: size=2, match exact pattern
       size := 2
       for i := 0; i+size <= n; i++ {
           for j := 0; j+size <= m; j++ {
               ok := true
               for a := 0; a < 2; a++ {
                   for b := 0; b < 2; b++ {
                       want := P[a][b]
                       if (grid[i+a][j+b] == '*') != want {
                           ok = false
                           break
                       }
                   }
                   if !ok {
                       break
                   }
               }
               dpPrev[i][j] = ok
           }
       }
       // higher levels
       for level := 2; level <= maxLevel; level++ {
           size = 1 << level
           half := size >> 1
           for i := 0; i+size <= n; i++ {
               for j := 0; j+size <= m; j++ {
                   ok := true
                   for a := 0; a < 2; a++ {
                       for b := 0; b < 2; b++ {
                           si := i + a*half
                           sj := j + b*half
                           if P[a][b] {
                               if !allStars(si, sj, half) {
                                   ok = false
                               }
                           } else {
                               if !dpPrev[si][sj] {
                                   ok = false
                               }
                           }
                           if !ok {
                               break
                           }
                       }
                       if !ok {
                           break
                       }
                   }
                   dpCurr[i][j] = ok
                   if ok && !counted[level][i][j] {
                       counted[level][i][j] = true
                       ans++
                   }
               }
           }
           // swap dpPrev, dpCurr and clear dpCurr
           for i := 0; i+size <= n; i++ {
               for j := 0; j+size <= m; j++ {
                   dpPrev[i][j] = dpCurr[i][j]
                   dpCurr[i][j] = false
               }
           }
       }
   }
   fmt.Fprintln(writer, ans)
}
