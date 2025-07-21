package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s1, s2 string
   if _, err := fmt.Fscan(reader, &s1); err != nil {
       return
   }
   if _, err := fmt.Fscan(reader, &s2); err != nil {
       return
   }
   switch {
   case s1 > s2:
       fmt.Println("TEAM 1 WINS")
   case s2 > s1:
       fmt.Println("TEAM 2 WINS")
   default:
       fmt.Println("TIE")
   }
}
