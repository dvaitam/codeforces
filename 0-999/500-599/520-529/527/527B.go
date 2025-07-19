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
   var s, t string
   fmt.Fscan(reader, &s, &t)
   // collect mismatches
   type pair struct{ i int; s, t byte }
   ms := make([]pair, 0, 16)
   for i := 0; i < n; i++ {
       if s[i] != t[i] {
           ms = append(ms, pair{i, s[i], t[i]})
       }
   }
   cnt := len(ms)
   // position matrix: pos[x][y] stores 1-based index of mismatch where s==x, t==y
   var pos [26][26]int
   for k, p := range ms {
       pos[p.s-'a'][p.t-'a'] = k + 1
   }
   // try reduce by 2
   for _, p := range ms {
       // find opposite pair to reduce two mismatches
       if opp := pos[p.t-'a'][p.s-'a']; opp != 0 {
           fmt.Println(cnt-2)
           fmt.Printf("%d %d\n", p.i+1, ms[opp-1].i+1)
           return
       }
   }
   // try reduce by 1
   for _, p := range ms {
       // find any mismatch to reduce one mismatch
       row := p.t - 'a'
       for c := 0; c < 26; c++ {
           if pos[row][c] != 0 {
               j := pos[row][c] - 1
               fmt.Println(cnt-1)
               fmt.Printf("%d %d\n", p.i+1, ms[j].i+1)
               return
           }
       }
   }
   // no improvement
   fmt.Println(cnt)
   fmt.Println("-1 -1")
}
