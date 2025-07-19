package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var T int
   if _, err := fmt.Fscan(reader, &T); err != nil {
       return
   }
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for tc := 0; tc < T; tc++ {
       var n, m int
       fmt.Fscan(reader, &n, &m)
       grid := make([][]int, n)
       for i := 0; i < n; i++ {
           var line string
           fmt.Fscan(reader, &line)
           grid[i] = make([]int, m)
           for j := 0; j < m; j++ {
               if line[j] == '1' {
                   grid[i][j] = 1
               }
           }
       }
       type op [6]int
       var ops []op
       // flip function, coordinates 0-based
       add := func(a, b, c, d, e, f int) {
           grid[a][b] ^= 1
           grid[c][d] ^= 1
           grid[e][f] ^= 1
           ops = append(ops, op{a, b, c, d, e, f})
       }
       // process rows except last two
       for i := 0; i <= n-3; i++ {
           for j := 0; j <= m-2; j++ {
               if grid[i][j] == 1 {
                   add(i, j, i+1, j, i, j+1)
               }
           }
           if grid[i][m-1] == 1 {
               add(i, m-1, i+1, m-2, i+1, m-1)
           }
       }
       // process columns except last two in last two rows
       for j := 0; j <= m-3; j++ {
           a := grid[n-2][j]
           b := grid[n-1][j]
           if a == 1 && b == 1 {
               add(n-2, j, n-1, j, n-1, j+1)
           } else if a == 1 {
               add(n-2, j, n-2, j+1, n-1, j+1)
           } else if b == 1 {
               add(n-1, j, n-2, j+1, n-1, j+1)
           }
       }
       // bottom-right 2x2
       var ones, zeros [][2]int
       for di := 0; di < 2; di++ {
           for dj := 0; dj < 2; dj++ {
               x, y := n-2+di, m-2+dj
               if grid[x][y] == 1 {
                   ones = append(ones, [2]int{x, y})
               } else {
                   zeros = append(zeros, [2]int{x, y})
               }
           }
       }
       switch len(ones) {
       case 1:
           // one one, three zeros
           o := ones[0]
           z0, z1, z2 := zeros[0], zeros[1], zeros[2]
           add(o[0], o[1], z0[0], z0[1], z1[0], z1[1])
           add(o[0], o[1], z0[0], z0[1], z2[0], z2[1])
           add(o[0], o[1], z1[0], z1[1], z2[0], z2[1])
       case 2:
           // two ones, two zeros
           o0, o1 := ones[0], ones[1]
           z0, z1 := zeros[0], zeros[1]
           add(o0[0], o0[1], z0[0], z0[1], z1[0], z1[1])
           add(o1[0], o1[1], z0[0], z0[1], z1[0], z1[1])
       case 3:
           // three ones, one zero
           o0, o1, o2 := ones[0], ones[1], ones[2]
           add(o0[0], o0[1], o1[0], o1[1], o2[0], o2[1])
       case 4:
           // four ones
           v0, v1, v2, v3 := ones[0], ones[1], ones[2], ones[3]
           add(v0[0], v0[1], v1[0], v1[1], v2[0], v2[1])
           add(v0[0], v0[1], v2[0], v2[1], v3[0], v3[1])
           add(v0[0], v0[1], v1[0], v1[1], v3[0], v3[1])
           add(v1[0], v1[1], v2[0], v2[1], v3[0], v3[1])
       }
       // output
       fmt.Fprintln(writer, len(ops))
       for _, op := range ops {
           // convert to 1-based
           fmt.Fprintf(writer, "%d %d %d %d %d %d\n",
               op[0]+1, op[1]+1, op[2]+1, op[3]+1, op[4]+1, op[5]+1)
       }
   }
}
