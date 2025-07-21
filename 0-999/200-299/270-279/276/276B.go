package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   counts := make(map[rune]int)
   for _, ch := range s {
       counts[ch]++
   }
   oddCount := 0
   for _, cnt := range counts {
       if cnt%2 != 0 {
           oddCount++
       }
   }
   if oddCount == 0 || oddCount%2 == 1 {
       fmt.Println("First")
   } else {
       fmt.Println("Second")
   }
}
