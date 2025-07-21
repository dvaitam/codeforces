package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   g := gcd(n, m)
   comp := make([]bool, g)
   var b int
   fmt.Fscan(reader, &b)
   for i := 0; i < b; i++ {
       var x int
       fmt.Fscan(reader, &x)
       comp[x%g] = true
   }
   var girls int
   fmt.Fscan(reader, &girls)
   for i := 0; i < girls; i++ {
       var y int
       fmt.Fscan(reader, &y)
       comp[y%g] = true
   }
   for i := 0; i < g; i++ {
       if !comp[i] {
           fmt.Println("No")
           return
       }
   }
   fmt.Println("Yes")
}

// gcd returns the greatest common divisor of a and b.
func gcd(a, b int) int {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}
