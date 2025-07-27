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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n int
       fmt.Fscan(reader, &n)
       b := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &b[i])
       }
       // number of odd positions in original a
       kOdd := (n + 1) / 2
       kEven := n / 2
       a := make([]int, n)
       // fill odd indices (1-based): positions 1,3,5... -> zero-based 0,2,4...
       for i := 0; i < kOdd; i++ {
           a[2*i] = b[i]
       }
       // fill even indices in reverse order: a[2*kEven], a[2*(kEven-1)],... a[2]
       // zero-based index: 2*(kEven-j)-1 for j from 0 to kEven-1
       for j := 0; j < kEven; j++ {
           idx := 2*(kEven-j) - 1
           a[idx] = b[kOdd+j]
       }
       // output a
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(fmt.Sprint(v))
       }
       writer.WriteByte('\n')
   }
}
