package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M, K int
   if _, err := fmt.Fscan(in, &N, &M, &K); err != nil {
       return
   }
   grid := make([][]rune, N)
   for i := 0; i < N; i++ {
       var line string
       fmt.Fscan(in, &line)
       grid[i] = []rune(line)
   }
   visited := make([][]bool, N)
   for i := range visited {
       visited[i] = make([]bool, M)
   }
   type node struct{ x, y, val int }
   lakes := make([]node, 0)
   dx := []int{0, 0, -1, 1}
   dy := []int{1, -1, 0, 0}
   var num int
   var flag bool
   var dfs func(x, y int)
   dfs = func(x, y int) {
       if visited[x][y] || grid[x][y] == '*' {
           return
       }
       visited[x][y] = true
       num++
       for dir := 0; dir < 4; dir++ {
           nx, ny := x+dx[dir], y+dy[dir]
           if nx < 0 || ny < 0 || nx >= N || ny >= M {
               flag = false
               continue
           }
           if !visited[nx][ny] && grid[nx][ny] == '.' {
               dfs(nx, ny)
           }
       }
   }
   // identify fully enclosed lakes
   for i := 0; i < N; i++ {
       for j := 0; j < M; j++ {
           if !visited[i][j] && grid[i][j] == '.' {
               num = 0
               flag = true
               dfs(i, j)
               if flag {
                   lakes = append(lakes, node{i, j, num})
               }
           }
       }
   }
   // sort lakes by size descending
   sort.Slice(lakes, func(i, j int) bool {
       return lakes[i].val > lakes[j].val
   })
   var ans int
   var dfsl func(x, y int)
   dfsl = func(x, y int) {
       if grid[x][y] == '*' {
           return
       }
       grid[x][y] = '*'
       for dir := 0; dir < 4; dir++ {
           nx, ny := x+dx[dir], y+dy[dir]
           if nx >= 0 && ny >= 0 && nx < N && ny < M && grid[nx][ny] == '.' {
               dfsl(nx, ny)
           }
       }
   }
   // fill all but the largest K lakes
   for idx := K; idx < len(lakes); idx++ {
       ans += lakes[idx].val
       dfsl(lakes[idx].x, lakes[idx].y)
   }
   // output
   fmt.Println(ans)
   for i := 0; i < N; i++ {
       fmt.Println(string(grid[i]))
   }
}
