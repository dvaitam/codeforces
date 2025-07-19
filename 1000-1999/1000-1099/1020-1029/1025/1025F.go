package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Pct struct {
   x, y int
}

type Lin struct {
   x, y int
   p    [2]int
}

func C2(x int64) int64 {
   if x < 2 {
       return 0
   }
   return x * (x - 1) / 2
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   v := make([]Pct, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &v[i].x, &v[i].y)
   }
   sort.Slice(v, func(i, j int) bool {
       if v[i].x != v[j].x {
           return v[i].x < v[j].x
       }
       return v[i].y < v[j].y
   })

   var ans int64
   wlin := make([]Lin, 0, n*(n-1)/2)
   for i := 0; i < n; i++ {
       for j := i + 1; j < n; j++ {
           if v[i].x == v[j].x {
               ans += C2(int64(i)) * C2(int64(n-1-j))
           } else {
               dx := v[i].y - v[j].y
               dy := v[i].x - v[j].x
               if dy < 0 {
                   dy = -dy
                   dx = -dx
               }
               wlin = append(wlin, Lin{x: i, y: j, p: [2]int{dx, dy}})
           }
       }
   }

   pos := make([]int, n)
   ord := make([]int, n)
   for i := 0; i < n; i++ {
       pos[i] = i
   }
   panta := [2]int{-1000000001, 1}
   sort.Slice(pos, func(i, j int) bool {
       a, b := pos[i], pos[j]
       return int64(v[a].y-v[b].y)*int64(panta[1]) < int64(panta[0])*int64(v[a].x-v[b].x)
   })
   for i, pi := range pos {
       ord[pi] = i
   }

   sort.Slice(wlin, func(i, j int) bool {
       a, b := wlin[i], wlin[j]
       return int64(a.p[0])*int64(b.p[1]) < int64(a.p[1])*int64(b.p[0])
   })
   for _, l := range wlin {
       // swap positions of neighbors in order
       if l.x < 0 || l.x >= n || l.y < 0 || l.y >= n {
           continue
       }
       // assume ord[l.x] and ord[l.y] differ by 1
       ord[l.x], ord[l.y] = ord[l.y], ord[l.x]
       mx := ord[l.x]
       if ord[l.y] > mx {
           mx = ord[l.y]
       }
       // translate to 1-based logic: C2(mx-1) * C2(n-1-mx)
       ans += C2(int64(mx-1)) * C2(int64(n-1-mx))
   }
   fmt.Fprintln(writer, ans)
}
