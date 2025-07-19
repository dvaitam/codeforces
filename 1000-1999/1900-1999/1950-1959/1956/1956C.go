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

   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n int
   fmt.Fscan(reader, &n)
   var sum int64
   for i := 1; i <= n; i++ {
       sum += int64(2*i-1) * int64(i)
   }

   // total operations and max sequence length
   writer.WriteString(strconv.FormatInt(sum, 10))
   writer.WriteByte(' ')
   writer.WriteString(strconv.Itoa(2*n - 1))
   writer.WriteByte('\n')

   // first operation
   writer.WriteString("1 ")
   writer.WriteString(strconv.Itoa(n))
   writer.WriteByte(' ')
   for i := 1; i <= n; i++ {
       writer.WriteString(strconv.Itoa(i))
       if i == n {
           writer.WriteByte('\n')
       } else {
           writer.WriteByte(' ')
       }
   }

   // subsequent operations
   for i := 1; i < n; i++ {
       // operation type 2
       writer.WriteString("2 ")
       writer.WriteString(strconv.Itoa(n - i))
       writer.WriteByte(' ')
       for j := 1; j <= n; j++ {
           writer.WriteString(strconv.Itoa(j))
           if j == n {
               writer.WriteByte('\n')
           } else {
               writer.WriteByte(' ')
           }
       }

       // operation type 1
       writer.WriteString("1 ")
       writer.WriteString(strconv.Itoa(n - i))
       writer.WriteByte(' ')
       for j := 1; j <= n; j++ {
           writer.WriteString(strconv.Itoa(j))
           if j == n {
               writer.WriteByte('\n')
           } else {
               writer.WriteByte(' ')
           }
       }
   }
}
