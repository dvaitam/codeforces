package main

import "fmt"

func main() {
   var n int
   if _, err := fmt.Scan(&n); err != nil {
       return
   }
   cards := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Scan(&cards[i])
   }
   left, right := 0, n-1
   var serejaSum, dimaSum int
   turn := 0 // 0 for Sereja, 1 for Dima
   for left <= right {
       var pick int
       if cards[left] >= cards[right] {
           pick = cards[left]
           left++
       } else {
           pick = cards[right]
           right--
       }
       if turn == 0 {
           serejaSum += pick
       } else {
           dimaSum += pick
       }
       turn ^= 1
   }
   fmt.Println(serejaSum, dimaSum)
}
