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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       flg := false
       for i := 2; i < n; i++ {
           if a[0]+a[1] <= a[i] {
               fmt.Fprintf(writer, "%d %d %d\n", 1, 2, i+1)
               flg = true
               break
           }
       }
       if !flg {
           fmt.Fprintln(writer, -1)
       }
   }
}
