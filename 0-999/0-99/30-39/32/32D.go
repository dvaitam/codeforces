package main

import (
   "bufio"
   "fmt"
   "os"
)

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, m int
   var k int64
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   grid := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &grid[i])
   }

   // Precompute continuous star lengths including self
   up := make([][]int, n)
   down := make([][]int, n)
   left := make([][]int, n)
   right := make([][]int, n)
   for i := 0; i < n; i++ {
       up[i] = make([]int, m)
       down[i] = make([]int, m)
       left[i] = make([]int, m)
       right[i] = make([]int, m)
   }
   // up and left
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if grid[i][j] == '*' {
               if i > 0 {
                   up[i][j] = up[i-1][j] + 1
               } else {
                   up[i][j] = 1
               }
               if j > 0 {
                   left[i][j] = left[i][j-1] + 1
               } else {
                   left[i][j] = 1
               }
           }
       }
   }
   // down and right
   for i := n - 1; i >= 0; i-- {
       for j := m - 1; j >= 0; j-- {
           if grid[i][j] == '*' {
               if i < n-1 {
                   down[i][j] = down[i+1][j] + 1
               } else {
                   down[i][j] = 1
               }
               if j < m-1 {
                   right[i][j] = right[i][j+1] + 1
               } else {
                   right[i][j] = 1
               }
           }
       }
   }
   // compute max radius per center (excluding self)
   maxR := 0
   rgrid := make([][]int, n)
   for i := 0; i < n; i++ {
       rgrid[i] = make([]int, m)
       for j := 0; j < m; j++ {
           if grid[i][j] == '*' {
               ru := up[i][j] - 1
               rd := down[i][j] - 1
               rl := left[i][j] - 1
               rr := right[i][j] - 1
               r := ru
               if rd < r {
                   r = rd
               }
               if rl < r {
                   r = rl
               }
               if rr < r {
                   r = rr
               }
               rgrid[i][j] = r
               if r > maxR {
                   maxR = r
               }
           }
       }
   }
   if maxR <= 0 {
       // no crosses
       fmt.Fprintln(writer, -1)
       return
   }
   // count crosses per radius
   freq := make([]int64, maxR+2)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           r := rgrid[i][j]
           if r > 0 {
               freq[1]++
               if r+1 <= maxR {
                   freq[r+1]--
               }
           }
       }
   }
   // build countPerRadius by prefix sum on freq
   countR := make([]int64, maxR+2)
   var cur int64
   for x := 1; x <= maxR; x++ {
       cur += freq[x]
       countR[x] = cur
   }
   // find target radius
   target := 0
   for x := 1; x <= maxR; x++ {
       if k <= countR[x] {
           target = x
           break
       }
       k -= countR[x]
   }
   if target == 0 {
       fmt.Fprintln(writer, -1)
       return
   }
   // find k-th center with rgrid>=target in order
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           if rgrid[i][j] >= target {
               k--
               if k == 0 {
                   // output 5 stars: center, up, down, left, right
                   ci, cj := i+1, j+1
                   fmt.Fprintf(writer, "%d %d\n", ci, cj)
                   fmt.Fprintf(writer, "%d %d\n", ci-target, cj)
                   fmt.Fprintf(writer, "%d %d\n", ci+target, cj)
                   fmt.Fprintf(writer, "%d %d\n", ci, cj-target)
                   fmt.Fprintf(writer, "%d %d\n", ci, cj+target)
                   return
               }
           }
       }
   }
   // not found
   fmt.Fprintln(writer, -1)
}
