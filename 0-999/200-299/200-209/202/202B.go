package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   lesha := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &lesha[i])
   }
   var m int
   fmt.Fscan(reader, &m)
   bestP := -1
   bestIdx := -1
   totalInv := n * (n - 1) / 2
   // generate all permutations of 0..n-1
   perm := make([]int, n)
   used := make([]bool, n)
   var perms [][]int
   var dfs func(int)
   dfs = func(pos int) {
       if pos == n {
           cp := make([]int, n)
           copy(cp, perm)
           perms = append(perms, cp)
           return
       }
       for i := 0; i < n; i++ {
           if !used[i] {
               used[i] = true
               perm[pos] = i
               dfs(pos + 1)
               used[i] = false
           }
       }
   }
   dfs(0)
   // process archive problems
   for idx := 1; idx <= m; idx++ {
       var k int
       fmt.Fscan(reader, &k)
       archive := make([]string, k)
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &archive[i])
       }
       // for this problem, find minimal inversions among perms that are subsequences
       minInv := totalInv + 1
       for _, p := range perms {
           // check subsequence
           pos := 0
           for _, w := range archive {
               if w == lesha[p[pos]] {
                   pos++
                   if pos == n {
                       break
                   }
               }
           }
           if pos == n {
               // compute inversions x
               x := 0
               for i := 0; i < n; i++ {
                   for j := i + 1; j < n; j++ {
                       if p[i] > p[j] {
                           x++
                       }
                   }
               }
               if x < minInv {
                   minInv = x
               }
           }
       }
       if minInv <= totalInv {
           p := totalInv - minInv + 1
           if p > bestP {
               bestP = p
               bestIdx = idx
           }
       }
   }
   if bestIdx == -1 {
       fmt.Println("Brand new problem!")
   } else {
       fmt.Println(bestIdx)
       // build similarity bar
       bar := strings.Repeat("|", bestP)
       fmt.Println(":" + bar + ":")
   }
}
