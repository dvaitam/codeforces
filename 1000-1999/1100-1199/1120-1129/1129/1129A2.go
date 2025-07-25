package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m int
   if _, err := fmt.Fscan(in, &n, &m); err != nil {
       return
   }
   // distances of candies grouped by start station
   dists := make([][]int, n)
   for i := 0; i < m; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       d := b - a
       if d < 0 {
           d += n
       }
       dists[a] = append(dists[a], d)
   }
   // precompute max extra time for each station
   mx := make([]int64, n)
   stations := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if len(dists[i]) == 0 {
           continue
       }
       stations = append(stations, i)
       ds := dists[i]
       // sort descending
       sort.Slice(ds, func(x, y int) bool { return ds[x] > ds[y] })
       var best int64
       for j, dd := range ds {
           t := int64(j)*int64(n) + int64(dd)
           if t > best {
               best = t
           }
       }
       mx[i] = best
   }
   // compute answer for each starting station s
   res := make([]int64, n)
   for s := 0; s < n; s++ {
       var ans int64
       for _, i := range stations {
           // distance from s to i
           di := i - s
           if di < 0 {
               di += n
           }
           t := int64(di) + mx[i]
           if t > ans {
               ans = t
           }
       }
       res[s] = ans
   }
   // output
   for i, v := range res {
       if i > 0 {
           fmt.Fprint(out, " ")
       }
       fmt.Fprint(out, v)
   }
   fmt.Fprintln(out)
}
