package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   b := make([]int, n)
   // count ones
   k1 := 0
   for k1 < n && a[k1] == 1 {
       k1++
   }
   if k1 < n {
       // we can replace one non-1 with 1, giving k1+1 ones
       // positions 0..k1 set to 1
       for i := 0; i <= k1 && i < n; i++ {
           b[i] = 1
       }
       // remaining positions shift original
       for i := k1 + 1; i < n; i++ {
           b[i] = a[i-1]
       }
   } else {
       // all ones: must replace one 1 with 2
       for i := 0; i < n-1; i++ {
           b[i] = 1
       }
       if n > 0 {
           b[n-1] = 2
       }
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range b {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprint(v))
   }
   writer.WriteByte('\n')
}
