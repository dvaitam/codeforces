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
   var t int
   fmt.Fscan(reader, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(reader, &n)
       n *= 2
       var trumpStr string
       fmt.Fscan(reader, &trumpStr)
       trump := trumpStr[0]
       cards := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &cards[i])
       }
       sort.Strings(cards)
       m := make(map[byte][]string)
       for _, card := range cards {
           if len(card) < 2 {
               continue
           }
           suit := card[1]
           m[suit] = append(m[suit], card)
       }
       // Prepare grouping
       var rem int
       var maxStr []string
       // order of suits deterministic
       var suits []int
       for k := range m {
           suits = append(suits, int(k))
       }
       sort.Ints(suits)
       var ans [][2]string
       var remStr []string
       maxCount := 0
       // process suits
       for _, ki := range suits {
           k := byte(ki)
           v := m[k]
           if k == trump {
               maxCount = len(v)
               maxStr = v
           } else {
               // pair within same suit
               for i := 0; i+1 < len(v); i += 2 {
                   ans = append(ans, [2]string{v[i], v[i+1]})
               }
               if len(v)%2 == 1 {
                   rem++
                   remStr = append(remStr, v[len(v)-1])
               }
           }
       }
       // check feasibility
       if maxCount < rem || (maxCount-rem)%2 != 0 {
           fmt.Fprintln(writer, "IMPOSSIBLE")
           continue
       }
       // match remStr with trumps
       idx := 0
       for _, c := range remStr {
           ans = append(ans, [2]string{c, maxStr[idx]})
           idx++
       }
       // pair remaining trumps
       for idx+1 < len(maxStr) {
           ans = append(ans, [2]string{maxStr[idx], maxStr[idx+1]})
           idx += 2
       }
       // output
       for _, p := range ans {
           fmt.Fprintln(writer, p[0], p[1])
       }
   }
}
