package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   _, err := fmt.Fscan(reader, &t)
   if err != nil {
       return
   }
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       // reverse the string
       runes := []rune(s)
       for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
           runes[l], runes[r] = runes[r], runes[l]
       }
       rev := string(runes)
       fmt.Fprintln(writer, rev+s)
   }
}
