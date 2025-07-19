package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var k int64
   fmt.Fscan(reader, &n, &k)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   L := int64(-4e18) - 100
   R := int64(1e9) + 100
   b := make([]int64, n)
   exact := make([]int, 0, n)

   // get_b computes max m such that 3*(m*m + m) + 1 <= border
   getB := func(border int64) int64 {
       if border <= 0 {
           return -1
       }
       bb := math.Sqrt(float64(border-1) / 3.0)
       lo := int64(bb) - 5
       if lo < 0 {
           lo = 0
       }
       hi := int64(bb) + 5
       for lo+1 < hi {
           m := (lo + hi) >> 1
           if 3*(m*m + m) + 1 <= border {
               lo = m
           } else {
               hi = m
           }
       }
       return lo
   }

   var fillB func(m, maxSum int64)
   fillB = func(m, maxSum int64) {
       exact = exact[:0]
       var s int64
       for i := 0; i < n; i++ {
           border := a[i] - m
           cur := getB(border)
           if cur > a[i]-1 {
               cur = a[i] - 1
           }
           if 3*(cur*cur + cur) + 1 == border {
               exact = append(exact, i)
           }
           b[i] = cur + 1
           s += b[i]
       }
       // reduce if sum exceeds maxSum
       for s > maxSum {
           idx := exact[len(exact)-1]
           exact = exact[:len(exact)-1]
           b[idx]--
           s--
       }
   }

   // binary search for threshold
   const INF = int64(1e18)
   for L+1 < R {
       m := (L + R) >> 1
       fillB(m, INF)
       var sumB int64
       for i := 0; i < n; i++ {
           sumB += b[i]
       }
       if sumB < k {
           R = m
       } else {
           L = m
       }
   }
   fillB(L, k)
   // output result
   for i := 0; i < n; i++ {
       if i > 0 {
           fmt.Fprint(writer, " ")
       }
       fmt.Fprint(writer, b[i])
   }
   fmt.Fprint(writer, "\n")
}
