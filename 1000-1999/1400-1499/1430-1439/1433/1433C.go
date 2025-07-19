package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   var x int
   var sign = 1
   ch, err := reader.ReadByte()
   for err == nil && (ch < '0' || ch > '9') {
       if ch == '-' {
           sign = -1
       }
       ch, err = reader.ReadByte()
   }
   for err == nil && ch >= '0' && ch <= '9' {
       x = x*10 + int(ch-'0')
       ch, err = reader.ReadByte()
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   t := readInt()
   const INF = 1000000005
   for t > 0 {
       t--
       n := readInt()
       arr := make([]int, n+2)
       for i := 1; i <= n; i++ {
           arr[i] = readInt()
       }
       arr[0], arr[n+1] = INF, INF
       mx := arr[1]
       for i := 2; i <= n; i++ {
           if arr[i] > mx {
               mx = arr[i]
           }
       }
       idx := -1
       for i := 1; i <= n; i++ {
           if arr[i] == mx && (arr[i-1] < mx || arr[i+1] < mx) {
               idx = i
               break
           }
       }
       if idx == -1 {
           fmt.Fprintln(writer, -1)
       } else {
           fmt.Fprintln(writer, idx)
       }
   }
}
