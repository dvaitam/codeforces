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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var a, b, c, d int64
       fmt.Fscan(reader, &a, &b, &c, &d)
       // divisors of a
       afacts := make(map[int64]int64)
       for i := int64(1); i*i <= a; i++ {
           if a%i == 0 {
               afacts[i] = a / i
               afacts[a/i] = i
           }
       }
       // divisors of b
       bfacts := make(map[int64]int64)
       for i := int64(1); i*i <= b; i++ {
           if b%i == 0 {
               bfacts[i] = b / i
               bfacts[b/i] = i
           }
       }
       found := false
       for x1, x2 := range afacts {
           if found {
               break
           }
           for y1, y2 := range bfacts {
               x := x1 * y1
               y := x2 * y2
               k1 := a/x + 1
               k2 := b/y + 1
               if k1*x <= c && k2*y <= d {
                   fmt.Fprintln(writer, k1*x, k2*y)
                   found = true
                   break
               }
           }
       }
       if !found {
           fmt.Fprintln(writer, -1, -1)
       }
   }
}
