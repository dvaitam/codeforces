package main

import (
   "bufio"
   "fmt"
   "os"
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, k int
   fmt.Fscan(reader, &n, &m, &k)
   // Read and build prefix sums and interval sums
   M := m - k + 1
   prefix := make([][]int, n+1)
   a := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       row := make([]int, m+1)
       prefix[i] = make([]int, m+1)
       for j := 1; j <= m; j++ {
           fmt.Fscan(reader, &row[j])
           prefix[i][j] = prefix[i][j-1] + row[j]
       }
       a[i] = make([]int, M+1)
       for s := 1; s <= M; s++ {
           a[i][s] = prefix[i][s+k-1] - prefix[i][s-1]
       }
   }
   // dp arrays
   dpPrev := make([]int, M+1)
   // Day 1
   for s := 1; s <= M; s++ {
       dpPrev[s] = a[1][s]
   }
   // Temp arrays
   bPrev := make([]int, M+2)
   preMax := make([]int, M+2)
   sufMax := make([]int, M+2)
   dpCurr := make([]int, M+1)

   // DP for days 2..n
   for day := 2; day <= n; day++ {
       // bPrev[s] = dpPrev[s] + a[day][s]
       for s := 1; s <= M; s++ {
           bPrev[s] = dpPrev[s] + a[day][s]
       }
       // prefix max
       preMax[0] = -1 << 60
       for i := 1; i <= M; i++ {
           if i == 1 {
               preMax[i] = bPrev[i]
           } else {
               preMax[i] = max(preMax[i-1], bPrev[i])
           }
       }
       // suffix max
       sufMax[M+1] = -1 << 60
       for i := M; i >= 1; i-- {
           if i == M {
               sufMax[i] = bPrev[i]
           } else {
               sufMax[i] = max(sufMax[i+1], bPrev[i])
           }
       }
       // compute dpCurr
       for sCurr := 1; sCurr <= M; sCurr++ {
           best := -1 << 60
           // non-overlapping sPrev
           if sCurr-k >= 1 {
               best = max(best, preMax[sCurr-k])
           }
           if sCurr+k <= M {
               best = max(best, sufMax[sCurr+k])
           }
           // overlapping sPrev
           lo := sCurr - k + 1
           if lo < 1 {
               lo = 1
           }
           hi := sCurr + k - 1
           if hi > M {
               hi = M
           }
           for sPrev := lo; sPrev <= hi; sPrev++ {
               // overlap range on day: from max(sPrev, sCurr) to min(sPrev, sCurr)+k-1
               // sum = prefix[day][min(sPrev,sCurr)+k-1] - prefix[day][max(sPrev,sCurr)-1]
               l := sPrev
               r := sCurr
               if l > r {
                   l, r = r, l
               }
               overlapSum := prefix[day][l+k-1] - prefix[day][r-1]
               cand := bPrev[sPrev] - overlapSum
               if cand > best {
                   best = cand
               }
           }
           dpCurr[sCurr] = a[day][sCurr] + best
       }
       // swap dpPrev and dpCurr
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // answer is max over dpPrev
   ans := 0
   for s := 1; s <= M; s++ {
       if dpPrev[s] > ans {
           ans = dpPrev[s]
       }
   }
   fmt.Fprintln(writer, ans)
}
