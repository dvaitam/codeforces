package main

import (
   "bufio"
   "fmt"
   "os"
)

const MOD = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n, k int
   fmt.Fscan(reader, &n, &k)
   weights := make([]int, n)
   total50, total100 := 0, 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &weights[i])
       if weights[i] == 50 {
           total50++
       } else {
           total100++
       }
   }
   // precompute combinations C[n][k]
   maxN := n
   C := make([][]int, maxN+1)
   for i := 0; i <= maxN; i++ {
       C[i] = make([]int, i+1)
       C[i][0] = 1
       for j := 1; j <= i; j++ {
           if j == i {
               C[i][j] = 1
           } else {
               C[i][j] = (C[i-1][j-1] + C[i-1][j]) % MOD
           }
       }
   }
   // possible moves (x50, y100) with weight <= k
   type move struct{ x, y int }
   moves := make([]move, 0)
   maxY := k / 100
   for y := 0; y <= maxY; y++ {
       rem := k - 100*y
       maxX := rem / 50
       for x := 0; x <= maxX; x++ {
           if x == 0 && y == 0 {
               continue
           }
           moves = append(moves, move{x, y})
       }
   }
   // dist[a50_left][a100_left][side], side: 0=left,1=right
   const INF = 1<<30
   dist := make([][][]int, total50+1)
   ways := make([][][]int, total50+1)
   for a := 0; a <= total50; a++ {
       dist[a] = make([][]int, total100+1)
       ways[a] = make([][]int, total100+1)
       for b := 0; b <= total100; b++ {
           dist[a][b] = []int{INF, INF}
           ways[a][b] = []int{0, 0}
       }
   }
   type state struct{ a, b, side int }
   // BFS
   queue := make([]state, 0, (total50+1)*(total100+1)*2)
   dist[total50][total100][0] = 0
   ways[total50][total100][0] = 1
   queue = append(queue, state{total50, total100, 0})
   head := 0
   for head < len(queue) {
       cur := queue[head]
       head++
       d := dist[cur.a][cur.b][cur.side]
       wcur := ways[cur.a][cur.b][cur.side]
       for _, mv := range moves {
           x, y := mv.x, mv.y
           var na, nb, nside int
           var c50, c100 int
           if cur.side == 0 {
               // move from left to right
               if x > cur.a || y > cur.b {
                   continue
               }
               na = cur.a - x
               nb = cur.b - y
               nside = 1
               c50 = cur.a
               c100 = cur.b
           } else {
               // move from right to left
               right50 := total50 - cur.a
               right100 := total100 - cur.b
               if x > right50 || y > right100 {
                   continue
               }
               na = cur.a + x
               nb = cur.b + y
               nside = 0
               c50 = right50
               c100 = right100
           }
           nd := d + 1
           addWays := int64(C[c50][x]) * int64(C[c100][y]) % MOD
           if dist[na][nb][nside] > nd {
               dist[na][nb][nside] = nd
               ways[na][nb][nside] = int((int64(wcur) * addWays) % MOD)
               queue = append(queue, state{na, nb, nside})
           } else if dist[na][nb][nside] == nd {
               ways[na][nb][nside] = int((int64(ways[na][nb][nside]) + int64(wcur)*addWays) % MOD)
           }
       }
   }
   resDist := dist[0][0][1]
   resWays := ways[0][0][1]
   if resDist >= INF {
       fmt.Fprintln(writer, -1)
       fmt.Fprintln(writer, 0)
   } else {
       fmt.Fprintln(writer, resDist)
       fmt.Fprintln(writer, resWays)
   }
}
