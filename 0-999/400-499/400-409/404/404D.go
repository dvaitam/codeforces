package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   // special case n == 1: only one cell
   if n == 1 {
       switch s[0] {
       case '?': fmt.Println(2)
       case '*': fmt.Println(1)
       case '0': fmt.Println(1)
       default: fmt.Println(0)
       }
       return
   }
   // pos by parity and index mapping
   pos := make([][]int, 2)
   idx := make([]int, n)
   for i := 0; i < n; i++ {
       p := i & 1
       idx[i] = len(pos[p])
       pos[p] = append(pos[p], i)
   }
   total := int64(1)
   // process each parity
   for p := 0; p < 2; p++ {
       m := len(pos[p])
       if m == 0 {
           continue
       }
       forced := make([]int, m)
       for i := 0; i < m; i++ {
           forced[i] = -1
       }
       c := make([]int, m-1)
       for i := range c {
           c[i] = -1
       }
       // build constraints for this parity and opposite
       for i := 0; i < n; i++ {
           ch := s[i]
           if ch == '*' {
               if i&1 == p {
                   j := idx[i]
                   if forced[j] == 0 {
                       fmt.Println(0)
                       return
                   }
                   forced[j] = 1
               }
           } else if ch >= '0' && ch <= '2' {
               d0 := int(ch - '0')
               if i&1 == p {
                   j := idx[i]
                   if forced[j] == 1 {
                       fmt.Println(0)
                       return
                   }
                   forced[j] = 0
               }
               // neighbor constraint sum
               neigh := []int{}
               if i-1 >= 0 {
                   neigh = append(neigh, i-1)
               }
               if i+1 < n {
                   neigh = append(neigh, i+1)
               }
               if len(neigh) == 1 {
                   // single neighbor value equals d0
                   if d0 > 1 {
                       fmt.Println(0)
                       return
                   }
                   ni := neigh[0]
                   if ni&1 != p {
                       continue
                   }
                   j := idx[ni]
                   if forced[j] >= 0 && forced[j] != d0 {
                       fmt.Println(0)
                       return
                   }
                   forced[j] = d0
               } else if len(neigh) == 2 {
                   // constraint between two neighbors
                   // both neighbors have parity 1-p
                   ni0, ni1 := neigh[0], neigh[1]
                   if ni0&1 != (p^1) || ni1&1 != (p^1) {
                       continue
                   }
                   // map to indices in that parity
                   // we only store constraints for opposite parity group here
                   // skip for current p
                   // but process only when i&1 != p to avoid duplicate? already enforced by ni parity
                   // compute in opposite parity later
               }
           }
       }
       // We built forced and c only for parity p; but neighbor-sum constraints at positions with digits affect opposite parity only
       // To process constraints linking positions of parity p, we need to scan digits at i with i&1 != p
       // Build those constraints:
       for i := 0; i < n; i++ {
           ch := s[i]
           if !(ch >= '0' && ch <= '2') {
               continue
           }
           d0 := int(ch - '0')
           // only digits at opposite parity produce constraints here
           if i&1 == p {
               continue
           }
           // neighbors are in parity p
           neigh := []int{}
           if i-1 >= 0 {
               neigh = append(neigh, i-1)
           }
           if i+1 < n {
               neigh = append(neigh, i+1)
           }
           if len(neigh) == 1 {
               // single neighbor in parity p
               ni := neigh[0]
               if ni&1 != p {
                   continue
               }
               if d0 > 1 {
                   fmt.Println(0)
                   return
               }
               j := idx[ni]
               if forced[j] >= 0 && forced[j] != d0 {
                   fmt.Println(0)
                   return
               }
               forced[j] = d0
           } else if len(neigh) == 2 {
               ni0, ni1 := neigh[0], neigh[1]
               // both must be parity p
               if ni0&1 != p || ni1&1 != p {
                   continue
               }
               j0 := idx[ni0]
               j1 := idx[ni1]
               // ensure j1 == j0+1 or vice versa
               if j1 == j0+1 {
                   if c[j0] >= 0 && c[j0] != d0 {
                       fmt.Println(0)
                       return
                   }
                   c[j0] = d0
               } else if j0 == j1+1 {
                   if c[j1] >= 0 && c[j1] != d0 {
                       fmt.Println(0)
                       return
                   }
                   c[j1] = d0
               } else {
                   // neighbors not adjacent: invalid
                   fmt.Println(0)
                   return
               }
           }
       }
       // compute ways for this parity
       ways := int64(1)
       visited := make([]bool, m)
       for j := 0; j < m; j++ {
           if visited[j] {
               continue
           }
           // find component [l..r]
           l, r := j, j
           for l > 0 && c[l-1] != -1 {
               l--
           }
           for r < m-1 && c[r] != -1 {
               r++
           }
           for k := l; k <= r; k++ {
               visited[k] = true
           }
           size := r - l + 1
           if size == 1 {
               if forced[l] == -1 {
                   ways = ways * 2 % mod
               }
               continue
           }
           // path component
           cnt := 0
           for start := 0; start <= 1; start++ {
               if forced[l] >= 0 && forced[l] != start {
                   continue
               }
               ok := true
               x := start
               for k := l; k < r; k++ {
                   d := c[k]
                   // neighbor constraint must exist
                   // compute x2
                   x2 := d - x
                   if x2 < 0 || x2 > 1 {
                       ok = false
                       break
                   }
                   if forced[k+1] >= 0 && forced[k+1] != x2 {
                       ok = false
                       break
                   }
                   x = x2
               }
               if ok {
                   cnt++
               }
           }
           if cnt == 0 {
               fmt.Println(0)
               return
           }
           ways = ways * int64(cnt) % mod
       }
       total = total * ways % mod
   }
   fmt.Println(total)
}
