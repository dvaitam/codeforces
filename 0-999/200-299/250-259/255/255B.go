package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       // error reading input
       return
   }
   // remove trailing newline or spaces
   s = strings.TrimSpace(s)
   var countX, countY int
   for i := 0; i < len(s); i++ {
       switch s[i] {
       case 'x':
           countX++
       case 'y':
           countY++
       }
   }
   // Determine the result length
   if countX > countY {
       // more x's remain
       rem := countX - countY
       fmt.Print(strings.Repeat("x", rem))
   } else {
       rem := countY - countX
       fmt.Print(strings.Repeat("y", rem))
   }
}
