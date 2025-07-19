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
   var s string
   fmt.Fscan(reader, &s)
   n := len(s)
   nx := make([][26]int, n+1)
   pr := make([][26]int, n+1)
   gl := make([][26]int, n+1)
   gr := make([][26]int, n+1)
   sm := make([]int, n+1)
   sm2 := make([]int, 26)
   // initialize next occurrences and previous occurrences
   for j := 0; j < 26; j++ {
       nx[n][j] = n
       pr[0][j] = -1
   }
   for i := n - 1; i >= 0; i-- {
       for j := 0; j < 26; j++ {
           nx[i][j] = nx[i+1][j]
       }
       nx[i][s[i]-'a'] = i
   }
   for i := 0; i < n; i++ {
       for j := 0; j < 26; j++ {
           pr[i+1][j] = pr[i][j]
       }
       pr[i+1][s[i]-'a'] = i
   }
   var tmp [30]bool
   type pair struct{ pos, c int }
   var vv []pair
   // compute gl and sm via backward DP
   for i := n - 1; i >= 0; i-- {
       sm[i] = sm2[s[i]-'a']
       vv = vv[:0]
       for j := 0; j < 26; j++ {
           if nx[i][j] < n {
               vv = append(vv, pair{nx[i][j], j})
           }
       }
       sort.Slice(vv, func(a, b int) bool { return vv[a].pos < vv[b].pos })
       for _, pc := range vv {
           j := pc.c
           rpos := pc.pos
           // mex of reachable states
           for k := range tmp { tmp[k] = false }
           for g := 0; g < 26; g++ {
               if nx[i][g] < rpos {
                   lst := pr[rpos][g]
                   goVal := gl[i][g] ^ sm[nx[i][g]] ^ sm[lst] ^ gl[lst+1][j]
                   tmp[goVal] = true
               }
           }
           mex := 0
           for tmp[mex] {
               mex++
           }
           gl[i][j] = mex
       }
       if i > 0 {
           sm2[s[i-1]-'a'] ^= gl[i][s[i-1]-'a']
       }
   }
   // compute gr via forward DP
   for i := 1; i <= n; i++ {
       vv = vv[:0]
       for j := 0; j < 26; j++ {
           if pr[i][j] != -1 {
               vv = append(vv, pair{pr[i][j], j})
           }
       }
       sort.Slice(vv, func(a, b int) bool { return vv[a].pos > vv[b].pos })
       for _, pc := range vv {
           j := pc.c
           lpos := pc.pos + 1
           for k := range tmp { tmp[k] = false }
           for g := 0; g < 26; g++ {
               if pr[i][g] >= lpos {
                   lst := nx[lpos][g]
                   goVal := gr[i][g] ^ sm[pr[i][g]] ^ sm[lst] ^ gr[lst][j]
                   tmp[goVal] = true
               }
           }
           mex := 0
           for tmp[mex] {
               mex++
           }
           gr[i][j] = mex
       }
   }
   var m int
   fmt.Fscan(reader, &m)
   for q := 0; q < m; q++ {
       var l, r int
       fmt.Fscan(reader, &l, &r)
       l--
       for k := range tmp { tmp[k] = false }
       for j := 0; j < 26; j++ {
           if nx[l][j] < r {
               lst := pr[r][j]
               goVal := sm[nx[l][j]] ^ sm[lst] ^ gl[l][j] ^ gr[r][j]
               tmp[goVal] = true
           }
       }
       mex := 0
       for tmp[mex] {
           mex++
       }
       if mex != 0 {
           fmt.Fprintln(writer, "Alice")
       } else {
           fmt.Fprintln(writer, "Bob")
       }
   }
}
