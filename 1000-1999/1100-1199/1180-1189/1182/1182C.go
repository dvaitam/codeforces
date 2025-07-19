package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n int
   fmt.Fscan(in, &n)
   lines := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &lines[i])
   }
   // vowels mapping a,e,i,o,u to 0..4
   vmap := map[rune]int{
       'a': 0, 'e': 1, 'i': 2, 'o': 3, 'u': 4,
   }
   num := make([]int, n)
   bel := make([]int, n)
   vowelGroups := make([][]int, 5)
   // compute count and last vowel
   for i, s := range lines {
       cnt := 0
       last := -1
       for _, ch := range s {
           if vi, ok := vmap[ch]; ok {
               cnt++
               last = vi
           }
       }
       num[i] = cnt
       bel[i] = last
       if last >= 0 {
           vowelGroups[last] = append(vowelGroups[last], i)
       }
   }
   // perfect rhyme pairs: same count and same last vowel
   perfect := make([][2]int, 0, n/2)
   used := make([]bool, n)
   // for each vowel group
   for v := 0; v < 5; v++ {
       // group by count
       mpc := make(map[int][]int)
       for _, idx := range vowelGroups[v] {
           mpc[num[idx]] = append(mpc[num[idx]], idx)
       }
       for _, lst := range mpc {
           for i := 0; i+1 < len(lst); i += 2 {
               perfect = append(perfect, [2]int{lst[i], lst[i+1]})
               used[lst[i]] = true
               used[lst[i+1]] = true
           }
       }
   }
   // count-only pairs from remaining
   countOnly := make([][2]int, 0, n/2)
   mpc2 := make(map[int][]int)
   for i := 0; i < n; i++ {
       if used[i] {
           continue
       }
       mpc2[num[i]] = append(mpc2[num[i]], i)
   }
   for _, lst := range mpc2 {
       for i := 0; i+1 < len(lst); i += 2 {
           countOnly = append(countOnly, [2]int{lst[i], lst[i+1]})
       }
   }
   // number of stanzas
   m := len(perfect)
   if len(countOnly) < m {
       m = len(countOnly)
   }
   fmt.Fprintln(out, m)
   // output m stanzas
   for i := 0; i < m; i++ {
       // first line: first words of count-only and perfect
       fmt.Fprintln(out, lines[countOnly[i][0]], lines[perfect[i][0]])
       // second line: second words
       fmt.Fprintln(out, lines[countOnly[i][1]], lines[perfect[i][1]])
   }
}
