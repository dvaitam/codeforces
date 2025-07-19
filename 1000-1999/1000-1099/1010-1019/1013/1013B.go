package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, x int
   if _, err := fmt.Fscan(reader, &n, &x); err != nil {
       return
   }
   const maxV = 100000
   input := make([]bool, maxV+1)
   dan := make([]bool, maxV+1)
   ans := 3
   for i := 0; i < n; i++ {
       var y int
       fmt.Fscan(reader, &y)
       if input[y] {
           ans = min(ans, 0)
       }
       if y <= maxV && dan[y] {
           ans = min(ans, 1)
       }
       t := y & x
       if t <= maxV && input[t] {
           ans = min(ans, 1)
       }
       if t <= maxV && dan[t] {
           ans = min(ans, 2)
       }
       if y <= maxV {
           input[y] = true
       }
       if t <= maxV {
           dan[t] = true
       }
   }
   if ans == 0 {
       fmt.Println(0)
   } else if ans == 1 {
       fmt.Println(1)
   } else if ans == 2 {
       fmt.Println(2)
   } else {
       fmt.Println(-1)
   }
}
