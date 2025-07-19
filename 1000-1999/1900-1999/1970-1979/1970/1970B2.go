package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type pair struct { a, idx int }

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   v := make([]pair, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &v[i].a)
       v[i].idx = i
   }
   sort.Slice(v, func(i, j int) bool { return v[i].a < v[j].a })

   ansx := make([]int, n)
   ansy := make([]int, n)
   dir := make([]int, n)
   taken := make([]int, n+1)
   takenID := make([]int, n+1)

   if v[0].a == 0 {
       id0 := v[0].idx
       ansx[id0] = 1
       ansy[id0] = 1
       dir[id0] = id0 + 1
       taken[1] = 1
       takenID[1] = id0
   } else {
       i := 0
       for i+1 < n && v[i].a != v[i+1].a {
           i++
       }
       dis := v[i].a
       id1 := v[i].idx
       id2 := v[i+1].idx
       // Assign first duplicate
       ansx[id1] = 1
       ansy[id1] = 1
       dir[id1] = id2 + 1
       taken[1] = 1
       takenID[1] = id1

       x := 1 + dis
       y := 1
       if x > n {
           y += x - n
           x = n
       }
       ansx[id2] = x
       ansy[id2] = y
       dir[id2] = id1 + 1
       taken[x] = y
       takenID[x] = id2
   }

   curX := 1
   for _, p := range v {
       a := p.a
       id := p.idx
       if dir[id] != 0 {
           continue
       }
       for curX <= n && taken[curX] != 0 {
           curX++
       }
       var y int
       if a == 0 {
           y = 1
           dir[id] = id + 1
       } else {
           if curX-a >= 1 {
               y = taken[curX-a]
               dir[id] = takenID[curX-a] + 1
           } else {
               y = a - curX + 2
               dir[id] = takenID[1] + 1
           }
       }
       ansx[id] = curX
       ansy[id] = y
       taken[curX] = y
       takenID[curX] = id
   }

   fmt.Fprintln(out, "YES")
   for i := 0; i < n; i++ {
       fmt.Fprintln(out, ansx[i], ansy[i])
   }
   for i := 0; i < n; i++ {
       fmt.Fprintln(out, dir[i])
   }
}
