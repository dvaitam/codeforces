package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

// query sends a guess (x, y) and returns the interactor's response:
// 1 means x < a, 2 means y < b, 3 means x > a or y > b.
func query(x, y int64) int {
   fmt.Fprintf(writer, "? %d %d\n", x, y)
   writer.Flush()
   var res int
   if _, err := fmt.Fscan(reader, &res); err != nil {
       os.Exit(0)
   }
   return res
}

func main() {
   defer writer.Flush()
   var n int64
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // binary search for a using y = n
   var la, ra int64 = 1, n
   for la < ra {
       mid := (la + ra) / 2
       res := query(mid, n)
       if res == 1 {
           // x < a
           la = mid + 1
       } else {
           // res == 3: x >= a (or y > b)
           ra = mid
       }
   }
   a := la
   // binary search for b using x = a
   var lb, rb int64 = 1, n
   for lb < rb {
       mid := (lb + rb) / 2
       res := query(a, mid)
       if res == 2 {
           // y < b
           lb = mid + 1
       } else {
           // res == 3: y >= b (or x > a)
           rb = mid
       }
   }
   b := lb
   // output the answer
   fmt.Fprintf(writer, "! %d %d\n", a, b)
}
