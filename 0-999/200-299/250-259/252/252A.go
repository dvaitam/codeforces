package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   ans := 0
   for i := 0; i < n; i++ {
       x := 0
       for j := i; j < n; j++ {
           x ^= a[j]
           if x > ans {
               ans = x
           }
       }
   }
   fmt.Println(ans)
}
