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
   for t > 0 {
       t--

       var n int
       fmt.Fscan(reader, &n)
       ans := make([]int, n)
       // find smallest x such that x*x > n
       x := 0
       for x*x <= n {
           x++
       }
       // fill ans from the end
       i := n - 1
       cur := n
       for i >= 0 {
           nx := cur
           for i >= 0 && x*x - i < cur {
               ans[i] = x*x - i
               nx--
               i--
           }
           cur = nx
           x--
       }
       // output
       for j := 0; j < n; j++ {
           if j > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, ans[j])
       }
       fmt.Fprint(writer, '\n')
   }
}
