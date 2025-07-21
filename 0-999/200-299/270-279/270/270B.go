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
   // Find longest strictly increasing suffix
   pos := n - 1
   for pos > 0 && a[pos-1] < a[pos] {
       pos--
   }
   // Threads in positions [0..pos-1] surely have new messages
   fmt.Println(pos)
}
