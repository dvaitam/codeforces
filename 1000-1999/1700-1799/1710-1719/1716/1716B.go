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
       // length of chain is n
       fmt.Fprintln(writer, n)
       // initial permutation
       ans := make([]int, n)
       for i := 0; i < n; i++ {
           ans[i] = i + 1
       }
       // output helper
       printPerm := func(p []int) {
           for i, v := range p {
               if i > 0 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, v)
           }
           writer.WriteByte('\n')
       }
       // first permutation
       printPerm(ans)
       // swap first and last
       ans[0], ans[n-1] = ans[n-1], ans[0]
       printPerm(ans)
       // subsequent adjacent swaps
       for i := 3; i <= n; i++ {
           // swap positions i-2 and i-1 in 1-based => indices i-3, i-2
           j1 := i - 3
           j2 := i - 2
           ans[j1], ans[j2] = ans[j2], ans[j1]
           printPerm(ans)
       }
   }
}
