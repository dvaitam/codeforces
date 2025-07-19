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

   var T int
   if _, err := fmt.Fscan(in, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       stack := make([]int, 0, n)
       found := false
       var x, y, z int
       for i := 0; i < n; i++ {
           for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
               if len(stack) > 1 && !found {
                   x = stack[len(stack)-2] + 1
                   y = stack[len(stack)-1] + 1
                   z = i + 1
                   found = true
               }
               stack = stack[:len(stack)-1]
           }
           stack = append(stack, i)
       }
       if found {
           fmt.Fprintln(out, "YES")
           fmt.Fprintf(out, "%d %d %d", x, y, z)
       } else {
           fmt.Fprintln(out, "NO")
       }
       if T > 0 {
           fmt.Fprintln(out)
       }
   }
}
