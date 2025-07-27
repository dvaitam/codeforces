package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // positions of shields
   P := make([]int, 0, n)
   for i, v := range a {
       if v == 1 {
           P = append(P, i+1)
       }
   }
   m := len(P)
   Z := n - m
   // total zero-zero pairs
   totalPairsZeros := Z * (Z - 1) / 2
   M := n * (n - 1) / 2
   // if fewer than 2 zeros, no protected pairs ever
   if Z < 2 {
       out := make([]string, M+1)
       for i := range out {
           out[i] = "0"
       }
       fmt.Println(strings.Join(out, " "))
       return
   }
   const INF = 1000000000
   maxC := M
   // dpPrev[j][k]: minimal unprotected cost excluding tail, last shield at j with k moves
   dpPrev := make([][]int, n+1)
   dpCurr := make([][]int, n+1)
   for j := 0; j <= n; j++ {
       dpPrev[j] = make([]int, maxC+1)
       dpCurr[j] = make([]int, maxC+1)
       for k := 0; k <= maxC; k++ {
           dpPrev[j][k] = INF
       }
   }
   dpPrev[0][0] = 0
   usedPrev := []int{0}
   prevMaxK := 0
   // DP over shields
   for i := 1; i <= m; i++ {
       // reset dpCurr
       for j := 0; j <= n; j++ {
           for k := 0; k <= maxC; k++ {
               dpCurr[j][k] = INF
           }
       }
       usedCurr := make([]int, 0, n)
       usedMap := make([]bool, n+1)
       currMaxK := 0
       minPos := i
       maxPos := n - (m - i)
       for _, j1 := range usedPrev {
           for j2 := minPos; j2 <= maxPos; j2++ {
               move := P[i-1] - j2
               if move < 0 {
                   move = -move
               }
               d := j2 - j1 - 1
               segCost := d * (d - 1) / 2
               for k1 := 0; k1 <= prevMaxK; k1++ {
                   prevCost := dpPrev[j1][k1]
                   if prevCost >= INF {
                       continue
                   }
                   k2 := k1 + move
                   if k2 > maxC {
                       continue
                   }
                   cost := prevCost + segCost
                   if cost < dpCurr[j2][k2] {
                       dpCurr[j2][k2] = cost
                       if !usedMap[j2] {
                           usedMap[j2] = true
                           usedCurr = append(usedCurr, j2)
                       }
                       if k2 > currMaxK {
                           currMaxK = k2
                       }
                   }
               }
           }
       }
       dpPrev, dpCurr = dpCurr, dpPrev
       usedPrev = usedCurr
       prevMaxK = currMaxK
   }
   // bestCost[k]: minimal total unprotected cost including tail with k moves
   bestCost := make([]int, maxC+1)
   for k := 0; k <= maxC; k++ {
       bestCost[k] = INF
   }
   // add tail segment cost
   for _, j := range usedPrev {
       tailZeros := n - j
       tailCost := tailZeros * (tailZeros - 1) / 2
       for k := 0; k <= prevMaxK; k++ {
           c := dpPrev[j][k]
           if c >= INF {
               continue
           }
           tot := c + tailCost
           if tot < bestCost[k] {
               bestCost[k] = tot
           }
       }
   }
   // prefix min of bestCost
   prefixMin := make([]int, maxC+1)
   curMin := INF
   for k := 0; k <= maxC; k++ {
       if bestCost[k] < curMin {
           curMin = bestCost[k]
       }
       prefixMin[k] = curMin
   }
   // output answers
   out := make([]string, M+1)
   for k := 0; k <= M; k++ {
       if k <= maxC && prefixMin[k] < INF {
           out[k] = strconv.Itoa(totalPairsZeros - prefixMin[k])
       } else {
           out[k] = strconv.Itoa(totalPairsZeros - prefixMin[maxC])
       }
   }
   fmt.Println(strings.Join(out, " "))
