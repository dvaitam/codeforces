package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   var l, x, y int
   fmt.Fscan(in, &n, &l, &x, &y)
   a := make([]int, n)
   for i := 0; i < n; i++ {
      fmt.Fscan(in, &a[i])
   }

   judge := func(d int) bool {
      r := 0
      for i := 0; i < n; i++ {
         for r < n && a[r]-a[i] < d {
            r++
         }
         if r < n && a[r]-a[i] == d {
            return true
         }
      }
      return false
   }

   f1 := judge(x)
   f2 := judge(y)
   if f1 && f2 {
      fmt.Fprintln(out, 0)
      return
   }
   if f1 || f2 {
      fmt.Fprintln(out, 1)
      if f1 {
         fmt.Fprintln(out, y)
      } else {
         fmt.Fprintln(out, x)
      }
      return
   }

   ck := func(pos int) bool {
      if pos < 0 || pos > l {
         return false
      }
      idx := sort.Search(n, func(i int) bool { return a[i] >= pos+y })
      if idx < n && a[idx] == pos+y {
         return true
      }
      idx = sort.Search(n, func(i int) bool { return a[i] >= pos-y })
      if idx < n && a[idx] == pos-y {
         return true
      }
      return false
   }
   for i := 0; i < n; i++ {
      p := a[i] + x
      if ck(p) {
         fmt.Fprintln(out, 1)
         fmt.Fprintln(out, p)
         return
      }
      p = a[i] - x
      if ck(p) {
         fmt.Fprintln(out, 1)
         fmt.Fprintln(out, p)
         return
      }
   }
   fmt.Fprintln(out, 2)
   fmt.Fprintln(out, x, y)
}
