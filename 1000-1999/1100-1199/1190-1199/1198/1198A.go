package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   var I int64
   fmt.Fscan(reader, &n, &I)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   sort.Ints(a)
   // compress values and count frequencies
   vals := make([]int, 0, n)
   cnt := make([]int, 0, n)
   for i := 0; i < n; {
       j := i + 1
       for j < n && a[j] == a[i] {
           j++
       }
       vals = append(vals, a[i])
       cnt = append(cnt, j-i)
       i = j
   }
   m := len(vals)
   // max bits per element
   totalBits := I * 8
   // B bits allowed per element (floor)
   B := totalBits / int64(n)
   // if too many bits, no need to change
   if B >= 31 {
       fmt.Fprintln(writer, 0)
       return
   }
   // maximum distinct values allowed
   var W int
   if B < 0 {
       W = 1
   } else {
       W = 1 << uint(B)
   }
   if W >= m {
       fmt.Fprintln(writer, 0)
       return
   }
   // prefix sums of counts
   pref := make([]int, m+1)
   for i := 0; i < m; i++ {
       pref[i+1] = pref[i] + cnt[i]
   }
   // sliding window of size W over distinct values
   best := 0
   for i := 0; i+W <= m; i++ {
       j := i + W - 1
       sum := pref[j+1] - pref[i]
       if sum > best {
           best = sum
       }
   }
   // minimal changes = total elements - max kept
   fmt.Fprintln(writer, n-best)
}
