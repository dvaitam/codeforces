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

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   if n == 3 || n == 5 {
       fmt.Fprintln(writer, "NO")
       return
   }
   ara := make([]int, n)
   if n%2 == 0 {
       p := []int{-1, 1, -1, -2, 2, 1}
       for i := 0; i < n; i++ {
           ara[i] = p[i%6]
       }
   } else {
       p := []int{2, -2, -3, 3, 1, -1}
       s := []int{3, 3, -2, -1, 1, -1, 2}
       for i := 0; i < n; i++ {
           if i < len(s) {
               ara[i] = s[i]
           } else {
               ara[i] = p[i%6]
           }
       }
   }
   fmt.Fprintln(writer, "YES")
   for i, v := range ara {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
