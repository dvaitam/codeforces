package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   rdr := bufio.NewReader(os.Stdin)
   var n int
   var l, t int64
   if _, err := fmt.Fscan(rdr, &n, &l, &t); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(rdr, &a[i])
   }
   // double time
   t *= 2
   full := t / l
   left := t % l
   var res int64
   j := 0
   phase := false
   for i := 0; i < n; i++ {
       if left + a[i] >= l {
           phase = true
           j = 0
           left -= l
       }
       for j < n && a[j] <= left + a[i] {
           j++
       }
       if !phase {
           res += (full+1)*int64(j-i-1) + full*int64(n-j+i)
       } else {
           res += (full+1)*int64(n+j-i-1) + full*int64(i-j)
       }
   }
   // result divided by 4.0 with 6 decimals
   ans := float64(res) / 4.0
   fmt.Printf("%.6f\n", ans)
}
