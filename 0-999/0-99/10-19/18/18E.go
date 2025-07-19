package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1000000000
const ALPHA = 26

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   str := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &str[i])
   }
   // cost[row][x][y]: cost for assigning x on evens, y on odds
   cost := make([][ALPHA][ALPHA]int, n)
   for i := 0; i < n; i++ {
       var even [ALPHA]int
       var odd [ALPHA]int
       row := str[i]
       for j := 0; j < m; j++ {
           c := int(row[j] - 'a')
           if j%2 == 0 {
               // even position: cost for any x != c
               for x := 0; x < ALPHA; x++ {
                   if x != c {
                       even[x]++
                   }
               }
           } else {
               for x := 0; x < ALPHA; x++ {
                   if x != c {
                       odd[x]++
                   }
               }
           }
       }
       for x := 0; x < ALPHA; x++ {
           for y := 0; y < ALPHA; y++ {
               cost[i][x][y] = even[x] + odd[y]
           }
       }
   }
   // dpPrev and dpCurr: dp for previous and current row
   var dpPrev [ALPHA][ALPHA]int
   var dpCurr [ALPHA][ALPHA]int
   // pre[row][x][y] = previous (px,py)
   pre := make([][ALPHA][ALPHA][2]int, n+1)
   // init dpPrev for row 0
   for x := 0; x < ALPHA; x++ {
       for y := 0; y < ALPHA; y++ {
           if x != y {
               dpPrev[x][y] = 0
           } else {
               dpPrev[x][y] = INF
           }
       }
   }
   // DP
   for i := 0; i < n; i++ {
       // find best a,b from dpPrev
       best := INF
       ba, bb := 0, 0
       for x := 0; x < ALPHA; x++ {
           for y := 0; y < ALPHA; y++ {
               if dpPrev[x][y] < best {
                   best = dpPrev[x][y]
                   ba, bb = x, y
               }
           }
       }
       // compute dpCurr
       for x := 0; x < ALPHA; x++ {
           for y := 0; y < ALPHA; y++ {
               dpCurr[x][y] = INF
               if x == y {
                   continue
               }
               // if no conflict with best
               if x != ba && y != bb {
                   dpCurr[x][y] = best + cost[i][x][y]
                   pre[i+1][x][y][0] = ba
                   pre[i+1][x][y][1] = bb
               } else {
                   // search other
                   curBest := INF
                   pa, pb := 0, 0
                   for u := 0; u < ALPHA; u++ {
                       if u == x {
                           continue
                       }
                       for v := 0; v < ALPHA; v++ {
                           if v == y {
                               continue
                           }
                           if dpPrev[u][v] < curBest {
                               curBest = dpPrev[u][v]
                               pa, pb = u, v
                           }
                       }
                   }
                   dpCurr[x][y] = curBest + cost[i][x][y]
                   pre[i+1][x][y][0] = pa
                   pre[i+1][x][y][1] = pb
               }
           }
       }
       // swap dpPrev and dpCurr
       for x := 0; x < ALPHA; x++ {
           for y := 0; y < ALPHA; y++ {
               dpPrev[x][y] = dpCurr[x][y]
           }
       }
   }
   // find best at row n
   ansCost := INF
   a, b := 0, 1
   for x := 0; x < ALPHA; x++ {
       for y := 0; y < ALPHA; y++ {
           if dpPrev[x][y] < ansCost {
               ansCost = dpPrev[x][y]
               a, b = x, y
           }
       }
   }
   fmt.Fprintln(out, ansCost)
   // reconstruct
   ans := make([][]byte, n)
   for i := n - 1; i >= 0; i-- {
       row := make([]byte, m)
       for j := 0; j < m; j++ {
           if j%2 == 0 {
               row[j] = byte('a' + a)
           } else {
               row[j] = byte('a' + b)
           }
       }
       ans[i] = row
       pa := pre[i+1][a][b][0]
       pb := pre[i+1][a][b][1]
       a, b = pa, pb
   }
   for i := 0; i < n; i++ {
       out.Write(ans[i])
       out.WriteByte('\n')
   }
}
