package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var g [5][5]int64
   for i := 0; i < 5; i++ {
       for j := 0; j < 5; j++ {
           fmt.Fscan(reader, &g[i][j])
       }
   }
   perm := []int{0, 1, 2, 3, 4}
   var maxH int64
   // generate all permutations
   var dfs func(int)
   dfs = func(idx int) {
       if idx == 5 {
           // compute happiness for this ordering
           p := perm
           var sum int64
           // stage 1
           sum += g[p[0]][p[1]] + g[p[1]][p[0]]
           sum += g[p[2]][p[3]] + g[p[3]][p[2]]
           // stage 2
           sum += g[p[1]][p[2]] + g[p[2]][p[1]]
           sum += g[p[3]][p[4]] + g[p[4]][p[3]]
           // stage 3
           sum += g[p[2]][p[3]] + g[p[3]][p[2]]
           // stage 4
           sum += g[p[3]][p[4]] + g[p[4]][p[3]]
           if sum > maxH {
               maxH = sum
           }
           return
       }
       for i := idx; i < 5; i++ {
           perm[idx], perm[i] = perm[i], perm[idx]
           dfs(idx + 1)
           perm[idx], perm[i] = perm[i], perm[idx]
       }
   }
   dfs(0)
   fmt.Println(maxH)
}
