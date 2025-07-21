package main

import (
   "fmt"
)

func main() {
   var a, b int
   if _, err := fmt.Scan(&a, &b); err != nil {
       return
   }
   candles := a
   leftovers := 0
   hours := 0
   for candles > 0 {
       hours += candles
       leftovers += candles
       candles = leftovers / b
       leftovers %= b
   }
   fmt.Println(hours)
}
