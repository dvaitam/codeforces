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

   var t int
   if _, err := fmt.Fscan(in, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       // collect zeros and positions of non-zeros
       zeros := make([]int, 0, n)
       type px struct{ pos, x int }
       posList := make([]px, 0, n)
       for i, v := range a {
           if v == 0 {
               zeros = append(zeros, i)
           } else {
               posList = append(posList, px{pos: i, x: v})
           }
       }
       sort.Ints(zeros)
       sort.Slice(posList, func(i, j int) bool {
           return posList[i].pos < posList[j].pos
       })
       zlen := len(zeros)
       used := make([]bool, n+1)
       l, r := 0, zlen-1
       p := 0
       cnt := 0
       for l < r {
           zl := zeros[l]
           zr := zeros[r]
           // find pos > zl
           for p < len(posList) && posList[p].pos <= zl {
               p++
           }
           // find pos < zr and unused x
           for p < len(posList) && (posList[p].pos >= zr || used[posList[p].x]) {
               p++
           }
           if p >= len(posList) {
               break
           }
           // use this triple
           used[posList[p].x] = true
           cnt++
           l++
           r--
           p++
       }
       fmt.Fprintln(out, cnt)
   }
}
