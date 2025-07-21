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
   var p float64
   if _, err := fmt.Fscan(reader, &p); err != nil {
       return
   }
   count1 := 0
   countQ := 0
   for _, ch := range s {
       switch ch {
       case '1':
           count1++
       case '?':
           countQ++
       }
   }
   n := float64(len(s))
   expected := (float64(count1) + float64(countQ)*p) / n
   fmt.Printf("%.5f\n", expected)
}
