package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type interval struct{ l, r int }

// merge sorted intervals
func mergeIntervals(a []interval) []interval {
   if len(a) == 0 {
       return a
   }
   sort.Slice(a, func(i, j int) bool { return a[i].l < a[j].l })
   res := make([]interval, 0, len(a))
   cur := a[0]
   for _, iv := range a[1:] {
       if iv.l <= cur.r+1 {
           if iv.r > cur.r {
               cur.r = iv.r
           }
       } else {
           res = append(res, cur)
           cur = iv
       }
   }
   res = append(res, cur)
   return res
}

// process a row with obstacles: subtract obstacles and expand right
func processRow(curr []interval, ys []int, n int) []interval {
   sort.Ints(ys)
   k := len(ys)
   // subtract obstacles
   initSeg := make([]interval, 0)
   j := 0
   for _, iv := range curr {
       l, r := iv.l, iv.r
       // skip obstacles before l
       for j < k && ys[j] < l {
           j++
       }
       start := l
       for j < k && ys[j] <= r {
           if start <= ys[j]-1 {
               initSeg = append(initSeg, interval{start, ys[j] - 1})
           }
           start = ys[j] + 1
           j++
       }
       if start <= r {
           initSeg = append(initSeg, interval{start, r})
       }
   }
   if len(initSeg) == 0 {
       return nil
   }
   // expand to right
   newSeg := make([]interval, 0, len(initSeg))
   for _, iv := range initSeg {
       l, r := iv.l, iv.r
       // find first obstacle > r
       idx := sort.Search(k, func(i int) bool { return ys[i] > r })
       end := n
       if idx < k {
           end = ys[idx] - 1
       }
       newSeg = append(newSeg, interval{l, end})
   }
   // merge and return
   return mergeIntervals(newSeg)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // edge n==1
   if n == 1 {
       fmt.Fprint(out, 0)
       return
   }
   obs := make(map[int][]int)
   for i := 0; i < m; i++ {
       var x, y int
       fmt.Fscan(in, &x, &y)
       obs[x] = append(obs[x], y)
   }
   // prepare rows with obstacles
   rows := make([]int, 0, len(obs))
   for r := range obs {
       rows = append(rows, r)
   }
   sort.Ints(rows)
   // initial at row 1
   curr := []interval{{1, 1}}
   // process row1
   if ys, ok := obs[1]; ok {
       curr = processRow(curr, ys, n)
   } else {
       // expand to right until n
       curr = []interval{{1, n}}
   }
   if len(curr) == 0 {
       fmt.Fprint(out, -1)
       return
   }
   currRow := 1
   // process other rows
   for _, r := range rows {
       if r == 1 {
           continue
       }
       if len(curr) == 0 {
           break
       }
       // move down to r
       if currRow+1 < r {
           // expansion in empty rows
           // collapse to [min_l, n]
           minl := curr[0].l
           for _, iv := range curr {
               if iv.l < minl {
                   minl = iv.l
               }
           }
           curr = []interval{{minl, n}}
       }
       currRow = r
       // process obstacles at row r
       curr = processRow(curr, obs[r], n)
   }
   if len(curr) == 0 {
       fmt.Fprint(out, -1)
       return
   }
   // if last obstacle row is the final row, ensure column n is reachable there
   if len(rows) > 0 && rows[len(rows)-1] == n {
       ok := false
       for _, iv := range curr {
           if iv.l <= n && n <= iv.r {
               ok = true
               break
           }
       }
       if !ok {
           fmt.Fprint(out, -1)
           return
       }
   }
   // reachable, path length fixed
   ans := int64(2)*int64(n) - 2
   fmt.Fprint(out, ans)
}
