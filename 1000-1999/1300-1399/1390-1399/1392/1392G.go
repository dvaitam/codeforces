package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

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
   var sa, sb string
   fmt.Fscan(reader, &sa, &sb)
   // initial masks
   maskA, maskB := 0, 0
   for i := 0; i < len(sa); i++ {
       maskA = maskA*2 + int(sa[i]-'0')
   }
   for i := 0; i < len(sb); i++ {
       maskB = maskB*2 + int(sb[i]-'0')
   }
   // read exchanges
   x := make([]int, n)
   y := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x[i], &y[i])
   }
   // prepare arrays
   N := 1 << k
   a := make([]int, N)
   b := make([]int, N)
   infA := n + 1
   // initialize a with infA, b with 0
   for i := 0; i < N; i++ {
       a[i] = infA
       b[i] = 0
   }
   // set initial
   a[maskA] = infA
   b[maskB] = infA
   // positions
   pos := make([]int, k+1)
   for i := 1; i <= k; i++ {
       pos[i] = i
   }
   // process in reverse
   for idx := n - 1; idx >= 0; idx-- {
       xi := x[idx]
       yi := y[idx]
       // swap bits in masks
       // bit positions: pos[p] from 1..k, bit index = k-pos[p]
       px := pos[xi]
       py := pos[yi]
       // extract bits
       bitA_x := (maskA >> (k - px)) & 1
       bitA_y := (maskA >> (k - py)) & 1
       // clear bits
       maskA &^= (1 << (k - px)) | (1 << (k - py))
       // set swapped
       maskA |= (bitA_x << (k - py)) | (bitA_y << (k - px))
       // same for maskB
       bitB_x := (maskB >> (k - px)) & 1
       bitB_y := (maskB >> (k - py)) & 1
       maskB &^= (1 << (k - px)) | (1 << (k - py))
       maskB |= (bitB_x << (k - py)) | (bitB_y << (k - px))
       // swap pos
       pos[xi], pos[yi] = pos[yi], pos[xi]
       // update arrays; use 1-based index
       iVal := idx + 1
       if a[maskA] > iVal {
           a[maskA] = iVal
       }
       if b[maskB] < iVal {
           b[maskB] = iVal
       }
   }
   // SOS DP: propagate min on a, max on b over supersets
   for bit := 0; bit < k; bit++ {
       for mask := 0; mask < N; mask++ {
           if mask&(1<<bit) == 0 {
               nxt := mask | (1 << bit)
               if a[mask] > a[nxt] {
                   a[mask] = a[nxt]
               }
               if b[mask] < b[nxt] {
                   b[mask] = b[nxt]
               }
           }
       }
   }
   // find best
   best := 0
   for mask := 0; mask < N; mask++ {
       if b[mask]-a[mask] >= m {
           pc := bits.OnesCount(uint(mask))
           if pc > best {
               best = pc
           }
       }
   }
   // compute result
   res := best*2 + k - bits.OnesCount(uint(maskA)) - bits.OnesCount(uint(maskB))
   fmt.Fprintln(writer, res)
   // find interval
   for mask := 0; mask < N; mask++ {
       if b[mask]-a[mask] >= m && bits.OnesCount(uint(mask)) == best {
           // start at a[mask], end at b[mask]-1
           fmt.Fprintln(writer, a[mask], b[mask]-1)
           return
       }
   }
}
