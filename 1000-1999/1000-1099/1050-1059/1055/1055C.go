package main

import (
   "bufio"
   "fmt"
   "os"
)

func gcd(a, b int64) int64 {
   for b != 0 {
       a, b = b, a%b
   }
   return a
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func min(a, b int64) int64 {
   if a < b {
       return a
   }
   return b
}

func abs(a int64) int64 {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var l1, r1, t1, l2, r2, t2 int64
   if _, err := fmt.Fscan(reader, &l1, &r1, &t1); err != nil {
       return
   }
   fmt.Fscan(reader, &l2, &r2, &t2)

   ans := max(0, min(r1, r2)-max(l1, l2)+1)
   f1 := r1 - l1 + 1
   f2 := r2 - l2 + 1
   g := gcd(t1, t2)
   k := abs(l1 - l2) % g

   if t1 == t2 {
       fmt.Fprint(writer, ans)
       return
   }
   if k == 0 {
       fmt.Fprint(writer, max(ans, min(f1, f2)))
       return
   }
   if l2 > l1 {
       if f2 + k <= f1 {
           fmt.Fprint(writer, max(ans, f2))
           return
       }
       if f1 + g - k <= f2 {
           fmt.Fprint(writer, max(ans, f1))
           return
       }
       fmt.Fprint(writer, max(ans, max(f1 - k, f2 - g + k)))
       return
   }
   if f1 + k <= f2 {
       fmt.Fprint(writer, max(ans, f1))
       return
   }
   if f2 + g - k <= f1 {
       fmt.Fprint(writer, max(ans, f2))
       return
   }
   fmt.Fprint(writer, max(ans, max(f2 - k, f1 - g + k)))
}
