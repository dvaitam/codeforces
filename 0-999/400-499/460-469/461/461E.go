package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int64
   // Read n
   if _, err := fmt.Fscan(reader, &n); err != nil {
       fmt.Fprintln(os.Stderr, "failed to read n:", err)
       return
   }
   var t string
   if _, err := fmt.Fscan(reader, &t); err != nil {
       fmt.Fprintln(os.Stderr, "failed to read t:", err)
       return
   }
   // Compute max run length for each character A-D
   runs := map[rune]int64{'A': 0, 'B': 0, 'C': 0, 'D': 0}
   var prev rune
   var curRun int64
   for i, ch := range t {
       if i == 0 || ch != prev {
           // start new run
           curRun = 1
       } else {
           curRun++
       }
       // update
       if curRun > runs[ch] {
           runs[ch] = curRun
       }
       prev = ch
   }
   // find minimal run among A-D
   rmin := runs['A']
   for _, c := range []rune{'B', 'C', 'D'} {
       if runs[c] < rmin {
           rmin = runs[c]
       }
   }
   if rmin <= 0 {
       rmin = 1 // safety, though per problem each appears
   }
   // Compute answer: ceil(n / rmin)
   moves := (n + rmin - 1) / rmin
   writer := bufio.NewWriter(os.Stdout)
   fmt.Fprintln(writer, moves)
   writer.Flush()
}
