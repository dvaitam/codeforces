package main

import (
   "bufio"
   "fmt"
   "os"
)

var a, b, ta, tb []int

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func solve(as, ae, bs, be int) int {
   cnt0, cnt1 := 0, 0
   for i := as; i <= ae; i++ {
       if a[i]&1 == 0 {
           cnt0++
       } else {
           cnt1++
       }
   }
   for i := bs; i <= be; i++ {
       if (b[i]&1)^1 == 0 {
           cnt0++
       } else {
           cnt1++
       }
   }
   ret := max(cnt0, cnt1)
   if as > ae || bs > be {
       return ret
   }
   if as == ae && bs == be {
       return max(ret, 2)
   }
   j := as
   for i := as; i <= ae; i++ {
       if a[i]&1 == 1 {
           ta[j] = a[i] >> 1
           j++
       }
   }
   A := j
   for i := as; i <= ae; i++ {
       if a[i]&1 == 0 {
           ta[j] = a[i] >> 1
           j++
       }
   }
   for i := as; i <= ae; i++ {
       a[i] = ta[i]
   }
   k := bs
   for i := bs; i <= be; i++ {
       if b[i]&1 == 1 {
           tb[k] = b[i] >> 1
           k++
       }
   }
   B := k
   for i := bs; i <= be; i++ {
       if b[i]&1 == 0 {
           tb[k] = b[i] >> 1
           k++
       }
   }
   for i := bs; i <= be; i++ {
       b[i] = tb[i]
   }
   ret = max(ret, solve(as, A-1, bs, B-1))
   ret = max(ret, solve(A, ae, B, be))
   return ret
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, y1, y2 int
   fmt.Fscan(reader, &n, &y1)
   a = make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   fmt.Fscan(reader, &m, &y2)
   b = make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &b[i])
   }
   ta = make([]int, n)
   tb = make([]int, m)
   ans := solve(0, n-1, 0, m-1)
   fmt.Println(ans)
}
