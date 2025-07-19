package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, a, b int
   fmt.Fscan(reader, &n, &a, &b)
   f := make([]int, n+1)
   for i := 0; i < a; i++ {
       var x int
       fmt.Fscan(reader, &x)
       f[x] = 1
   }
   for i := 0; i < b; i++ {
       var x int
       fmt.Fscan(reader, &x)
       f[x] = 2
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i := 1; i <= n; i++ {
       fmt.Fprint(writer, f[i])
       if i < n {
           fmt.Fprint(writer, " ")
       }
   }
}
