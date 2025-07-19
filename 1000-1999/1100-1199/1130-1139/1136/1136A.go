package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i], &b[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   cnt := 0
   for i := 0; i < n; i++ {
       if a[i] >= m || b[i] >= m {
           cnt++
       }
   }
   fmt.Fprintln(writer, cnt)
}
