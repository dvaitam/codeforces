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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for i := 0; i < q; i++ {
       var a, b, c int64
       fmt.Fscan(reader, &a, &b, &c)
       if a > b {
           a, b = b, a
       }
       // minimal sum reachable is b
       if b > c {
           fmt.Fprintln(writer, -1)
           continue
       }
       var ans int64
       if a != b {
           ans = a
           cha := b - a
           // if difference is odd, best add c-a-1
           if cha&1 != 0 {
               ans += c - a - 1
           } else {
               // cha even
               if (c - a)&1 == 0 {
                   ans += c - a
               } else {
                   ans += c - a - 2
               }
           }
       } else {
           if a == c {
               ans = a
           } else {
               if (c - a)&1 == 0 {
                   ans = c
               } else {
                   ans = c - 2
               }
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
