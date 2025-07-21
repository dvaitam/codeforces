package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int
   if _, err := fmt.Fscan(reader, &n, &a, &b); err != nil {
       return
   }
   count := 0
   for i := 1; i <= n; i++ {
       if i-1 >= a && n-i <= b {
           count++
       }
   }
   fmt.Println(count)
}
