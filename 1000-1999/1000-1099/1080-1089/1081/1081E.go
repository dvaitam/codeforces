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
   a := make([]int64, n)
   // read odd positions
   for i := 1; i < n; i += 2 {
       fmt.Fscan(reader, &a[i])
   }
   var cur int64 = 100000000000001
   var i int
   // fill even positions backward
   for i = n - 1; i >= 1; i -= 2 {
       ok := false
       ai := a[i]
       // try divisors k
       start := int64(2 - ai%2)
       for k := start; k*k <= ai; k += 2 {
           if ai%k != 0 {
               continue
           }
           miden := ai / k
           if (miden-k)%2 != 0 {
               continue
           }
           upper := miden + (k - 1)
           if upper+2 < cur {
               ok = true
               t1 := upper + 2
               // cur-2
               t2 := cur - 2
               sum := t1 + t2
               count := (cur - t1) / 2
               a[i+1] = sum * count / 2
               cur = miden - (k - 1)
               break
           }
       }
       if !ok {
           break
       }
   }
   // compute a[0]
   // sum = 1 + (cur-2) = cur-1
   sum0 := cur - 1
   count0 := (cur - 1) / 2
   a[0] = sum0 * count0 / 2
   if a[0] > 0 && i < 1 {
       fmt.Fprintln(writer, "Yes")
       for idx := 0; idx < n; idx++ {
           if idx+1 < n {
               fmt.Fprintf(writer, "%d ", a[idx])
           } else {
               fmt.Fprintf(writer, "%d\n", a[idx])
           }
       }
   } else {
       fmt.Fprintln(writer, "No")
   }
}
