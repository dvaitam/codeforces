package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   total := int64(0)
   a := make([]int64, n*n)
   for i := range a {
       fmt.Fscan(reader, &a[i])
       total += a[i]
   }
   // magic sum
   s := total / int64(n)
   if n == 1 {
       // trivial
       fmt.Println(a[0])
       fmt.Printf("%d\n", a[0])
       return
   }
   // count distinct values
   counts := make(map[int64]int)
   for _, v := range a {
       counts[v]++
   }
   uniq := make([]int64, 0, len(counts))
   for v := range counts {
       uniq = append(uniq, v)
   }
   // prepare parallel count slice
   cnts := make([]int, len(uniq))
   for i, v := range uniq {
       cnts[i] = counts[v]
   }
   // grid and sums
   grid := make([][]int64, n)
   for i := range grid {
       grid[i] = make([]int64, n)
   }
   rowSum := make([]int64, n)
   colSum := make([]int64, n)
   diag := make([]int64, 2)
   var found bool
   var dfs func(pos int)
   dfs = func(pos int) {
       if found {
           return
       }
       if pos == n*n {
           // solution
           fmt.Println(s)
           w := bufio.NewWriter(os.Stdout)
           for i := 0; i < n; i++ {
               for j := 0; j < n; j++ {
                   if j > 0 {
                       w.WriteByte(' ')
                   }
                   fmt.Fprint(w, grid[i][j])
               }
               w.WriteByte('\n')
           }
           w.Flush()
           found = true
           return
       }
       r := pos / n
       c := pos % n
       for i, v := range uniq {
           if cnts[i] == 0 {
               continue
           }
           // try place v at (r,c)
           nsr := rowSum[r] + v
           nsc := colSum[c] + v
           if nsr > s || nsc > s {
               continue
           }
           // diag checks
           needD0 := (r == c)
           needD1 := (r+c == n-1)
           if needD0 && diag[0]+v > s {
               continue
           }
           if needD1 && diag[1]+v > s {
               continue
           }
           // if end of row, must equal
           if c == n-1 && nsr != s {
               continue
           }
           // if end of column, but only if last row
           if r == n-1 && nsc != s {
               continue
           }
           // place
           grid[r][c] = v
           cnts[i]--
           rowSum[r] = nsr
           colSum[c] = nsc
           if needD0 {
               diag[0] += v
           }
           if needD1 {
               diag[1] += v
           }
           dfs(pos + 1)
           // undo
           if found {
               return
           }
           if needD0 {
               diag[0] -= v
           }
           if needD1 {
               diag[1] -= v
           }
           rowSum[r] -= v
           colSum[c] -= v
           cnts[i]++
       }
   }
   dfs(0)
}
