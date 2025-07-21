package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var t int64
   fmt.Fscan(reader, &n, &t)
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var sum int64
   ans := 0
   l := 0
   for r := 0; r < n; r++ {
       sum += a[r]
       for sum > t {
           sum -= a[l]
           l++
       }
       if r-l+1 > ans {
           ans = r - l + 1
       }
   }
   fmt.Println(ans)
}
