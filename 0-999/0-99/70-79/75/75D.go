package main

import (
   "bufio"
   "fmt"
   "os"
)

const INF = 1 << 60

// block holds precomputed sums for a small array
type block struct {
   sum, pref, suf, best int64
}

func max(a, b int64) int64 {
   if a > b {
      return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   info := make([]block, n)
   for i := 0; i < n; i++ {
      var l int
      fmt.Fscan(reader, &l)
      arr := make([]int64, l)
      for j := 0; j < l; j++ {
         fmt.Fscan(reader, &arr[j])
      }
      // total sum
      var s int64
      for _, v := range arr {
         s += v
      }
      // max prefix sum
      curr := int64(0)
      p := arr[0]
      for j := 0; j < l; j++ {
         curr += arr[j]
         if curr > p {
            p = curr
         }
      }
      // max suffix sum
      curr = 0
      q := arr[l-1]
      for j := l - 1; j >= 0; j-- {
         curr += arr[j]
         if curr > q {
            q = curr
         }
      }
      // max subarray sum (Kadane)
      best := arr[0]
      curr = arr[0]
      for j := 1; j < l; j++ {
         if curr > 0 {
            curr += arr[j]
         } else {
            curr = arr[j]
         }
         if curr > best {
            best = curr
         }
      }
      info[i] = block{sum: s, pref: p, suf: q, best: best}
   }
   ans := -INF
   var curSuffix int64
   for k := 0; k < m; k++ {
      var idx int
      fmt.Fscan(reader, &idx)
      b := info[idx-1]
      if b.best > ans {
         ans = b.best
      }
      if curSuffix+b.pref > ans {
         ans = curSuffix + b.pref
      }
      if curSuffix+b.sum > b.suf {
         curSuffix = curSuffix + b.sum
      } else {
         curSuffix = b.suf
      }
   }
   fmt.Fprintln(writer, ans)
