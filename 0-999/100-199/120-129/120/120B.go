package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // adjust to 0-based index
   pos := k - 1
   // search from pos to end
   for i := pos; i < n; i++ {
       if a[i] == 1 {
           fmt.Println(i + 1)
           return
       }
   }
   // wrap around
   for i := 0; i < pos; i++ {
       if a[i] == 1 {
           fmt.Println(i + 1)
           return
       }
   }
}
