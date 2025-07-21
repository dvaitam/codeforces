package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   names := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &names[i])
   }
   surnames := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &surnames[i])
   }
   // initial counts of starting letters
   cntN := make([]int, 26)
   cntS := make([]int, 26)
   for i := 0; i < n; i++ {
       cntN[names[i][0]-'A']++
       cntS[surnames[i][0]-'A']++
   }
   // max special matches
   K := 0
   for c := 0; c < 26; c++ {
       if cntN[c] < cntS[c] {
           K += cntN[c]
       } else {
           K += cntS[c]
       }
   }
   // prepare all pairs
   type pair struct {
       i, j    int
       special int
       str     string
   }
   pairs := make([]pair, 0, n*n)
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           s := names[i] + " " + surnames[j]
           sp := 0
           if names[i][0] == surnames[j][0] {
               sp = 1
           }
           pairs = append(pairs, pair{i, j, sp, s})
       }
   }
   sort.Slice(pairs, func(a, b int) bool {
       return pairs[a].str < pairs[b].str
   })
   usedN := make([]bool, n)
   usedS := make([]bool, n)
   res := make([]string, 0, n)
   remK := K
   // mutable counts
   curN := make([]int, 26)
   curS := make([]int, 26)
   copy(curN, cntN)
   copy(curS, cntS)
   // greedy pick
   for t := 0; t < n; t++ {
       for _, p := range pairs {
           if usedN[p.i] || usedS[p.j] {
               continue
           }
           if remK-p.special < 0 {
               continue
           }
           // simulate remove
           ci := names[p.i][0] - 'A'
           cj := surnames[p.j][0] - 'A'
           curN[ci]--
           curS[cj]--
           // compute possible special
           poss := 0
           for c := 0; c < 26; c++ {
               if curN[c] < curS[c] {
                   poss += curN[c]
               } else {
                   poss += curS[c]
               }
           }
           if poss >= remK-p.special {
               // accept
               usedN[p.i] = true
               usedS[p.j] = true
               res = append(res, p.str)
               remK -= p.special
               break
           }
           // rollback
           curN[ci]++
           curS[cj]++
       }
   }
   // output sorted result (already in lex order)
   out := strings.Join(res, ", ")
   fmt.Fprintln(os.Stdout, out)
}
