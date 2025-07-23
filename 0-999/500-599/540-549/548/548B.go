package main

import (
   "bufio"
   "fmt"
   "os"
)

func calcMaxRow(row []int) int {
   maxCnt := 0
   cnt := 0
   for _, v := range row {
       if v == 1 {
           cnt++
           if cnt > maxCnt {
               maxCnt = cnt
           }
       } else {
           cnt = 0
       }
   }
   return maxCnt
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m, q int
   if _, err := fmt.Fscan(reader, &n, &m, &q); err != nil {
       return
   }
   grid := make([][]int, n)
   rowMax := make([]int, n)
   for i := 0; i < n; i++ {
       grid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           fmt.Fscan(reader, &grid[i][j])
       }
       rowMax[i] = calcMaxRow(grid[i])
   }

   for k := 0; k < q; k++ {
       var i, j int
       fmt.Fscan(reader, &i, &j)
       i--
       j--
       // flip state
       if grid[i][j] == 1 {
           grid[i][j] = 0
       } else {
           grid[i][j] = 1
       }
       // recompute this row's max
       rowMax[i] = calcMaxRow(grid[i])
       // compute global max
       ans := 0
       for r := 0; r < n; r++ {
           if rowMax[r] > ans {
               ans = rowMax[r]
           }
       }
       fmt.Fprintln(writer, ans)
   }
}
