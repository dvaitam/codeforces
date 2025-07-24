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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   evenCount := 0
   for i := 0; i < n; i++ {
       var a int
       fmt.Fscan(reader, &a)
       if a%2 == 0 {
           evenCount++
       }
       // First player wins if number of even cycles is odd
       if evenCount%2 == 1 {
           writer.WriteString("1")
       } else {
           writer.WriteString("2")
       }
       if i < n-1 {
           writer.WriteByte('\n')
       }
   }
}
