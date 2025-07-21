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

   var v int
   if _, err := fmt.Fscan(reader, &v); err != nil {
       return
   }
   const maxD = 9
   cost := make([]int, maxD+1)
   for d := 1; d <= maxD; d++ {
       fmt.Fscan(reader, &cost[d])
   }
   // find minimal cost
   minCost := cost[1]
   for d := 2; d <= maxD; d++ {
       if cost[d] < minCost {
           minCost = cost[d]
       }
   }
   if v < minCost {
       fmt.Fprintln(writer, -1)
       return
   }
   // maximum length
   length := v / minCost
   rem := v - length*minCost
   result := make([]byte, length)
   // build number
   for i := 0; i < length; i++ {
       // try to put the largest digit possible
       for d := maxD; d >= 1; d-- {
           extra := cost[d] - minCost
           if extra <= rem {
               result[i] = byte('0' + d)
               rem -= extra
               break
           }
       }
   }
   writer.Write(result)
}
