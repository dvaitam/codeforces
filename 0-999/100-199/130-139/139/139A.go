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
   days := make([]int, 7)
   for i := 0; i < 7; i++ {
       fmt.Fscan(reader, &days[i])
   }
   remaining := n
   for {
       for i, pages := range days {
           if remaining <= pages {
               // days are 1-indexed: Monday=1, ..., Sunday=7
               fmt.Println(i + 1)
               return
           }
           remaining -= pages
       }
   }
}
