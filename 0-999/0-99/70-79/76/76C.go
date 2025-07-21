package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, K int
   var Tlimit int64
   fmt.Fscan(in, &N, &K, &Tlimit)
   var s string
   fmt.Fscan(in, &s)
   tAll := make([]int64, K)
   for i := 0; i < K; i++ {
       fmt.Fscan(in, &tAll[i])
   }
   aAll := make([][]int64, K)
   for i := 0; i < K; i++ {
       aAll[i] = make([]int64, K)
       for j := 0; j < K; j++ {
           fmt.Fscan(in, &aAll[i][j])
       }
   }
   // map present types
   present := make([]bool, K)
   for i := 0; i < N; i++ {
       present[s[i]-'A'] = true
   }
   typeMap := make([]int, K)
   revMap := make([]int, 0, K)
   for i := 0; i < K; i++ {
       if present[i] {
           typeMap[i] = len(revMap)
           revMap = append(revMap, i)
       }
   }
   m := len(revMap)
   if m == 0 {
       fmt.Println(0)
       return
   }
   // map string to indices in [0,m)
   A := make([]int, N)
   for i := 0; i < N; i++ {
       A[i] = typeMap[s[i]-'A']
   }
   // build nextpos
   nextpos := make([][]int, m)
   for x := 0; x < m; x++ {
       nextpos[x] = make([]int, N+1)
       next := N
       for i := N - 1; i >= 0; i-- {
           if A[i] == x {
               next = i
           }
           nextpos[x][i] = next
       }
       nextpos[x][N] = N
   }
   // accumulate total_out and updates per j
   size := 1 << m
   totalOut := make([]int64, size)
   updates := make([][]struct{mask uint32; val int64}, m)
   type pair struct{pos, idx int}
   ord := make([]pair, 0, m)
   // iterate positions
   for p := 0; p < N; p++ {
       i := A[p]
       ord = ord[:0]
       for x := 0; x < m; x++ {
           np := nextpos[x][p]
           if np < N {
               ord = append(ord, pair{np, x})
           }
       }
       sort.Slice(ord, func(i, j int) bool { return ord[i].pos < ord[j].pos })
       var mask uint32
       for _, pr := range ord {
           j := pr.idx
           // cost a[i][j]
           cost := aAll[ revMap[i] ][ revMap[j] ]
           totalOut[mask] += cost
           updates[j] = append(updates[j], struct{mask uint32; val int64}{mask, cost})
           mask |= 1 << j
       }
   }
   // f1: zeta transform on totalOut
   for b := 0; b < m; b++ {
       for mask := 0; mask < size; mask++ {
           if mask&(1<<b) != 0 {
               totalOut[mask] += totalOut[mask^(1<<b)]
           }
       }
   }
   // h[M] = sum_{j in M} sum_{S subset M} W[S][j]
   h := make([]int64, size)
   // process each j
   for j := 0; j < m; j++ {
       W := make([]int64, size)
       for _, u := range updates[j] {
           W[u.mask] += u.val
       }
       updates[j] = nil
       // zeta on W
       for b := 0; b < m; b++ {
           for mask := 0; mask < size; mask++ {
               if mask&(1<<b) != 0 {
                   W[mask] += W[mask^(1<<b)]
               }
           }
       }
       // add to h where bit j set
       bit := 1 << j
       for mask := bit; mask < size; mask++ {
           if mask&bit != 0 {
               h[mask] += W[mask]
           }
       }
   }
   // precompute tMask
   tMap := make([]int64, m)
   for idx, orig := range revMap {
       tMap[idx] = tAll[orig]
   }
   tMask := make([]int64, size)
   for mask := 1; mask < size; mask++ {
       lsb := mask & -mask
       b := bitsTrailing(uint(lsb))
       tMask[mask] = tMask[mask^int(lsb)] + tMap[b]
   }
   // count
   full := size - 1
   var ans int64
   for mask := 0; mask < size; mask++ {
       if mask == full {
           continue
       }
       riskAdj := totalOut[mask] - h[mask]
       if riskAdj + tMask[mask] <= Tlimit {
           ans++
       }
   }
   fmt.Println(ans)
}

// bitsTrailing returns index of trailing bit (0-based)
func bitsTrailing(x uint) int {
   return bitsTrailingDeBruijn(x)
}
var deBruijn = uint(0x077CB531)
var idx32 = [32]int{
   0, 1, 28, 2, 29, 14, 24, 3,
   30, 22, 20, 15, 25, 17, 4, 8,
   31, 27, 13, 23, 21, 19, 16, 7,
   26, 12, 18, 6, 11, 5, 10, 9,
}
func bitsTrailingDeBruijn(v uint) int {
   return idx32[((v&-v)*deBruijn)>>27]
}
