package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   piles := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &piles[i])
   }
   // match if same value or same suit
   match := func(a, b string) bool {
       return a[0] == b[0] || a[1] == b[1]
   }
   // greedy merge any pile i onto i-3 or i-1 when possible
   for {
       moved := false
       for i := range piles {
           // try i onto i-3
           if i >= 3 && match(piles[i], piles[i-3]) {
               // merge i into i-3: new top is piles[i]
               piles[i-3] = piles[i]
               // remove pile i
               piles = append(piles[:i], piles[i+1:]...)
               moved = true
               break
           }
           // try i onto i-1
           if i >= 1 && match(piles[i], piles[i-1]) {
               piles[i-1] = piles[i]
               piles = append(piles[:i], piles[i+1:]...)
               moved = true
               break
           }
       }
       if !moved {
           break
       }
   }
   if len(piles) == 1 {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
