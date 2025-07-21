package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   var s string
   if _, err := fmt.Fscan(reader, &n, &s); err != nil {
       return
   }
   count := 0
   // s is 0-based, turn j corresponds to index i = j-1
   for i := 3; i < len(s); i++ {
       // check if this turn is Vasya's: (j-1) mod n == 0 => i mod n == 0
       if i % n != 0 {
           continue
       }
       // check last three moves are identical
       if s[i-1] == s[i-2] && s[i-2] == s[i-3] {
           count++
       }
   }
   fmt.Println(count)
}
