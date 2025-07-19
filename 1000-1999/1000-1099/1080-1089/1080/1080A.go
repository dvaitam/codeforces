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

   var n, k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   // Compute notebooks needed for each color: red(2), green(5), blue(8)
   red := (n*2 + k - 1) / k
   green := (n*5 + k - 1) / k
   blue := (n*8 + k - 1) / k
   result := red + green + blue
   fmt.Fprintln(writer, result)
}
