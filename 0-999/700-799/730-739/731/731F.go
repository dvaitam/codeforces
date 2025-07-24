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
   const maxVal = 200000
   freq := make([]int, maxVal+1)
   var x int
   var actualMax int
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &x)
       if x >= 0 && x <= maxVal {
           freq[x]++
           if x > actualMax {
               actualMax = x
           }
       }
   }
   // prefix sums of counts
   pref := make([]int, actualMax+1)
   for i := 1; i <= actualMax; i++ {
       pref[i] = pref[i-1] + freq[i]
   }
   var ans int64
   // try each possible leading power value present
   for v := 1; v <= actualMax; v++ {
       if freq[v] == 0 {
           continue
       }
       var total int64
       // sum over blocks [j .. j+v-1]
       for j := v; j <= actualMax; j += v {
           r := j + v - 1
           if r > actualMax {
               r = actualMax
           }
           cnt := pref[r] - pref[j-1]
           total += int64(cnt) * int64(j)
       }
       if total > ans {
           ans = total
       }
   }
   fmt.Fprintln(writer, ans)
}
