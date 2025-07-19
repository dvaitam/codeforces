package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type point struct {
   x, y, z, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   pts := make([]point, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &pts[i].x, &pts[i].y, &pts[i].z)
       pts[i].id = i + 1
   }
   sort.Slice(pts, func(i, j int) bool {
       if pts[i].x != pts[j].x {
           return pts[i].x < pts[j].x
       }
       if pts[i].y != pts[j].y {
           return pts[i].y < pts[j].y
       }
       return pts[i].z < pts[j].z
   })
   vis := make([]bool, n)

   // First pass: same x and y
   l := 0
   for i := 1; i < n; i++ {
       if l >= 0 && pts[i].x == pts[l].x && pts[i].y == pts[l].y {
           fmt.Fprintln(writer, pts[i].id, pts[l].id)
           vis[i], vis[l] = true, true
           l = -1
       } else {
           l = i
       }
   }
   // Second pass: same x
   l = -1
   for i := 0; i < n; i++ {
       if !vis[i] {
           l = i
           break
       }
   }
   for i := l + 1; i < n; i++ {
       if vis[i] {
           continue
       }
       if l >= 0 && pts[i].x == pts[l].x {
           fmt.Fprintln(writer, pts[i].id, pts[l].id)
           vis[i], vis[l] = true, true
           l = -1
       } else {
           l = i
       }
   }
   // Third pass: arbitrary
   l = -1
   for i := 0; i < n; i++ {
       if !vis[i] {
           l = i
           break
       }
   }
   for i := l + 1; i < n; i++ {
       if vis[i] {
           continue
       }
       if l >= 0 {
           fmt.Fprintln(writer, pts[i].id, pts[l].id)
           vis[i], vis[l] = true, true
           l = -1
       } else {
           l = i
       }
   }
}
