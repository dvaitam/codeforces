package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// DSU for union-find
type DSU struct {
   p []int
}

func NewDSU(n int) *DSU {
   d := &DSU{p: make([]int, n)}
   for i := range d.p {
       d.p[i] = -1
   }
   return d
}

func (d *DSU) Find(x int) int {
   if d.p[x] < 0 {
       return x
   }
   d.p[x] = d.Find(d.p[x])
   return d.p[x]
}

func (d *DSU) Union(a, b int) bool {
   a = d.Find(a)
   b = d.Find(b)
   if a == b {
       return false
   }
   if d.p[a] > d.p[b] {
       a, b = b, a
   }
   d.p[a] += d.p[b]
   d.p[b] = a
   return true
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(in, &n)
   pts := make([][2]int, n)
   rowMap := make(map[int][][2]int)
   colMap := make(map[int][][2]int)
   xsSet := make(map[int]struct{})
   ysSet := make(map[int]struct{})
   exists := make(map[int64]struct{})
   for i := 0; i < n; i++ {
       x, y := 0, 0
       fmt.Fscan(in, &x, &y)
       pts[i] = [2]int{x, y}
       rowMap[y] = append(rowMap[y], [2]int{x, i})
       colMap[x] = append(colMap[x], [2]int{y, i})
       xsSet[x] = struct{}{}
       ysSet[y] = struct{}{}
       key := (int64(x) << 32) | int64(uint32(y))
       exists[key] = struct{}{}
   }
   xs := make([]int, 0, len(xsSet))
   ys := make([]int, 0, len(ysSet))
   for x := range xsSet {
       xs = append(xs, x)
   }
   for y := range ysSet {
       ys = append(ys, y)
   }
   sort.Ints(xs)
   sort.Ints(ys)
   // sort rows and cols
   for y, arr := range rowMap {
       sort.Slice(arr, func(i, j int) bool { return arr[i][0] < arr[j][0] })
       rowMap[y] = arr
   }
   for x, arr := range colMap {
       sort.Slice(arr, func(i, j int) bool { return arr[i][0] < arr[j][0] })
       colMap[x] = arr
   }

   // check if some t works
   var check func(t int) bool
   check = func(t int) bool {
       dsu := NewDSU(n)
       // row unions
       for _, arr := range rowMap {
           for i := 0; i+1 < len(arr); i++ {
               x0, id0 := arr[i][0], arr[i][1]
               x1, id1 := arr[i+1][0], arr[i+1][1]
               if x1-x0 <= t {
                   dsu.Union(id0, id1)
               }
           }
       }
       // col unions
       for _, arr := range colMap {
           for i := 0; i+1 < len(arr); i++ {
               y0, id0 := arr[i][0], arr[i][1]
               y1, id1 := arr[i+1][0], arr[i+1][1]
               if y1-y0 <= t {
                   dsu.Union(id0, id1)
               }
           }
       }
       // collect components
       compMap := make(map[int]int)
       compCnt := 0
       compID := make([]int, n)
       for i := 0; i < n; i++ {
           r := dsu.Find(i)
           if _, ok := compMap[r]; !ok {
               compMap[r] = compCnt
               compCnt++
           }
           compID[i] = compMap[r]
       }
       if compCnt == 1 {
           fmt.Fprintln(os.Stderr, "debug: compCnt==1 for t=", t)
           return true
       }
       if compCnt > 4 {
           return false
       }
       // try adding one node
       var cover bool
       for _, x := range xs {
           colArr := colMap[x]
           ysCol := make([]int, len(colArr))
           idsCol := make([]int, len(colArr))
           for i, p := range colArr {
               ysCol[i] = p[0]
               idsCol[i] = p[1]
           }
           for _, y := range ys {
               key := (int64(x) << 32) | int64(uint32(y))
               if _, ok := exists[key]; ok {
                   continue
               }
               comps := make(map[int]struct{})
               // col neighbors
               j := sort.SearchInts(ysCol, y)
               if j > 0 {
                   if y-ysCol[j-1] <= t {
                       comps[compID[idsCol[j-1]]] = struct{}{}
                   }
               }
               if j < len(ysCol) {
                   if ysCol[j]-y <= t {
                       comps[compID[idsCol[j]]] = struct{}{}
                   }
               }
               // row neighbors
               rowArr := rowMap[y]
               xsRow := make([]int, len(rowArr))
               idsRow := make([]int, len(rowArr))
               for i, p := range rowArr {
                   xsRow[i] = p[0]
                   idsRow[i] = p[1]
               }
               k := sort.SearchInts(xsRow, x)
               if k > 0 {
                   if x-xsRow[k-1] <= t {
                       comps[compID[idsRow[k-1]]] = struct{}{}
                   }
               }
               if k < len(xsRow) {
                   if xsRow[k]-x <= t {
                       comps[compID[idsRow[k]]] = struct{}{}
                   }
               }
               if len(comps) == compCnt {
                   cover = true
                   break
               }
           }
           if cover {
               break
           }
       }
       return cover
   }

   // binary search t
   lo, hi := 0, int(4e9)
   ans := -1
   for lo <= hi {
       mid := lo + (hi-lo)/2
       if check(mid) {
           ans = mid
           hi = mid - 1
       } else {
           lo = mid + 1
       }
   }
   fmt.Println(ans)
}
