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
   curr := 'a'
   total := 0
   for _, ch := range s {
       diff := int(ch - curr)
       if diff < 0 {
           diff = -diff
       }
       if diff > 26-diff {
           diff = 26-diff
       }
       total += diff
       curr = ch
   }
   fmt.Println(total)
}
