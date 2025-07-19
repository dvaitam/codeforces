package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n+1)
   t := 1
   for t <= n && m >= (t-1)/2 {
       a[t] = t
       m -= (t - 1) / 2
       t++
   }
   if t > n && m > 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   if m > 0 {
       // place a[t] to satisfy remaining m
       a[t] = 2*t - 2*m - 1
       t++
   }
   // fill the rest with large values
   for ; t <= n; t++ {
       a[t] = n*n + t*n
   }
   // output the array
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteString(" ")
       }
       writer.WriteString(strconv.Itoa(a[i]))
   }
   writer.WriteString("\n")
}
