package main

import (
   "fmt"
)

type Pos struct{ x, y int }

func posFrom(s string) Pos {
   return Pos{int(s[0] - 'a'), int(s[1] - '1')}
}

// pathClear returns true if all squares between a and b (excluding endpoints) are empty
func pathClear(a, b Pos, occ [][]int) bool {
   dx, dy := 0, 0
   if b.x > a.x {
       dx = 1
   } else if b.x < a.x {
       dx = -1
   }
   if b.y > a.y {
       dy = 1
   } else if b.y < a.y {
       dy = -1
   }
   cx, cy := a.x+dx, a.y+dy
   for cx != b.x || cy != b.y {
       if occ[cy][cx] != 0 {
           return false
       }
       cx += dx
       cy += dy
   }
   return true
}

// isAttacked checks if square p is attacked by any white piece
func isAttacked(occ [][]int, rooks []Pos, wk Pos, p Pos) bool {
   // by rooks
   for _, r := range rooks {
       if r.x == p.x || r.y == p.y {
           if pathClear(r, p, occ) {
               return true
           }
       }
   }
   // by white king
   dx := wk.x - p.x
   if dx < 0 {
       dx = -dx
   }
   dy := wk.y - p.y
   if dy < 0 {
       dy = -dy
   }
   if dx <= 1 && dy <= 1 {
       return true
   }
   return false
}

func main() {
   var s1, s2, s3, s4 string
   if _, err := fmt.Scan(&s1, &s2, &s3, &s4); err != nil {
       return
   }
   rk1 := posFrom(s1)
   rk2 := posFrom(s2)
   wk := posFrom(s3)
   bk := posFrom(s4)
   // initial occupancy: 0 empty, 1 rook, 2 wking, 3 bking
   occ := make([][]int, 8)
   for i := range occ {
       occ[i] = make([]int, 8)
   }
   occ[rk1.y][rk1.x] = 1
   occ[rk2.y][rk2.x] = 1
   occ[wk.y][wk.x] = 2
   occ[bk.y][bk.x] = 3
   rooks := []Pos{rk1, rk2}
   // Check if black king is currently in check
   if !isAttacked(occ, rooks, wk, bk) {
       fmt.Println("OTHER")
       return
   }
   // try all king moves
   dirs := []int{-1, 0, 1}
   for _, dx := range dirs {
       for _, dy := range dirs {
           if dx == 0 && dy == 0 {
               continue
           }
           nx, ny := bk.x+dx, bk.y+dy
           if nx < 0 || nx > 7 || ny < 0 || ny > 7 {
               continue
           }
           // cannot move onto white king
           if nx == wk.x && ny == wk.y {
               continue
           }
           // determine new rooks after possible capture
           newRooks := make([]Pos, 0, 2)
           for _, r := range rooks {
               if r.x == nx && r.y == ny {
                   continue
               }
               newRooks = append(newRooks, r)
           }
           // build new occupancy
           occ2 := make([][]int, 8)
           for i := range occ2 {
               occ2[i] = make([]int, 8)
           }
           for _, r := range newRooks {
               occ2[r.y][r.x] = 1
           }
           occ2[wk.y][wk.x] = 2
           occ2[ny][nx] = 3
           // if not attacked in new position, it's not mate
           if !isAttacked(occ2, newRooks, wk, Pos{nx, ny}) {
               fmt.Println("OTHER")
               return
           }
       }
   }
   fmt.Println("CHECKMATE")
}
