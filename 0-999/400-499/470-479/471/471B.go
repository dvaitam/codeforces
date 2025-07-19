package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct { val, pos int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   pairs := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pairs[i].val)
       pairs[i].pos = i + 1
   }
   sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })
   l1, l2 := -1, -1
   for i := 0; i < n-1; i++ {
       if pairs[i].val == pairs[i+1].val {
           if l1 == -1 {
               l1 = i
           } else if l2 == -1 {
               l2 = i
           }
       }
   }
   if l2 == -1 {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   base := make([]int, n)
   for i := 0; i < n; i++ {
       base[i] = pairs[i].pos
   }
   // first sequence
   for i, v := range base {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
   // second sequence: swap at l1
   seq2 := append([]int(nil), base...)
   seq2[l1], seq2[l1+1] = seq2[l1+1], seq2[l1]
   for i, v := range seq2 {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
   // third sequence: swap at l2
   seq3 := append([]int(nil), base...)
   seq3[l2], seq3[l2+1] = seq3[l2+1], seq3[l2]
   for i, v := range seq3 {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   fmt.Fprintln(writer)
}
