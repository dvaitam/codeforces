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
   mid := n / 2
   sum := 0
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           if i == j || i+j == n-1 || i == mid || j == mid {
               sum += x
           }
       }
   }
   fmt.Println(sum)
}
