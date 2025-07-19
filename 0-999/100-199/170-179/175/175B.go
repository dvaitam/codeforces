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
   names := make([]string, n)
   scores := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &names[i], &scores[i])
   }
   // initial indices
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       idx[i] = i
   }
   // sort by name
   sort.Slice(idx, func(i, j int) bool {
       return names[idx[i]] < names[idx[j]]
   })
   // deduplicate names, keep highest score per name
   unique := make([]int, 0, n)
   for i := 0; i < len(idx); {
       j := i
       best := idx[i]
       // group with same name
       for j+1 < len(idx) && names[idx[j+1]] == names[idx[i]] {
           j++
           if scores[idx[j]] > scores[best] {
               best = idx[j]
           }
       }
       unique = append(unique, best)
       i = j + 1
   }
   idx = unique
   n = len(idx)
   // sort by score ascending
   sort.Slice(idx, func(i, j int) bool {
       return scores[idx[i]] < scores[idx[j]]
   })
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, n)
   // classify each
   for i := 0; i < n; i++ {
       s := scores[idx[i]]
       // find rightmost with same score
       r := i
       for r+1 < n && scores[idx[r+1]] == s {
           r++
       }
       // proportion less or equal
       nwtr := float64(r+1) / float64(n)
       // proportion strictly greater
       btr := float64(n-(r+1)) / float64(n)
       name := names[idx[i]]
       switch {
       case nwtr >= 0.99:
           fmt.Fprintf(writer, "%s pro\n", name)
       case nwtr >= 0.9 && btr > 0.01:
           fmt.Fprintf(writer, "%s hardcore\n", name)
       case nwtr >= 0.8 && btr > 0.1:
           fmt.Fprintf(writer, "%s average\n", name)
       case nwtr >= 0.5 && btr > 0.2:
           fmt.Fprintf(writer, "%s random\n", name)
       case btr > 0.5:
           fmt.Fprintf(writer, "%s noob\n", name)
       }
       // skip i to r
       // but since we print for each, we continue next
   }
}
