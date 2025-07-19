package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Cube struct {
   id    int
   color int
   size  int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var n int
   fmt.Fscan(in, &n)
   cubes := make([]Cube, n)
   for i := 0; i < n; i++ {
       cubes[i].id = i + 1
       fmt.Fscan(in, &cubes[i].color, &cubes[i].size)
   }
   sort.Slice(cubes, func(i, j int) bool {
       if cubes[i].color != cubes[j].color {
           return cubes[i].color < cubes[j].color
       }
       return cubes[i].size > cubes[j].size
   })
   // dp arrays
   const INF = int64(-9e18)
   mx := make([]int64, n+5)
   mxidx := make([]int, n+5)
   for i := range mx {
       mx[i] = INF
   }
   start := make([]int, n)
   var ans int64
   var ansidx1, ansidx2 int
   updateMx := func(length int, v int64, idx int) {
       if mx[length] < v {
           mx[length] = v
           mxidx[length] = idx
       }
   }
   updateAns := func(v int64, i1, i2 int) {
       if ans < v {
           ans = v
           ansidx1 = i1
           ansidx2 = i2
       }
   }
   // process groups by color
   for r := 0; r < n; {
       l := r
       for r < n && cubes[r].color == cubes[l].color {
           r++
       }
       // combine with previous
       var sum int64
       sum = 0
       for i := l; i < r; i++ {
           start[i] = l
           sum += int64(cubes[i].size)
           length := i - l + 1
           // try len-1, len, len+1
           if length-1 >= 0 {
               updateAns(mx[length-1]+sum, mxidx[length-1], i)
           }
           updateAns(mx[length]+sum, mxidx[length], i)
           updateAns(mx[length+1]+sum, mxidx[length+1], i)
       }
       // update mx
       sum = 0
       for i := l; i < r; i++ {
           sum += int64(cubes[i].size)
           length := i - l + 1
           updateMx(length, sum, i)
       }
   }
   // output
   fmt.Fprintln(out, ans)
   len1 := ansidx1 - start[ansidx1] + 1
   len2 := ansidx2 - start[ansidx2] + 1
   total := len1 + len2
   fmt.Fprintln(out, total)
   // ensure len1 >= len2
   if len1 < len2 {
       len1, len2 = len2, len1
       ansidx1, ansidx2 = ansidx2, ansidx1
   }
   // print sequence
   for len1 > 0 || len2 > 0 {
       if len1 > len2 {
           fmt.Fprintf(out, "%d ", cubes[ansidx1].id)
           ansidx1--
           len1--
       } else {
           fmt.Fprintf(out, "%d ", cubes[ansidx2].id)
           ansidx2--
           len2--
       }
   }
   fmt.Fprintln(out)
}
