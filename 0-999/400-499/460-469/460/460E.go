package main

import (
   "fmt"
   "sort"
)

type point struct{ x, y int }

var (
   n, r     int
   v        []point
   ans, aaa []point
   maxx     int64 = -1 << 60
)

func dfs(k, d int) {
   if d == n {
       var sum int64
       for i := 0; i < len(ans); i++ {
           for j := i + 1; j < len(ans); j++ {
               dx := int64(ans[i].x - ans[j].x)
               dy := int64(ans[i].y - ans[j].y)
               sum += dx*dx + dy*dy
           }
       }
       if sum > maxx {
           maxx = sum
           // copy ans to aaa
           aaa = make([]point, len(ans))
           copy(aaa, ans)
       }
       return
   }
   // limit iterations: i < len(v) && i < 30-2*n
   limit := len(v)
   cut := 30 - 2*n
   if cut < limit {
       limit = cut
   }
   for i := k; i < limit; i++ {
       ans = append(ans, v[i])
       dfs(i, d+1)
       ans = ans[:len(ans)-1]
   }
}

func main() {
   if _, err := fmt.Scan(&n, &r); err != nil {
       return
   }
   // generate points on circle ring [r-1, r]
   for i := -r; i <= r; i++ {
       for j := -r; j <= r; j++ {
           d2 := i*i + j*j
           if d2 <= r*r && d2 >= (r-1)*(r-1) {
               v = append(v, point{i, j})
           }
       }
   }
   // sort by descending squared radius
   sort.Slice(v, func(i, j int) bool {
       di := v[i].x*v[i].x + v[i].y*v[i].y
       dj := v[j].x*v[j].x + v[j].y*v[j].y
       return di > dj
   })
   dfs(0, 0)
   // output
   fmt.Println(maxx)
   for _, p := range aaa {
       fmt.Println(p.x, p.y)
   }
}
