package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   fmt.Fscan(reader, &n, &m)
   a := make([][]int, n)
   for i := 0; i < n; i++ {
       a[i] = make([]int, m)
       for j := 0; j < m; j++ {
           var v int
           fmt.Fscan(reader, &v)
           a[i][j] = v - 1
       }
   }
   total := n * m
   rMax := make([]bool, total)
   cMax := make([]bool, total)
   for i := 0; i < n; i++ {
       mx := 0
       for j := 0; j < m; j++ {
           if a[i][j] > mx {
               mx = a[i][j]
           }
       }
       rMax[mx] = true
   }
   for j := 0; j < m; j++ {
       mx := 0
       for i := 0; i < n; i++ {
           if a[i][j] > mx {
               mx = a[i][j]
           }
       }
       cMax[mx] = true
   }
   px := make([]int, total)
   py := make([]int, total)
   x, y := -1, -1
   for i := total - 1; i >= 0; i-- {
       if rMax[i] || cMax[i] {
           if rMax[i] {
               x++
           }
           if cMax[i] {
               y++
           }
           a[x][y] = i
           px[i] = x
           py[i] = y
       }
   }
   cur := 0
   for i := 0; i < total; i++ {
       if rMax[i] {
           for y0 := 0; y0 < py[i]; y0++ {
               for rMax[cur] || cMax[cur] {
                   cur++
               }
               a[px[i]][y0] = cur
               cur++
           }
       }
       if cMax[i] {
           for x0 := 0; x0 < px[i]; x0++ {
               for rMax[cur] || cMax[cur] {
                   cur++
               }
               a[x0][py[i]] = cur
               cur++
           }
       }
   }
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           fmt.Fprint(writer, a[i][j]+1)
           if j+1 < m {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
