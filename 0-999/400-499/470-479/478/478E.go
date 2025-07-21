package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var wavies []int64

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int64
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   generateWavies()
   sort.Slice(wavies, func(i, j int) bool { return wavies[i] < wavies[j] })
   var cnt int64
   for _, v := range wavies {
       if v%n == 0 {
           cnt++
           if cnt == k {
               fmt.Println(v)
               return
           }
       }
   }
   fmt.Println(-1)
}

// generateWavies populates the global wavies slice with all wavy numbers <= 10^14
func generateWavies() {
   const maxLen = 14
   // length 1
   for d := 1; d <= 9; d++ {
       wavies = append(wavies, int64(d))
   }
   // length 2 seeds
   type seed struct{ num int64; p1, p2 int }
   var seeds []seed
   for d1 := 1; d1 <= 9; d1++ {
       for d2 := 0; d2 <= 9; d2++ {
           num := int64(d1*10 + d2)
           wavies = append(wavies, num)
           seeds = append(seeds, seed{num, d1, d2})
       }
   }
   // build longer numbers
   for _, s := range seeds {
       dfs(s.num, s.p1, s.p2, 2, maxLen)
   }
}

// dfs builds wavy numbers by extending current number
func dfs(num int64, p1, p2, length, maxLen int) {
   if length >= maxLen {
       return
   }
   for d := 0; d <= 9; d++ {
       // p2 must be a peak or valley relative to p1 and d
       if (p2 > p1 && p2 > d) || (p2 < p1 && p2 < d) {
           newNum := num*10 + int64(d)
           wavies = append(wavies, newNum)
           dfs(newNum, p2, d, length+1, maxLen)
       }
   }
}
