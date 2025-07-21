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
   fmt.Fscan(reader, &n)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   cnt := 0
   pos := -1
   for i := 0; i < n-1; i++ {
       if a[i] > a[i+1] {
           cnt++
           pos = i
       }
   }
   if cnt == 0 {
       fmt.Fprintln(writer, 0)
       return
   }
   // More than one drop or suffix cannot be rotated to front
   if cnt > 1 || a[n-1] > a[0] {
       fmt.Fprintln(writer, -1)
       return
   }
   // Number of right rotations needed is length of suffix
   k := n - pos - 1
   fmt.Fprintln(writer, k)
}
