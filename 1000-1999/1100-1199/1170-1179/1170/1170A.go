package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var q int
   if _, err := fmt.Fscan(in, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var x, y int64
       fmt.Fscan(in, &x, &y)
       // Try three cases, pick minimal sum a+b+c
       bestSum := int64(1<<62 - 1)
       var ra, rb, rc int64
       // Case1: x = a+b, y = a+c
       {
           a := x
           if y < a {
               a = y
           }
           a-- // maximize a = min(x-1, y-1)
           if a >= 1 {
               b := x - a
               c := y - a
               if b >= 1 && c >= 1 {
                   s := a + b + c
                   if s < bestSum {
                       bestSum = s
                       ra, rb, rc = a, b, c
                   }
               }
           }
       }
       // Case2: x = a+b, y = b+c
       {
           b := x
           if y < b {
               b = y
           }
           b-- // maximize b
           if b >= 1 {
               a := x - b
               c := y - b
               if a >= 1 && c >= 1 {
                   s := a + b + c
                   if s < bestSum {
                       bestSum = s
                       ra, rb, rc = a, b, c
                   }
               }
           }
       }
       // Case3: x = a+c, y = b+c
       {
           c := x
           if y < c {
               c = y
           }
           c-- // maximize c
           if c >= 1 {
               a := x - c
               b := y - c
               if a >= 1 && b >= 1 {
                   s := a + b + c
                   if s < bestSum {
                       bestSum = s
                       ra, rb, rc = a, b, c
                   }
               }
           }
       }
       // Output result
       fmt.Fprintf(out, "%d %d %d\n", ra, rb, rc)
   }
}
