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

   var n, q int
   fmt.Fscan(reader, &n, &q)
   a := make([]int64, 2*n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // extend array for circular
   for i := 0; i < n; i++ {
       a[n+i] = a[i]
   }
   r := make([]int, 2*n)
   dq := make([]int, 0, n+2)
   for qi := 0; qi < q; qi++ {
       var b int64
       fmt.Fscan(reader, &b)
       // two pointers to build r
       var sum int64
       j := 0
       for i := 0; i < 2*n; i++ {
           for j < 2*n && sum + a[j] <= b {
               sum += a[j]
               j++
           }
           r[i] = j
           sum -= a[i]
       }
       // compute minimal jumps covering n length for each start
       dq = dq[:0]
       head := 0
       // start at 0
       cur := 0
       dq = append(dq, cur)
       for cur < n {
           cur = r[cur]
           dq = append(dq, cur)
       }
       ans := len(dq) - head - 1
       // slide start
       for s := 1; s < n; s++ {
           // pop while front < s
           for head < len(dq) && dq[head] < s {
               head++
           }
           // extend back until cover s+n
           last := dq[len(dq)-1]
           for last < s+n {
               last = r[last]
               dq = append(dq, last)
           }
           // update answer
           cnt := len(dq) - head - 1
           if cnt < ans {
               ans = cnt
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
