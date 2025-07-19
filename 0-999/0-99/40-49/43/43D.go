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

   var n, m int
   fmt.Fscan(reader, &n, &m)
   // initial output: number of extra moves and optional endpoints
   if (n == 1 && m == 1) || (n == 1 && m == 2) || (n == 2 && m == 1) {
       fmt.Fprintln(writer, 0)
   } else if n == 1 || m == 1 {
       fmt.Fprintln(writer, 1)
       fmt.Fprintf(writer, "%d %d 1 1\n", n, m)
   } else if (n*m)%2 == 0 {
       fmt.Fprintln(writer, 0)
   } else {
       fmt.Fprintln(writer, 1)
       fmt.Fprintf(writer, "%d %d 1 1\n", n, m)
   }

   // traversal path
   if n == 1 && m == 1 {
       fmt.Fprintln(writer, "1 1")
   } else if n == 1 {
       for j := 1; j <= m; j++ {
           fmt.Fprintf(writer, "1 %d\n", j)
       }
   } else if m == 1 {
       for i := 1; i <= n; i++ {
           fmt.Fprintf(writer, "%d 1\n", i)
       }
   } else if m%2 == 0 {
       // even columns
       fmt.Fprintln(writer, "1 1")
       for c := 1; c <= m; c += 2 {
           for r := 2; r <= n; r++ {
               fmt.Fprintf(writer, "%d %d\n", r, c)
           }
           for r := n; r >= 2; r-- {
               fmt.Fprintf(writer, "%d %d\n", r, c+1)
           }
       }
       for c := m; c >= 2; c-- {
           fmt.Fprintf(writer, "1 %d\n", c)
       }
   } else if n%2 == 0 {
       // even rows
       fmt.Fprintln(writer, "1 1")
       for r := 1; r <= n; r += 2 {
           for c := 2; c <= m; c++ {
               fmt.Fprintf(writer, "%d %d\n", r, c)
           }
           for c := m; c >= 2; c-- {
               fmt.Fprintf(writer, "%d %d\n", r+1, c)
           }
       }
       for r := n; r >= 2; r-- {
           fmt.Fprintf(writer, "%d 1\n", r)
       }
   } else {
       // both odd
       for r := 1; r <= n; r += 2 {
           for c := 1; c <= m; c++ {
               fmt.Fprintf(writer, "%d %d\n", r, c)
           }
           if r == n {
               break
           }
           for c := m; c >= 1; c-- {
               fmt.Fprintf(writer, "%d %d\n", r+1, c)
           }
       }
   }
   // return to start
   fmt.Fprintln(writer, "1 1")
}
