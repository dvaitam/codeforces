package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

// Aho-Corasick node
type node struct {
   next   []int
   fail   int
   output int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, k int
   fmt.Fscan(in, &n, &m, &k)
   readNum := func() []int {
       var l int
       fmt.Fscan(in, &l)
       a := make([]int, l)
       for i := 0; i < l; i++ {
           fmt.Fscan(in, &a[i])
       }
       return a
   }
   lnum := readNum()
   rnum := readNum()
   patterns := make([][]int, n)
   vals := make([]int, n)
   for i := 0; i < n; i++ {
       s := readNum()
       patterns[i] = s
       fmt.Fscan(in, &vals[i])
   }
   // build AC
   // initial node
   ac := []node{{next: make([]int, m), fail: 0, output: 0}}
   // insert patterns
   for i, pat := range patterns {
       cur := 0
       for _, d := range pat {
           if ac[cur].next[d] == 0 {
               ac[cur].next[d] = len(ac)
               ac = append(ac, node{next: make([]int, m)})
           }
           cur = ac[cur].next[d]
       }
       ac[cur].output += vals[i]
   }
   // build failure
   q := make([]int, 0, len(ac))
   // first level
   for d := 0; d < m; d++ {
       v := ac[0].next[d]
       if v != 0 {
           ac[v].fail = 0
           q = append(q, v)
       }
   }
   for i := 0; i < len(q); i++ {
       u := q[i]
       for d := 0; d < m; d++ {
           v := ac[u].next[d]
           if v != 0 {
               f := ac[u].fail
               for f != 0 && ac[f].next[d] == 0 {
                   f = ac[f].fail
               }
               if ac[f].next[d] != 0 {
                   f = ac[f].next[d]
               }
               ac[v].fail = f
               ac[v].output += ac[f].output
               q = append(q, v)
           } else {
               ac[u].next[d] = ac[ac[u].fail].next[d]
           }
       }
   }
   // precompute go and weight
   S := len(ac)
   goTo := make([][]int, S)
   gain := make([][]int, S)
   for i := 0; i < S; i++ {
       goTo[i] = make([]int, m)
       gain[i] = make([]int, m)
       for d := 0; d < m; d++ {
           nxt := ac[i].next[d]
           goTo[i][d] = nxt
           gain[i][d] = ac[nxt].output
       }
   }
   // solve for bound
   // solve calculates number of valid integers <= bound
   solve := func(bound []int) int {
       L := len(bound)
       // dimensions
       sumDim := k + 1
       stateDim := S
       tightDim := 2
       startDim := 2
       total := tightDim * stateDim * sumDim * startDim
       dp := make([]int, total)
       dpn := make([]int, total)
       // index helper
       idx := func(t, s, su, st int) int {
           return ((t*stateDim + s)*sumDim + su)*startDim + st
       }
       dp[idx(1, 0, 0, 0)] = 1
       for pos := 0; pos < L; pos++ {
           // clear next
           for i := range dpn {
               dpn[i] = 0
           }
           for t := 0; t < 2; t++ {
               for s := 0; s < stateDim; s++ {
                   for su := 0; su <= k; su++ {
                       for st := 0; st < 2; st++ {
                           v := dp[idx(t, s, su, st)]
                           if v == 0 {
                               continue
                           }
                           maxd := m - 1
                           if t == 1 {
                               maxd = bound[pos]
                           }
                           for d := 0; d <= maxd; d++ {
                               nt := t
                               if t == 1 && d < maxd {
                                   nt = 0
                               }
                               var ns, nsu, nst int
                               if st == 0 && d == 0 {
                                   ns, nsu, nst = 0, 0, 0
                               } else {
                                   nst = 1
                                   cs := s
                                   if st == 0 {
                                       cs = 0
                                   }
                                   ns = goTo[cs][d]
                                   nsu = su + gain[cs][d]
                                   if nsu > k {
                                       continue
                                   }
                               }
                               j := idx(nt, ns, nsu, nst)
                               dpn[j] = (dpn[j] + v) % MOD
                           }
                       }
                   }
               }
           }
           dp, dpn = dpn, dp
       }
       res := 0
       for t := 0; t < 2; t++ {
           for s := 0; s < stateDim; s++ {
               for su := 0; su <= k; su++ {
                   res = (res + dp[idx(t, s, su, 1)]) % MOD
               }
           }
       }
       return res
   }
   // subtract one from lnum
   dec := func(a []int) []int {
       b := make([]int, len(a))
       copy(b, a)
       i := len(b) - 1
       for i >= 0 {
           if b[i] > 0 {
               b[i]--
               break
           }
           b[i] = m - 1
           i--
       }
       // remove leading zeros
       if len(b) > 1 && b[0] == 0 {
           j := 0
           for j < len(b) && b[j] == 0 {
               j++
           }
           if j == len(b) {
               return []int{0}
           }
           b = b[j:]
       }
       return b
   }
   ldec := dec(lnum)
   ans := solve(rnum) - solve(ldec)
   ans = (ans%MOD + MOD) % MOD
   fmt.Println(ans)
}
