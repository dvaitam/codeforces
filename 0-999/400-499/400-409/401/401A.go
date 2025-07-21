package main

import (
   "bufio"
   "fmt"
   "os"
)

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   fmt.Fscan(reader, &n, &x)
   sum := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       sum += a
   }
   if sum == 0 {
       fmt.Println(0)
       return
   }
   s := abs(sum)
   ans := s / x
   if s%x != 0 {
       ans++
   }
   fmt.Println(ans)
}
