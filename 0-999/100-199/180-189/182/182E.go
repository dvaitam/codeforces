package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, l int
   if _, err := fmt.Fscan(in, &n, &l); err != nil {
       return
   }
   a := make([]int, n)
   b := make([]int, n)
   maxLen := 0
   maxW := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i])
       if a[i] > maxLen {
           maxLen = a[i]
       }
       if b[i] > maxLen {
           maxLen = b[i]
       }
       if a[i] > maxW {
           maxW = a[i]
       }
       if b[i] > maxW {
           maxW = b[i]
       }
   }
   modLen := maxLen + 1
   // oriented boards: for each orientation (length, width, type)
   type entry struct{ L, W, k int }
   var oriented []entry
   for i := 0; i < n; i++ {
       oriented = append(oriented, entry{a[i], b[i], i})
       if a[i] != b[i] {
           oriented = append(oriented, entry{b[i], a[i], i})
       }
   }
   // dp[s_mod][width][type] and g[s_mod][width]
   dp := make([][][]int, modLen)
   g := make([][]int, modLen)
   for i := 0; i < modLen; i++ {
       dp[i] = make([][]int, maxW+1)
       g[i] = make([]int, maxW+1)
       for w := 0; w <= maxW; w++ {
           dp[i][w] = make([]int, n)
       }
   }
   // DP over total length s
   for s := 1; s <= l; s++ {
       cur := s % modLen
       // clear current
       for w := 0; w <= maxW; w++ {
           g[cur][w] = 0
           for k := 0; k < n; k++ {
               dp[cur][w][k] = 0
           }
       }
       // transitions and initial
       for _, e := range oriented {
           L, W, k := e.L, e.W, e.k
           if L > s {
               continue
           }
           var val int
           if s == L {
               val = 1
           } else {
               prev := (s - L) % modLen
               // total ways with prev width == L, excluding same type
               val = g[prev][L] - dp[prev][L][k]
               if val < 0 {
                   val += mod
               }
           }
           if val != 0 {
               dp[cur][W][k] += val
               if dp[cur][W][k] >= mod {
                   dp[cur][W][k] -= mod
               }
               g[cur][W] += val
               if g[cur][W] >= mod {
                   g[cur][W] -= mod
               }
           }
       }
   }
   // sum up answers at length l
   res := 0
   last := l % modLen
   for w := 0; w <= maxW; w++ {
       res += g[last][w]
       if res >= mod {
           res -= mod
       }
   }
   fmt.Println(res)
}
