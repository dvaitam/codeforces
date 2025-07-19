package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m int
   a [][]int
   b [][]int
   bld, nw []int
   ans int
)

func gen(pos int, curSum int, fb []int) {
   if pos == n {
       if curSum > ans {
           ans = curSum
           copy(bld, nw)
       }
       return
   }
   if len(fb) == 0 {
       // no forbidden, take best for all remaining
       for i := pos; i < n; i++ {
           last := len(a[i]) - 1
           curSum += a[i][last]
           nw[i] = last + 1
       }
       if curSum > ans {
           ans = curSum
           copy(bld, nw)
       }
       return
   }
   // group forbidden by choice at pos
   mp := make(map[int][]int)
   for _, ind := range fb {
       v := b[ind][pos]
       mp[v] = append(mp[v], ind)
   }
   // try choices from largest to smallest
   for i := len(a[pos]) - 1; i >= 0; i-- {
       nw[pos] = i + 1
       nextSum := curSum + a[pos][i]
       fb2 := mp[i]
       gen(pos+1, nextSum, fb2)
       if len(fb2) == 0 {
           break
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n)
   a = make([][]int, n)
   for i := 0; i < n; i++ {
       var c int
       fmt.Fscan(reader, &c)
       a[i] = make([]int, c)
       for j := 0; j < c; j++ {
           fmt.Fscan(reader, &a[i][j])
       }
   }
   fmt.Fscan(reader, &m)
   b = make([][]int, m)
   for i := 0; i < m; i++ {
       b[i] = make([]int, n)
       for j := 0; j < n; j++ {
           var x int
           fmt.Fscan(reader, &x)
           b[i][j] = x - 1
       }
   }
   bld = make([]int, n)
   nw = make([]int, n)
   ans = 0
   // initial forbidden list: all sequences
   fb := make([]int, m)
   for i := 0; i < m; i++ {
       fb[i] = i
   }
   gen(0, 0, fb)
   // output
   for i, v := range bld {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v)
   }
   writer.WriteByte('\n')
}
