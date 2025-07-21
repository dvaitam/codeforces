package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, c int
   if _, err := fmt.Fscan(in, &n, &c); err != nil {
       return
   }
   prices := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &prices[i])
   }
   maxProfit := 0
   for i := 0; i+1 < n; i++ {
       profit := prices[i] - prices[i+1] - c
       if profit > maxProfit {
           maxProfit = profit
       }
   }
   fmt.Println(maxProfit)
}
