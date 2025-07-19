package main

import (
   "bufio"
   "fmt"
   "os"
)

const MAX = 505

var (
   n int
   a [MAX][MAX]int
   visit [MAX][MAX]bool
   total int
   opt []byte
   dx = [4]int{-1, 1, 0, 0}
   dy = [4]int{0, 0, -1, 1}
   opts = [4]byte{'U', 'D', 'L', 'R'}
)

func canGo(x, y, dir int) bool {
   x += dx[dir]
   y += dy[dir]
   for x > 0 && y > 0 && x <= n && y <= n {
       if visit[x][y] {
           return false
       }
       if a[x][y] == 1 {
           return true
       }
       x += dx[dir]
       y += dy[dir]
   }
   return false
}

func dfs(x, y int) {
   visit[x][y] = true
   if a[x][y] == 0 {
       total++
       opt = append(opt, '1')
   }
   for i := 0; i < 4; i++ {
       if canGo(x, y, i) {
           opt = append(opt, opts[i])
           dfs(x+dx[i], y+dy[i])
           opt = append(opt, opts[i^1])
       }
   }
   opt = append(opt, '2')
   total--
   a[x][y] = 0
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var x0, y0 int
   fmt.Fscan(reader, &n, &x0, &y0)
   opt = make([]byte, 0, 3000000)
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           fmt.Fscan(reader, &a[i][j])
           if a[i][j] == 1 {
               total++
           }
       }
   }
   dfs(x0, y0)
   if total > 0 {
       fmt.Fprintln(writer, "NO")
   } else {
       fmt.Fprintln(writer, "YES")
       writer.Write(opt)
       writer.WriteByte('\n')
   }
}
