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
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       cnt := make(map[int]int)
       ans := -1
       for i := 0; i < n; i++ {
           var x int
           fmt.Fscan(reader, &x)
           if ans == -1 {
               cnt[x]++
               if cnt[x] >= 3 {
                   ans = x
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
