package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   board := make([]string, 8)
   for i := 0; i < 8; i++ {
       line, _ := reader.ReadString('\n')
       // remove trailing newline or carriage return
       if len(line) > 0 && (line[len(line)-1] == '\n' || line[len(line)-1] == '\r') {
           line = line[:len(line)-1]
       }
       if len(line) > 0 && line[len(line)-1] == '\r' {
           line = line[:len(line)-1]
       }
       board[i] = line
   }
   // movements: stay, 8 directions
   moves := [9][2]int{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
   // visited[r][c][t]
   var visited [8][8][9]bool
   type state struct{ r, c, t int }
   queue := []state{{7, 0, 0}}
   visited[7][0][0] = true
   // BFS
   for qi := 0; qi < len(queue); qi++ {
       s := queue[qi]
       r, c, t := s.r, s.c, s.t
       // if current cell has a statue at time t, skip
       if hasStone(board, r, c, t) {
           continue
       }
       // reached Anna
       if r == 0 && c == 7 {
           fmt.Println("WIN")
           return
       }
       // after 8 moves, no stones remain
       if t >= 8 {
           fmt.Println("WIN")
           return
       }
       // try all moves
       for _, mv := range moves {
           nr, nc := r+mv[0], c+mv[1]
           nt := t + 1
           if nr < 0 || nr >= 8 || nc < 0 || nc >= 8 {
               continue
           }
           // cell must be free at time t and t+1
           if hasStone(board, nr, nc, t) || hasStone(board, nr, nc, nt) {
               continue
           }
           if visited[nr][nc][nt] {
               continue
           }
           visited[nr][nc][nt] = true
           queue = append(queue, state{nr, nc, nt})
       }
   }
   fmt.Println("LOSE")
}

// hasStone checks if a statue occupies (r,c) at time t
func hasStone(board []string, r, c, t int) bool {
   pr := r - t
   if pr < 0 {
       return false
   }
   return board[pr][c] == 'S'
}
