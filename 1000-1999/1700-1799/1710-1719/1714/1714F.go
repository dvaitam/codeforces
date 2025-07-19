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
   for tt > 0 {
       tt--
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, d12, d23, d13 int
   fmt.Fscan(reader, &n, &d12, &d23, &d13)
   // Case root = 1
   if d23 == d12 + d13 {
       if d12 + d13 + 1 > n {
           fmt.Fprintln(writer, "NO")
           return
       }
       fmt.Fprintln(writer, "YES")
       node := 4
       second := 1
       // chain to 2
       for i := 1; i <= d12-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 2)
       // chain to 3
       second = 1
       for i := 1; i <= d13-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 3)
       // attach remaining to root 1
       for node <= n {
           fmt.Fprintln(writer, 1, node)
           node++
       }
       return
   }
   // Case root = 3
   if d12 == d23 + d13 {
       if d23 + d13 + 1 > n {
           fmt.Fprintln(writer, "NO")
           return
       }
       fmt.Fprintln(writer, "YES")
       node := 4
       second := 3
       // chain to 1
       for i := 1; i <= d13-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 1)
       // chain to 2
       second = 3
       for i := 1; i <= d23-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 2)
       for node <= n {
           fmt.Fprintln(writer, 1, node)
           node++
       }
       return
   }
   // Case root = 2
   if d13 == d12 + d23 {
       if d12 + d23 + 1 > n {
           fmt.Fprintln(writer, "NO")
           return
       }
       fmt.Fprintln(writer, "YES")
       node := 4
       second := 2
       // chain to 1
       for i := 1; i <= d12-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 1)
       // chain to 3
       second = 2
       for i := 1; i <= d23-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 3)
       for node <= n {
           fmt.Fprintln(writer, 1, node)
           node++
       }
       return
   }
   // General root = 4
   for j := 1; j < d12; j++ {
       d41 := j
       d42 := d12 - j
       d43 := d13 - j
       if d43 <= 0 || d42 <= 0 || d42 + d43 != d23 || d41 + d42 + d43 + 1 > n {
           continue
       }
       fmt.Fprintln(writer, "YES")
       node := 5
       // branch to 1
       second := 4
       for i := 1; i <= d41-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 1)
       // branch to 3
       second = 4
       for i := 1; i <= d43-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 3)
       // branch to 2
       second = 4
       for i := 1; i <= d42-1; i++ {
           first := second
           second = node
           fmt.Fprintln(writer, first, second)
           node++
       }
       fmt.Fprintln(writer, second, 2)
       // attach remaining to 4
       for node <= n {
           fmt.Fprintln(writer, 4, node)
           node++
       }
       return
   }
   fmt.Fprintln(writer, "NO")
}
