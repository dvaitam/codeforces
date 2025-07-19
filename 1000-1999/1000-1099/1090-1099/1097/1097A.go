package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var table string
   if _, err := fmt.Fscan(reader, &table); err != nil {
       return
   }
   yes := false
   for i := 0; i < 5; i++ {
       var card string
       if _, err := fmt.Fscan(reader, &card); err != nil {
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
