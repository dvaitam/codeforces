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

   var t int
   fmt.Fscan(reader, &t)
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       b := make([]struct{v, idx int}, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i].v)
           b[i].idx = i
       }
       // b is already non-decreasing
       // binary search on d
       lo, hi, best := 0, 1000001, 0
       for lo <= hi {
           mid := (lo + hi) >> 1
           if feasible(a, b, mid, nil) {
               best = mid
               lo = mid + 1
           } else {
               hi = mid - 1
           }
       }
       // build answer
       p := make([]int, n)
       ok := feasible(a, b, best, p)
       _ = ok // should be true
       // output p (convert to 1-based)
       for i := 0; i < n; i++ {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, p[i]+1)
       }
       writer.WriteByte('\n')
   }
}

// feasible checks if minimal difference d can be achieved and optionally fills p (len n)
func feasible(a []int, b []struct{v, idx int}, d int, p []int) bool {
   n := len(a)
   l, r := 0, n-1
   for i := 0; i < n; i++ {
       av := a[i]
       frontOK := l <= r && b[l].v <= av-d
       backOK := l <= r && b[r].v >= av+d
       if frontOK && backOK {
           // choose side with larger diff, tie favor back
           diffL := av - b[l].v
           diffR := b[r].v - av
           if diffR >= diffL {
               if p != nil {
                   p[i] = b[r].idx
               }
               r--
           } else {
               if p != nil {
                   p[i] = b[l].idx
               }
               l++
           }
       } else if frontOK {
           if p != nil {
               p[i] = b[l].idx
           }
           l++
       } else if backOK {
           if p != nil {
               p[i] = b[r].idx
           }
           r--
       } else {
           return false
       }
   }
   return true
}
