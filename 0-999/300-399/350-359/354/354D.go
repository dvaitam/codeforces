package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   // Row to sorted list of changed cols
   rowCols := make(map[int][]int, k)
   // existence map
   exist := make(map[int]map[int]bool, k)
   pts := make([][2]int, 0, k)
   for i := 0; i < k; i++ {
       var r, c int
       fmt.Fscan(in, &r, &c)
       if exist[r] == nil {
           exist[r] = make(map[int]bool)
       }
       exist[r][c] = true
       rowCols[r] = append(rowCols[r], c)
       pts = append(pts, [2]int{r, c})
   }
   // sort columns in each row
   rows := make([]int, 0, len(rowCols))
   for r, cols := range rowCols {
       sort.Ints(cols)
       rowCols[r] = cols
       rows = append(rows, r)
   }
   sort.Ints(rows)
   // compute dp: max layers count (h+1)
   dpRow := make(map[int]map[int]int, len(rows))
   // iterate rows descending
   for i := len(rows) - 1; i >= 0; i-- {
       r := rows[i]
       dpRow[r] = make(map[int]int, len(rowCols[r]))
       next := dpRow[r+1]
       for _, c := range rowCols[r] {
           var d1, d2 int
           if next != nil {
               d1 = next[c]
               d2 = next[c+1]
           }
           m := d1
           if d2 < m {
               m = d2
           }
           dpRow[r][c] = m + 1
       }
   }
   // build candidates
   type cand struct{ dp, r, c int }
   cands := make([]cand, 0, k)
   for _, p := range pts {
       r, c := p[0], p[1]
       d := dpRow[r][c]
       if d > 0 {
           cands = append(cands, cand{d, r, c})
       }
   }
   sort.Slice(cands, func(i, j int) bool {
       return cands[i].dp > cands[j].dp
   })
   // covered map
   covered := make(map[int]map[int]bool, len(rowCols))
   for r := range rowCols {
       covered[r] = make(map[int]bool, len(rowCols[r]))
   }
   rem := k
   var cost int64 = 0
   // process candidates
   for _, cd := range cands {
       if rem <= 0 {
           break
       }
       r, c, d := cd.r, cd.c, cd.dp
       if covered[r][c] {
           continue
       }
       // use triangle of height h = d-1
       h := d - 1
       // size of triangle: T(h) = (h+1)*(h+2)/2 = d*(d+1)/2
       tsize := d * (d + 1) / 2
       cost += int64(tsize + 2)
       // mark covered
       for p := 0; p <= h; p++ {
           rr := r + p
           cols := rowCols[rr]
           if cols == nil {
               continue
           }
           // find first index >= c
           l := sort.Search(len(cols), func(i int) bool { return cols[i] >= c })
           // iterate until cols[i] <= c+p
           for i := l; i < len(cols) && cols[i] <= c+p; i++ {
               cc := cols[i]
               if !covered[rr][cc] {
                   covered[rr][cc] = true
                   rem--
               }
           }
       }
   }
   // any remaining singletons
   if rem > 0 {
       // each costs 3
       cost += int64(rem * 3)
   }
   // output cost
   w := bufio.NewWriter(os.Stdout)
   fmt.Fprint(w, cost)
   w.Flush()
}
