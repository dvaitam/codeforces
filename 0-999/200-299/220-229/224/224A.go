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
   var x, y, z int64
   if _, err := fmt.Fscan(reader, &x, &y, &z); err != nil {
       return
   }
   a := int64(math.Sqrt(float64(x)*float64(y)/float64(z)))
   b := int64(math.Sqrt(float64(x)*float64(z)/float64(y)))
   c := int64(math.Sqrt(float64(y)*float64(z)/float64(x)))
   result := 4*(a+b+c)
   fmt.Fprintln(writer, result)
}
