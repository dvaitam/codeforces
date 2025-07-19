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

   var k int64
   if _, err := fmt.Fscan(reader, &k); err != nil {
       return
   }
   const n = 2000
   ans := make([]int64, n)
   ans[0] = -1
   ans[1] = n
   lft := k
   for i := 2; i < n; i++ {
       x := lft
       if x > 1000000 {
           x = 1000000
       }
       ans[i] = x
       lft -= x
   }

   fmt.Fprintln(writer, n)
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
