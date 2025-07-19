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
   arr := make([]int, n)
   rem := [3]int{}
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &arr[i])
       rem[arr[i]%3]++
   }
   half := n / 2
   var pick0, pick1, pick2 int
   var resultZ int
   switch {
   case rem[0] <= half:
       pick0 = rem[0]
       if pick0+rem[1] <= half {
           pick1 = rem[1]
           pick2 = half - pick0 - pick1
       } else {
           pick1 = half - pick0
           pick2 = 0
       }
       resultZ = 0
   case rem[0] == 0 || rem[0] == n:
       resultZ = 1
       if rem[0] == n {
           pick0 = half
           pick1 = 0
           pick2 = 0
       } else {
           if rem[1] <= half {
               pick0 = 0
               pick1 = rem[1]
               pick2 = half - pick1
           } else {
               pick0 = 0
               pick1 = half
               pick2 = 0
           }
       }
   case rem[0] >= half:
       pick0 = half
       pick1 = 0
       pick2 = 0
       resultZ = 2
   default:
       fmt.Fprintln(writer, -1)
       return
   }
   fmt.Fprintln(writer, resultZ)
   // build selection string
   s := make([]byte, n)
   for i := range s {
       s[i] = '0'
   }
   // mark picks in order
   for i := 0; i < n; i++ {
       m := arr[i] % 3
       if m == 0 && pick0 > 0 {
           s[i] = '1'
           pick0--
       } else if m == 1 && pick1 > 0 {
           s[i] = '1'
           pick1--
       } else if m == 2 && pick2 > 0 {
           s[i] = '1'
           pick2--
       }
   }
   fmt.Fprintln(writer, string(s))
}
