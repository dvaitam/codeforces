package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, s int
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       if _, err := fmt.Fscan(reader, &a[i]); err != nil {
           return
       }
   }
   sum := 0
   maxv := 0
   for _, v := range a {
       sum += v
       if v > maxv {
           maxv = v
       }
   }
   if sum-maxv <= s {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
