package main

import (
   "bufio"
   "fmt"
   "math/bits"
   "os"
)

func main() {
   // Precompute candidate masks for n from 5 to 19
   imp := make([][]int, 20)
   for n := 5; n <= 19; n++ {
       for mask := 1; mask < (1 << n); mask++ {
           var a [20]byte
           for d := 0; d < n; d++ {
               if mask&(1<<d) != 0 {
                   for j := d + 1; j <= 19; j += d + 1 {
                       a[j] ^= 1
                   }
               }
           }
           var z int
           for j := 1; j <= n; j++ {
               if a[j] != 0 {
                   z++
               }
           }
           if z <= n/5 {
               imp[n] = append(imp[n], mask)
           }
       }
   }

   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for ; t > 0; t-- {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       // Read edges
       var edges [][2]int
       for i := 0; i < m; i++ {
           var x, y int
           fmt.Fscan(reader, &x, &y)
           if x < 20 && y < 20 {
               edges = append(edges, [2]int{x, y})
           }
       }
       if n < 5 {
           fmt.Fprintln(writer, -1)
           continue
       }
       if n > 19 {
           fmt.Fprintln(writer, n)
           for i := 1; i <= n; i++ {
               if i > 1 {
                   writer.WriteByte(' ')
               }
               fmt.Fprint(writer, i)
           }
           writer.WriteByte('\n')
           continue
       }
       // find mask
       ans := -1
       for _, mask := range imp[n] {
           bad := false
           for _, e := range edges {
               x, y := e[0], e[1]
               if mask&(1<<(x-1)) != 0 && mask&(1<<(y-1)) == 0 {
                   bad = true
                   break
               }
           }
           if !bad {
               ans = mask
               break
           }
       }
       if ans < 0 {
           fmt.Fprintln(writer, -1)
           continue
       }
       cnt := bits.OnesCount(uint(ans))
       fmt.Fprintln(writer, cnt)
       first := true
       for i := 0; i < n; i++ {
           if ans&(1<<i) != 0 {
               if !first {
                   writer.WriteByte(' ')
               }
               first = false
               fmt.Fprint(writer, i+1)
           }
       }
       writer.WriteByte('\n')
   }
}
