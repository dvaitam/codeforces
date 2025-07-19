package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var k int
   fmt.Fscan(reader, &k)
   ans := make([]int, 0, k)
   var y int
   f, f1 := false, false
   for i := 0; i < k; i++ {
       var x int
       fmt.Fscan(reader, &x)
       if x == 0 || x == 100 {
           ans = append(ans, x)
       } else if x < 10 && !f {
           ans = append(ans, x)
           f = true
       } else if x%10 == 0 && !f1 {
           ans = append(ans, x)
           f1 = true
       } else {
           y = x
       }
   }
   if !f && !f1 && y != 0 {
       ans = append(ans, y)
   }
   fmt.Println(len(ans))
   for i, v := range ans {
       if i > 0 {
           fmt.Print(" ")
       }
       fmt.Print(v)
   }
   fmt.Println()
}
