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

   var q int
   if _, err := fmt.Fscan(reader, &q); err != nil {
       return
   }
   for qi := 0; qi < q; qi++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       if n == 0 {
           fmt.Fprintln(writer, 0)
           continue
       }
       sort.Ints(a)
       // count frequencies
       freqs := make([]int, 0, n)
       cnt := 1
       for i := 1; i < n; i++ {
           if a[i] == a[i-1] {
               cnt++
           } else {
               freqs = append(freqs, cnt)
               cnt = 1
           }
       }
       freqs = append(freqs, cnt)
       // sort descending
       sort.Sort(sort.Reverse(sort.IntSlice(freqs)))
       // greedy pick distinct counts
       last := n + 1
       total := 0
       for _, f := range freqs {
           if last <= 1 {
               break
           }
           // we can pick at most last-1
           pick := f
           if pick > last-1 {
               pick = last - 1
           }
           if pick <= 0 {
               break
           }
           total += pick
           last = pick
       }
       fmt.Fprintln(writer, total)
   }
}
