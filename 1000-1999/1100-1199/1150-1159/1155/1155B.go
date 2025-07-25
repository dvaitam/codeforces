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
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // number of moves each player makes
   moves := (n - 11) / 2
   // count '8's in the prefix of length n-10 (indices 0 to n-11)
   limit := n - 10
   count8 := 0
   for i := 0; i < limit && i < len(s); i++ {
       if s[i] == '8' {
           count8++
       }
   }
   if count8 > moves {
       fmt.Println("YES")
   } else {
       fmt.Println("NO")
   }
}
