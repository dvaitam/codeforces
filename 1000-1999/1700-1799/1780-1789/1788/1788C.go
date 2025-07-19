package main

import (
   "bufio"
   "fmt"
   "os"
)

// findx finds smallest x >=1 such that (x+n)*(x+n-1) - x*(x-1) == 2*n*(2*n+1), or -1 if none.
func findx(n int64) int64 {
   var l, h int64 = 1, 1000000
   for l <= h {
       x := (l + h) / 2
       if (x+n)*(x+n-1) - x*(x-1) >= 2*n*(2*n+1) {
           h = x - 1
       } else {
           l = x + 1
       }
   }
   if (l+n)*(l+n-1) - l*(l-1) == 2*n*(2*n+1) {
       return l
   }
   return -1
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ti := 0; ti < t; ti++ {
       var n int64
       fmt.Fscan(reader, &n)
       x := findx(n)
       if x == -1 {
           fmt.Fprintln(writer, "NO")
       } else {
           fmt.Fprintln(writer, "YES")
           var i1 int64 = 1
           j := x - 1
           // first segment
           for ; j <= 2*n; i1, j = i1+1, j+1 {
               fmt.Fprintln(writer, i1, j)
           }
           // second segment
           j = x + 1 - i1
           for ; j < x-1; i1, j = i1+1, j+1 {
               fmt.Fprintln(writer, i1, j)
           }
       }
   }
}
