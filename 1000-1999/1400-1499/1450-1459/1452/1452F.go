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
   if _, err := fmt.Fscanf(reader, "%d %d", &n, &q); err != nil {
       return
   }
   cnt := make([]int64, n)
   for i := 0; i < n; i++ {
       fmt.Fscanf(reader, "%d", &cnt[i])
   }
   for qi := 0; qi < q; qi++ {
       var typ int
       fmt.Fscanf(reader, "%d", &typ)
       if typ == 1 {
           var pos int
           var val int64
           fmt.Fscanf(reader, "%d %d", &pos, &val)
           cnt[pos] = val
       } else if typ == 2 {
           var x int
           var k int64
           fmt.Fscanf(reader, "%d %d", &x, &k)
           m := n - x - 1
           // c[j]: counts by level j, where j=0 small pieces
           c := make([]int64, max(1, m+1))
           // c[0]: pieces <=2^x
           var small int64
           for i := 0; i <= x; i++ {
               small += cnt[i]
           }
           c[0] = small
           for j := 1; j <= m; j++ {
               c[j] = cnt[x+j]
           }
           need := k
           if c[0] >= need {
               fmt.Fprintln(writer, 0)
               continue
           }
           // total potential small units
           var totUnits int64
           for j, v := range c {
               if v == 0 {
                   continue
               }
               totUnits += v << j
           }
           if totUnits < need {
               fmt.Fprintln(writer, -1)
               continue
           }
           // greedy operations
           var ans int64
           need2 := need - c[0]
           for need2 > 0 {
               // B1 splitting to small
               r1 := (need2 + 1) / 2
               var avail1 int64
               if len(c) > 1 {
                   avail1 = c[1]
               }
               need1 := r1 - avail1
               if need1 <= 0 {
                   ans += r1
                   break
               }
               // produce B1 from deeper
               jj := 2
               for jj <= m && (jj >= len(c) || c[jj] == 0) {
                   jj++
               }
               if jj > m {
                   ans = -1
                   break
               }
               // one c[jj] yields sj B1 units if fully split
               sj := int64(1) << (jj - 1)
               needCj := (need1 + sj - 1) / sj
               if needCj > c[jj] {
                   needCj = c[jj]
               }
               ans += needCj
               c[jj] -= needCj
               if jj-1 >= 0 {
                   c[jj-1] += needCj * 2
               }
           }
           fmt.Fprintln(writer, ans)
       }
   }
}

// max returns the maximum of two ints
func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}
