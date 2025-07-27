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
   for i := 0; i < t; i++ {
       var x, y, k int64
       fmt.Fscan(reader, &x, &y, &k)
       // Total sticks needed: k for torches + y*k for coals
       // We start with 1 stick, and each stick trade gives net (x-1) sticks
       // Let neededSticks = k*(y+1) - 1 (subtract initial stick)
       needed := k*(y+1) - 1
       denom := x - 1
       // Number of stick trades to cover needed sticks
       stickTrades := (needed + denom - 1) / denom
       // You also need k coal trades
       result := stickTrades + k
       fmt.Fprintln(writer, result)
   }
}
