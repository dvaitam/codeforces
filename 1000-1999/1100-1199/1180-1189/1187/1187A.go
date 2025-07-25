package main

import (
   "fmt"
)

func main() {
   var T int
   if _, err := fmt.Scan(&T); err != nil {
       return
   }
   for i := 0; i < T; i++ {
       var n, s, t int64
       fmt.Scan(&n, &s, &t)
       // Number of eggs with only sticker = n - t
       // Number of eggs with only toy = n - s
       // Need max(onlySticker, onlyToy) + 1
       onlySticker := n - t
       onlyToy := n - s
       ans := onlySticker
       if onlyToy > ans {
           ans = onlyToy
       }
       fmt.Println(ans + 1)
   }
}
