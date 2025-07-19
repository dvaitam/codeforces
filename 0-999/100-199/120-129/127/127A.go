package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   var a, b int
   fmt.Fscan(reader, &n, &k, &a, &b)
   var d float64
   for i := 1; i < n; i++ {
       var x, y int
       fmt.Fscan(reader, &x, &y)
       dx := float64(a - x)
       dy := float64(b - y)
       d += math.Hypot(dx, dy)
       a = x
       b = y
   }
   result := d * float64(k) / 50.0
   fmt.Printf("%.8f\n", result)
}
