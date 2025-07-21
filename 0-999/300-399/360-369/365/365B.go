package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   a := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   if n <= 2 {
       fmt.Println(n)
       return
   }
   cur := 2
   ans := 2
   for i := 2; i < n; i++ {
       if a[i] == a[i-1] + a[i-2] {
           cur++
       } else {
           if cur > ans {
               ans = cur
           }
           cur = 2
       }
   }
   if cur > ans {
       ans = cur
   }
   fmt.Println(ans)
}
