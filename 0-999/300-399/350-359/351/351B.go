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
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i])
   }
   var k int64
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if a[j] < a[i] {
               k++
           }
       }
   }
   res := (k/2)*4 + k%2
   fmt.Println(res)
}
