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
   var n, k int
   fmt.Fscan(reader, &n, &k)
   a := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
   }
   // order positions by descending skill
   ord := make([]int, n)
   for i := 0; i < n; i++ {
       ord[i] = i
   }
   sort.Slice(ord, func(i, j int) bool {
       return a[ord[i]] > a[ord[j]]
   })
   // doubly linked list of positions
   prev := make([]int, n)
   next := make([]int, n)
   for i := 0; i < n; i++ {
       prev[i] = i - 1
       next[i] = i + 1
   }
   next[n-1] = -1
   removed := make([]bool, n)
   ans := make([]byte, n)
   team := byte('1')
   cur := 0
   // assign students to teams
   for cur < n {
       center := ord[cur]
       cur++
       if removed[center] {
           continue
       }
       // collect up to k neighbors on each side
       ls := make([]int, 0, k)
       rs := make([]int, 0, k)
       x := prev[center]
       for j := 0; j < k && x != -1; j++ {
           ls = append(ls, x)
           x = prev[x]
       }
       x = next[center]
       for j := 0; j < k && x != -1; j++ {
           rs = append(rs, x)
           x = next[x]
       }
       // determine block boundaries
       L, R := center, center
       if len(ls) > 0 {
           L = ls[len(ls)-1]
       }
       if len(rs) > 0 {
           R = rs[len(rs)-1]
       }
       // mark removed and assign team
       removed[center] = true
       ans[center] = team
       for _, p := range ls {
           removed[p] = true
           ans[p] = team
       }
       for _, p := range rs {
           removed[p] = true
           ans[p] = team
       }
       // unlink block [L..R]
       leftNeighbor := prev[L]
       rightNeighbor := next[R]
       if leftNeighbor != -1 {
           next[leftNeighbor] = rightNeighbor
       }
       if rightNeighbor != -1 {
           prev[rightNeighbor] = leftNeighbor
       }
       // alternate team
       if team == '1' {
           team = '2'
       } else {
           team = '1'
       }
   }
   // output result
   writer.Write(ans)
   writer.WriteByte('\n')
}
