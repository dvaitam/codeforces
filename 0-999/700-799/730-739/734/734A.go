package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   var s string
   fmt.Fscan(reader, &s)
   aCount, dCount := 0, 0
   for i := 0; i < len(s); i++ {
       switch s[i] {
       case 'A':
           aCount++
       case 'D':
           dCount++
       }
   }
   switch {
   case aCount > dCount:
       fmt.Println("Anton")
   case dCount > aCount:
       fmt.Println("Danik")
   default:
       fmt.Println("Friendship")
   }
}
