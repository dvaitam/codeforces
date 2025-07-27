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
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       // frequency and first index for each value
       freq := make([]int, n+2)
       idx := make([]int, n+2)
       for i := 1; i <= n; i++ {
           var a int
           fmt.Fscan(reader, &a)
           if freq[a] == 0 {
               freq[a] = 1
               idx[a] = i
           } else if freq[a] == 1 {
               freq[a] = 2
           }
       }
       ans := -1
       for v := 1; v <= n; v++ {
           if freq[v] == 1 {
               ans = idx[v]
               break
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
