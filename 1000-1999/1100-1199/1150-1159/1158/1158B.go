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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   if k == 1 {
       fmt.Fprint(writer, "0")
       for i := 1; i < n; i++ {
           fmt.Fprint(writer, "1")
       }
       fmt.Fprint(writer, "\n")
       return
   }
   d := (n - k) / 2
   l0 := (d + 2) / 2
   l1 := (d + 1) / 2
   res := make([]byte, 0, n)
   for len(res) < n {
       for i := 0; i < l0 && len(res) < n; i++ {
           res = append(res, '0')
       }
       for i := 0; i < l1 && len(res) < n; i++ {
           res = append(res, '1')
       }
   }
   fmt.Fprint(writer, string(res))
   fmt.Fprint(writer, "\n")
}
