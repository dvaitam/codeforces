package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, d int
   if _, err := fmt.Fscan(reader, &n, &d); err != nil {
       return
   }
   const maxPos = 30000
   // count gems on each island
   gems := make([]int16, maxPos+1)
   for i := 0; i < n; i++ {
       var p int
       fmt.Fscan(reader, &p)
       if p >= 0 && p <= maxPos {
           gems[p]++
       }
   }
   // dp state: dp[pos][idx], idx = l - d + M
   const M = 300
   const W = 2*M + 1
   total := (maxPos + 1) * W
   dp := make([]int16, total)
   // initialize to -1
   for i := range dp {
       dp[i] = -1
   }
   // for pruning: track min and max idx at each pos
   minIdx := make([]int, maxPos+1)
   maxIdx := make([]int, maxPos+1)
   for i := 0; i <= maxPos; i++ {
       minIdx[i] = W
       maxIdx[i] = -1
   }
   // initial jump to d
   if d <= maxPos {
       idx0 := M
       base := d*W + idx0
       dp[base] = gems[d]
       minIdx[d] = idx0
       maxIdx[d] = idx0
   }
   ans := int16(0)
   // threshold for allowing l-1 (l > 1): idx > M - d + 1
   thr := M - d + 1
   // dp transitions
   for pos := d; pos <= maxPos; pos++ {
       lo := minIdx[pos]
       hi := maxIdx[pos]
       if lo > hi {
           continue
       }
       basePos := pos * W
       for idx := lo; idx <= hi; idx++ {
           cur := dp[basePos+idx]
           if cur < 0 {
               continue
           }
           if cur > ans {
               ans = cur
           }
           // compute jump length l = idx - M + d
           // transitions: l-1, l, l+1 => idx-1, idx, idx+1
           // l-1
           if idx > thr {
               idx2 := idx - 1
               pos2 := pos + (idx2 - M + d)
               if pos2 <= maxPos {
                   off2 := pos2*W + idx2
                   val := cur + gems[pos2]
                   if dp[off2] < val {
                       dp[off2] = val
                       if idx2 < minIdx[pos2] {
                           minIdx[pos2] = idx2
                       }
                       if idx2 > maxIdx[pos2] {
                           maxIdx[pos2] = idx2
                       }
                   }
               }
           }
           // l
           idx2 := idx
           pos2 := pos + (idx2 - M + d)
           if pos2 <= maxPos {
               off2 := pos2*W + idx2
               val := cur + gems[pos2]
               if dp[off2] < val {
                   dp[off2] = val
                   if idx2 < minIdx[pos2] {
                       minIdx[pos2] = idx2
                   }
                   if idx2 > maxIdx[pos2] {
                       maxIdx[pos2] = idx2
                   }
               }
           }
           // l+1
           idx3 := idx + 1
           if idx3 < W {
               pos3 := pos + (idx3 - M + d)
               if pos3 <= maxPos {
                   off3 := pos3*W + idx3
                   val := cur + gems[pos3]
                   if dp[off3] < val {
                       dp[off3] = val
                       if idx3 < minIdx[pos3] {
                           minIdx[pos3] = idx3
                       }
                       if idx3 > maxIdx[pos3] {
                           maxIdx[pos3] = idx3
                       }
                   }
               }
           }
       }
   }
   // final answer
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
}
