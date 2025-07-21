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

   // Read available quantities of reagents
   var stocks []int
   for {
       var x int
       if _, err := fmt.Fscan(reader, &x); err != nil {
           break
       }
       stocks = append(stocks, x)
   }
   // Required quantities per batch: Aqua Fortis, Aqua Regia, Amalgama, Minium, Vitriol
   req := []int{1, 1, 2, 7, 4}
   // Compute maximum number of full batches
   ans := 0
   for i, need := range req {
       if i >= len(stocks) {
           ans = 0
           break
       }
       batches := stocks[i] / need
       if i == 0 || batches < ans {
           ans = batches
       }
   }
   fmt.Fprintln(writer, ans)
}
