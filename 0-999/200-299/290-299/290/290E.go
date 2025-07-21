package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   s, err := reader.ReadString('\n')
   if err != nil && len(s) == 0 {
       return
   }
   for i := 0; i < len(s); i++ {
       switch s[i] {
       case 'H', 'Q', '9':
           fmt.Println("Yes")
           return
       }
   }
   fmt.Println("No")
}
