package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   // Read the username string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // Track distinct lowercase letters
   seen := make([]bool, 26)
   count := 0
   for _, ch := range s {
       idx := ch - 'a'
       if idx >= 0 && idx < 26 {
           if !seen[idx] {
               seen[idx] = true
               count++
           }
       }
   }
   // Determine gender by parity of distinct count
   if count%2 == 1 {
       fmt.Println("IGNORE HIM!")
   } else {
       fmt.Println("CHAT WITH HER!")
   }
}
