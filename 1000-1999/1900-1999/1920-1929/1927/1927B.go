package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       counts := make(map[rune]int)
       letters := make([]rune, 0, n)
       var sb strings.Builder
       next := 'a'
       for i := 0; i < n; i++ {
           var a int
           fmt.Fscan(reader, &a)
           if a == 0 {
               counts[next]++
               letters = append(letters, next)
               sb.WriteRune(next)
               next++
           } else {
               for _, r := range letters {
                   if counts[r] == a {
                       sb.WriteRune(r)
                       counts[r]++
                       break
                   }
               }
           }
       }
       fmt.Fprintln(writer, sb.String())
   }
}
