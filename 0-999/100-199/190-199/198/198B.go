package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   fmt.Fscan(reader, &n, &k)
   // read wall descriptions
   left := make([]byte, n+2)
   right := make([]byte, n+2)
   var s string
   fmt.Fscan(reader, &s)
   for i := 1; i <= n; i++ {
       left[i] = s[i-1]
   }
   fmt.Fscan(reader, &s)
   for i := 1; i <= n; i++ {
       right[i] = s[i-1]
   }
   // visited[wall][position]
   visited := make([][]bool, 2)
   visited[0] = make([]bool, n+2)
   visited[1] = make([]bool, n+2)
   type state struct{ wall, pos, time int }
   q := make([]state, 0, n*2)
   // start at left wall (0), position 1, time 0
   q = append(q, state{0, 1, 0})
   visited[0][1] = true
   for head := 0; head < len(q); head++ {
       cur := q[head]
       w, p, t := cur.wall, cur.pos, cur.time
       // possible moves: up, down, jump to other wall
       moves := []int{p + 1, p - 1, p + k}
       walls := []int{w, w, 1 - w}
       for i, np := range moves {
           nw := walls[i]
           nt := t + 1
           // success if out of canyon
           if np > n {
               fmt.Println("YES")
               return
           }
           // cannot move into flooded or invalid positions
           if np <= nt || np < 1 {
               continue
           }
           if visited[nw][np] {
               continue
           }
           // cannot step on dangerous area
           if (nw == 0 && left[np] == 'X') || (nw == 1 && right[np] == 'X') {
               continue
           }
           visited[nw][np] = true
           q = append(q, state{nw, np, nt})
       }
   }
   fmt.Println("NO")
}
