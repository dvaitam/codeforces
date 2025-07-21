package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(reader, &n, &m)
   // read horizontal capacities: m-1 lines of n ints
   h := make([][]int64, m-1)
   for j := 0; j < m-1; j++ {
       row := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &row[i])
       }
       h[j] = row
   }
   // read vertical capacities: m lines of n ints
   v := make([][]int64, m)
   for j := 0; j < m; j++ {
       row := make([]int64, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &row[i])
       }
       v[j] = row
   }
   // precompute bit difference positions for horizontal and vertical
   maxMask := 1 << n
   posDiff := make([][]int, maxMask)
   for mask := 0; mask < maxMask; mask++ {
       // horizontal diff bits: positions where two masks differ
       // for XOR values, reused for horizontal
       var bits []int
       for i := 0; i < n; i++ {
           if mask&(1<<i) != 0 {
               bits = append(bits, i)
           }
       }
       posDiff[mask] = bits
   }
   // vertical diff positions: for each mask, positions where mask bit != next row bit
   posVert := make([][]int, maxMask)
   for mask := 0; mask < maxMask; mask++ {
       var bits []int
       for i := 0; i < n; i++ {
           ni := (i + 1) % n
           if ((mask>>i)&1) != ((mask>>ni)&1) {
               bits = append(bits, i)
           }
       }
       posVert[mask] = bits
   }
   // DP arrays
   INF := int64(9e18)
   dpPrev := make([]int64, maxMask)
   dpCur := make([]int64, maxMask)
   // initial column 1: mask must be zero
   for mask := 0; mask < maxMask; mask++ {
       dpPrev[mask] = INF
   }
   dpPrev[0] = 0
   // include vertical cuts in column 1 if any (mask=0 means no vertical cuts)
   // iterate columns 2..m
   for j := 2; j <= m; j++ {
       // compute costHoriz for boundary j-1
       costH := make([]int64, maxMask)
       rowH := h[j-2]
       for xor := 0; xor < maxMask; xor++ {
           sum := int64(0)
           for _, i := range posDiff[xor] {
               sum += rowH[i]
           }
           costH[xor] = sum
       }
       // compute costVert for column j
       costV := make([]int64, maxMask)
       rowV := v[j-1]
       for mask := 0; mask < maxMask; mask++ {
           sum := int64(0)
           for _, i := range posVert[mask] {
               sum += rowV[i]
           }
           costV[mask] = sum
       }
       // DP transition
       for mask2 := 0; mask2 < maxMask; mask2++ {
           dpCur[mask2] = INF
       }
       for mask1 := 0; mask1 < maxMask; mask1++ {
           base := dpPrev[mask1]
           if base >= INF {
               continue
           }
           // try all mask2
           for mask2 := 0; mask2 < maxMask; mask2++ {
               cost := base + costH[mask1^mask2] + costV[mask2]
               if cost < dpCur[mask2] {
                   dpCur[mask2] = cost
               }
           }
       }
       // swap
       dpPrev, dpCur = dpCur, dpPrev
   }
   // result is dpPrev[all ones]
   full := (1 << n) - 1
   res := dpPrev[full]
   fmt.Println(res)
}
