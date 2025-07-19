package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var x float64
   if _, err := fmt.Fscan(reader, &x); err != nil {
       return
   }
   // compute log base 2 of x
   result := math.Log2(x)
   fmt.Fprint(writer, result)
}
