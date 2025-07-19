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
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       var a, b int64
       for i := 0; i < n; i++ {
           var x int64
           fmt.Fscan(reader, &x)
           a += x
       }
       for i := 0; i < n; i++ {
           var x int64
           fmt.Fscan(reader, &x)
           b += x
       }
       if a >= b {
           fmt.Fprintln(writer, "Yes")
       } else {
           fmt.Fprintln(writer, "No")
       }
   }
}
