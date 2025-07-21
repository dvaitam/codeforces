package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   // counts: 0:G,1:H,2:R,3:S
   c := [4]int{}
   k := 0
   for i := 0; i < n && i < len(s); i++ {
       switch s[i] {
       case 'G':
           c[0]++
       case 'H':
           c[1]++
       case 'R':
           c[2]++
       case 'S':
           c[3]++
       case '?':
           k++
       }
   }
   // Balanced distribution: sort houses by initial counts
   type pair struct{ val, idx int }
   p := []pair{{c[0], 0}, {c[1], 1}, {c[2], 2}, {c[3], 3}}
   // sort by val asc
   for i := 0; i < 4; i++ {
       for j := i + 1; j < 4; j++ {
           if p[j].val < p[i].val {
               p[i], p[j] = p[j], p[i]
           }
       }
   }
   kLeft := k
   level := p[0].val
   idxGroup := 1
   // bring groups to same level
   for idxGroup < 4 {
       need := idxGroup * (p[idxGroup].val - level)
       if kLeft >= need {
           kLeft -= need
           level = p[idxGroup].val
           idxGroup++
       } else {
           break
       }
   }
   // distribute remaining among current group
   if idxGroup > 0 {
       // full rounds possible
       f := kLeft / idxGroup
       level += f
       kLeft -= f * idxGroup
       // remaining kLeft < idxGroup
   }
   // houses that remain minimal are any of first idxGroup houses
   names := []string{"Gryffindor", "Hufflepuff", "Ravenclaw", "Slytherin"}
   res := make([]string, 0, idxGroup)
   for i := 0; i < idxGroup; i++ {
       res = append(res, names[p[i].idx])
   }
   // output in alphabetical order
   // names are already alphabetical, but res may be subset
   // so sort res
   for i := 0; i < len(res); i++ {
       for j := i + 1; j < len(res); j++ {
           if res[j] < res[i] {
               res[i], res[j] = res[j], res[i]
           }
       }
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for _, name := range res {
       fmt.Fprintln(w, name)
   }
}
