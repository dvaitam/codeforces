package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   isColor := false
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           var pixel byte
           if _, err := fmt.Fscan(reader, &pixel); err != nil {
               return
           }
           switch pixel {
           case 'C', 'M', 'Y':
               isColor = true
           }
       }
   }
   if isColor {
       fmt.Println("#Color")
   } else {
       fmt.Println("#Black&White")
   }
}
