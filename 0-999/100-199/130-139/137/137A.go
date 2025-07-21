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
   var visits int
   var currType rune
   var currCount int
   for _, r := range s {
       if currCount == 0 {
           currType = r
           currCount = 1
           continue
       }
       if r == currType {
           if currCount < 5 {
               currCount++
           } else {
               visits++
               // start new batch of same type
               currCount = 1
           }
       } else {
           // different type: go drop current
           visits++
           currType = r
           currCount = 1
       }
   }
   if currCount > 0 {
       visits++
   }
   fmt.Println(visits)
}
