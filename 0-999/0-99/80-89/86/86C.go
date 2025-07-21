package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000009

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   patterns := make([]string, m)
   maxLen := 0
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &patterns[i])
       if len(patterns[i]) > maxLen {
           maxLen = len(patterns[i])
       }
   }
   // Build Aho-Corasick
   // Alphabet: A,C,G,T -> 0,1,2,3
   toIndex := func(c byte) int {
       switch c {
       case 'A': return 0
       case 'C': return 1
       case 'G': return 2
       case 'T': return 3
       }
       return 0
   }
   // trie
   next := make([][4]int, 1)
   out := make([][]int, 1)
   link := []int{0}
   for _, p := range patterns {
       v := 0
       for i := 0; i < len(p); i++ {
           ci := toIndex(p[i])
           if next[v][ci] == 0 {
               next = append(next, [4]int{})
               out = append(out, nil)
               link = append(link, 0)
               next[v][ci] = len(next) - 1
           }
           v = next[v][ci]
       }
       out[v] = append(out[v], len(p))
   }
   // build links
   queue := make([]int, 0, len(next))
   // init depth 1
   for c := 0; c < 4; c++ {
       if next[0][c] != 0 {
           queue = append(queue, next[0][c])
           link[next[0][c]] = 0
       }
   }
   for qi := 0; qi < len(queue); qi++ {
       v := queue[qi]
       // merge outputs
       lv := link[v]
       if len(out[lv]) > 0 {
           out[v] = append(out[v], out[lv]...)
       }
       for c := 0; c < 4; c++ {
           u := next[v][c]
           if u != 0 {
               link[u] = next[link[v]][c]
               queue = append(queue, u)
           } else {
               next[v][c] = next[link[v]][c]
           }
       }
   }
   // DP arrays: dpCur[k][state]
   numStates := len(next)
   M := maxLen
   dpCur := make([][]int, M)
   dpNext := make([][]int, M)
   for k := 0; k < M; k++ {
       dpCur[k] = make([]int, numStates)
       dpNext[k] = make([]int, numStates)
   }
   dpCur[0][0] = 1
   // iterate positions
   for pos := 0; pos < n; pos++ {
       // clear dpNext
       for k := 0; k < M; k++ {
           for s := 0; s < numStates; s++ {
               dpNext[k][s] = 0
           }
       }
       for k := 0; k < M; k++ {
           for s := 0; s < numStates; s++ {
               ways := dpCur[k][s]
               if ways == 0 {
                   continue
               }
               for c := 0; c < 4; c++ {
                   ns := next[s][c]
                   // find longest pattern ending here
                   maxP := 0
                   for _, L := range out[ns] {
                       if L > maxP {
                           maxP = L
                       }
                   }
                   t := k + 1
                   if maxP > t {
                       maxP = t
                   }
                   k2 := t - maxP
                   if k2 < M {
                       dpNext[k2][ns] = (dpNext[k2][ns] + ways) % MOD
                   }
               }
           }
       }
       // swap
       dpCur, dpNext = dpNext, dpCur
   }
   // sum dpCur[0][*]
   ans := 0
   for s := 0; s < len(dpCur[0]); s++ {
       ans = (ans + dpCur[0][s]) % MOD
   }
   fmt.Println(ans)
}
