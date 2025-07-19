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

   var tt int
   fmt.Fscan(reader, &tt)
   for ; tt > 0; tt-- {
       var n int
       fmt.Fscan(reader, &n)
       // Handle impossible case for n=3
       if n%2 == 1 && n == 3 {
           fmt.Fprintln(writer, "NO")
           continue
       }
       // Always possible otherwise
       fmt.Fprintln(writer, "YES")
       arr := make([]int, n)
       if n%2 == 1 {
           p := n/2 - 1
           q := n / 2
           for i := 0; i < n; i++ {
               if i%2 == 0 {
                   arr[i] = p
               } else {
                   arr[i] = -q
               }
           }
       } else {
           for i := 0; i < n; i++ {
               if i%2 == 0 {
                   arr[i] = 1
               } else {
                   arr[i] = -1
               }
           }
       }
       for i, v := range arr {
           sep := " "
           if i == n-1 {
               sep = "\n"
           }
           fmt.Fprint(writer, v, sep)
       }
   }
}
