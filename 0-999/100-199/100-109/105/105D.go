package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const MAXP = 605

type temp struct {
   x, y, c int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   // mapping of original values to compressed indices
   g := make(map[int]int)
   s := 1
   g[0] = 1
   // slices for counts and vectors
   c := make([]int, 2)
   v := make([][]temp, 2)
   // grids a and b, 1-indexed
   a := make([][]int, n+2)
   b := make([][]int, n+2)
   for i := range a {
       a[i] = make([]int, m+2)
       b[i] = make([]int, m+2)
   }
   // read grid a
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           var k int
           fmt.Fscan(reader, &k)
           id, ok := g[k]
           if !ok {
               s++
               id = s
               g[k] = id
               // extend slices if needed
               if id >= len(c) {
                   nc := make([]int, id+1)
                   copy(nc, c)
                   c = nc
                   nv := make([][]temp, id+1)
                   copy(nv, v)
                   v = nv
               }
           }
           a[i][j] = id
           c[id]++
       }
   }
   // read grid b
   for i := 1; i <= n; i++ {
       for j := 1; j <= m; j++ {
           var k int
           fmt.Fscan(reader, &k)
           if k > -1 {
               id, ok := g[k]
               if !ok {
                   s++
                   id = s
                   g[k] = id
                   if id >= len(c) {
                       nc := make([]int, id+1)
                       copy(nc, c)
                       c = nc
                       nv := make([][]temp, id+1)
                       copy(nv, v)
                       v = nv
                   }
               }
               b[i][j] = id
               // append to vector for a[i][j]
               ai := a[i][j]
               v[ai] = append(v[ai], temp{i, j, id})
           }
       }
   }
   // starting position
   var px, py int
   fmt.Fscan(reader, &px, &py)
   // build spiral priority p
   var p [MAXP][MAXP]int
   p[n][m] = 1
   k := 1
   x, y := n, m
   for step := 2; step <= MAXP; step += 2 {
       x--
       y--
       // right
       for i := 0; i < step; i++ {
           y++
           k++
           p[x][y] = k
       }
       // down
       for i := 0; i < step; i++ {
           x++
           k++
           p[x][y] = k
       }
       // left
       for i := 0; i < step; i++ {
           y--
           k++
           p[x][y] = k
       }
       // up
       for i := 0; i < step; i++ {
           x--
           k++
           p[x][y] = k
       }
       // stop if filled enough
       if step > 2*n && step > 2*m {
           break
       }
   }
   // processing
   ans := int64(0)
   // queue of temp
   q := make([]temp, 0, 1024)
   // initial
   q = append(q, temp{px, py, b[px][py]})
   curColor := a[px][py]
   for l := 0; l < len(q); l++ {
       it := q[l]
       newColor := it.c
       if curColor > 1 && newColor != curColor && c[curColor] > 0 {
           // convert curColor to newColor
           ans += int64(c[curColor])
           c[newColor] += c[curColor]
           c[curColor] = 0
           // update position and color
           px = it.x
           py = it.y
           // sort v[curColor] by spiral distance
           vs := v[curColor]
           sort.Slice(vs, func(i, j int) bool {
               xi, yi := vs[i].x, vs[i].y
               xj, yj := vs[j].x, vs[j].y
               return p[xi-px+n][yi-py+m] < p[xj-px+n][yj-py+m]
           })
           // append all
           for _, tt := range vs {
               q = append(q, tt)
           }
           // clear
           v[curColor] = nil
           curColor = newColor
       }
   }
   fmt.Fprintln(writer, ans)
}
