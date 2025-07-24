package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   var ans int64
   for j := 0; j < m; j++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       var sum int64
       for i := l - 1; i < r; i++ {
           sum += a[i]
       }
       if sum > 0 {
           ans += sum
       }
   }
   fmt.Println(ans)
}
