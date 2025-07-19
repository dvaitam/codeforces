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
   if n%4 == 2 || n%4 == 3 {
       fmt.Fprintln(writer, "NO")
       return
   }
   fmt.Fprintln(writer, "YES")
   if n == 1 {
       return
   }
   // Predefined patterns
   ans4 := [][2]int{{2, 4}, {1, 3}, {2, 3}, {1, 4}, {1, 2}, {3, 4}}
   ans5 := [][2]int{{3, 5}, {1, 3}, {4, 5}, {2, 3}, {1, 2}, {1, 5}, {2, 5}, {3, 4}, {2, 4}, {1, 4}}
   t44 := [][2]int{
       {4, 5}, {1, 6}, {3, 6}, {1, 7}, {1, 5}, {3, 8}, {3, 5}, {1, 8},
       {2, 6}, {4, 7}, {2, 7}, {2, 5}, {3, 7}, {4, 6}, {4, 8}, {2, 8},
   }
   t45 := [][2]int{
       {3, 7}, {1, 6}, {2, 9}, {1, 7}, {1, 5}, {3, 6}, {2, 7}, {4, 9}, {3, 9}, {3, 8},
       {3, 5}, {4, 5}, {4, 6}, {2, 6}, {4, 7}, {4, 8}, {1, 8}, {1, 9}, {2, 8}, {2, 5},
   }
   // Build segments of size 4 or 5
   type seg struct{ l, r int }
   var v []seg
   for i := 1; i <= n; i += 4 {
       if i+3 == n-1 {
           v = append(v, seg{i, i + 4})
           break
       }
       v = append(v, seg{i, i + 3})
   }
   // Output intra-segment swaps
   for i := 0; i < len(v); i++ {
       s := v[i]
       delta := s.l - 1
       size := s.r - s.l + 1
       if size == 4 {
           for _, p := range ans4 {
               fmt.Fprintf(writer, "%d %d\n", p[0]+delta, p[1]+delta)
           }
       } else {
           for _, p := range ans5 {
               fmt.Fprintf(writer, "%d %d\n", p[0]+delta, p[1]+delta)
           }
       }
       // Inter-segment swaps
       for j := i + 1; j < len(v); j++ {
           t := v[j]
           dj := t.l - 5
           size2 := t.r - t.l + 1
           if size2 == 4 {
               for _, p := range t44 {
                   fmt.Fprintf(writer, "%d %d\n", p[0]+delta, p[1]+dj)
               }
           } else {
               for _, p := range t45 {
                   fmt.Fprintf(writer, "%d %d\n", p[0]+delta, p[1]+dj)
               }
           }
       }
   }
