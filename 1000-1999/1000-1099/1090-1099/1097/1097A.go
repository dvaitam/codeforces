package main

import "fmt"

func main() {
   var table string
   if _, err := fmt.Scan(&table); err != nil {
       return
   }
   yes := false
   for i := 0; i < 5; i++ {
       var card string
       if _, err := fmt.Scan(&card); err != nil {
           return
       }
       if len(card) >= 2 && (card[0] == table[0] || card[1] == table[1]) {
           yes = true
       }
   }
   if yes {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
