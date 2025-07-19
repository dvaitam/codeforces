package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, s int
   fmt.Fscan(reader, &n, &s)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &b[i])
   }
   s-- // convert to 0-based index
   if a[0] == 0 {
       fmt.Println("NO")
       return
   }
   if a[s] == 1 {
       fmt.Println("YES")
       return
   }
   if b[s] == 1 {
       for i := s; i < n; i++ {
           if a[i] == 1 && b[i] == 1 {
               fmt.Println("YES")
               return
           }
       }
   }
   fmt.Println("NO")
}
