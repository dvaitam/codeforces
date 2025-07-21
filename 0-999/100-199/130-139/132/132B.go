package main

import (
   "bufio"
   "fmt"
   "os"
)

type Block struct {
   color byte
   minx, miny, maxx, maxy int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var m int
   var n int64
   if _, err := fmt.Fscan(in, &m, &n); err != nil {
       return
   }
   grid := make([][]byte, m)
   var line string
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &line)
       grid[i] = []byte(line)
   }
   h := m
   w := len(grid[0])
   // label blocks
   blkid := make([][]int, h)
   for i := range blkid {
       blkid[i] = make([]int, w)
       for j := range blkid[i] {
           blkid[i][j] = -1
       }
   }
   var blocks []Block
   dirs := [4][2]int{{1,0},{-1,0},{0,1},{0,-1}}
   for y := 0; y < h; y++ {
       for x := 0; x < w; x++ {
           if grid[y][x] != '0' && blkid[y][x] < 0 {
               // flood fill
               col := grid[y][x]
               q := [][2]int{{x,y}}
               blkid[y][x] = len(blocks)
               minx, maxx, miny, maxy := x, x, y, y
               for qi := 0; qi < len(q); qi++ {
                   cx, cy := q[qi][0], q[qi][1]
                   if cx < minx { minx = cx }
                   if cx > maxx { maxx = cx }
                   if cy < miny { miny = cy }
                   if cy > maxy { maxy = cy }
                   for _, d := range dirs {
                       nx, ny := cx+d[0], cy+d[1]
                       if nx>=0 && nx<w && ny>=0 && ny<h && blkid[ny][nx]<0 && grid[ny][nx]==col {
                           blkid[ny][nx] = len(blocks)
                           q = append(q, [2]int{nx, ny})
                       }
                   }
               }
               blocks = append(blocks, Block{color: col, minx: minx, miny: miny, maxx: maxx, maxy: maxy})
           }
       }
   }
   B := len(blocks)
   // transitions
   // DP: 0=right,1=down,2=left,3=up; CP: 0=left,1=right
   dx := [4]int{1, 0, -1, 0}
   dy := [4]int{0, 1, 0, -1}
   trans := make([]int, B*8)
   for id := 0; id < B; id++ {
       blk := blocks[id]
       for dp := 0; dp < 4; dp++ {
           for cp := 0; cp < 2; cp++ {
               state := id*8 + dp*2 + cp
               // find target pixel
               var tx, ty int
               // edge in DP
               switch dp {
               case 0: // right
                   tx = blk.maxx
                   // among y in [miny,maxy], pick by CP
                   if cp == 0 { // left: CPdir = up
                       ty = blk.miny
                   } else {
                       ty = blk.maxy
                   }
               case 2: // left
                   tx = blk.minx
                   if cp == 0 {
                       ty = blk.maxy
                   } else {
                       ty = blk.miny
                   }
               case 1: // down
                   ty = blk.maxy
                   if cp == 0 {
                       tx = blk.maxx
                   } else {
                       tx = blk.minx
                   }
               case 3: // up
                   ty = blk.miny
                   if cp == 0 {
                       tx = blk.minx
                   } else {
                       tx = blk.maxx
                   }
               }
               nx, ny := tx+dx[dp], ty+dy[dp]
               if nx>=0 && nx<w && ny>=0 && ny<h && grid[ny][nx] != '0' {
                   nid := blkid[ny][nx]
                   trans[state] = nid*8 + dp*2 + cp
               } else {
                   // bounce
                   if cp == 0 {
                       // left->right
                       trans[state] = id*8 + dp*2 + 1
                   } else {
                       // right->left, rotate dp
                       ndp := (dp+1)&3
                       trans[state] = id*8 + ndp*2 + 0
                   }
               }
           }
       }
   }
   // initial state
   startBlk := blkid[0][0]
   state := startBlk*8 + 0*2 + 0
   // simulate with cycle detection
   visited := make([]int, B*8)
   for i := range visited {
       visited[i] = -1
   }
   order := make([]int, 0, 10000)
   var step int64
   for step = 0; step < n && visited[state] < 0; step++ {
       visited[state] = int(step)
       order = append(order, state)
       state = trans[state]
   }
   if step < n {
       // found cycle
       cycleStart := visited[state]
       cycleLen := int(step) - cycleStart
       rem := (n - step) % int64(cycleLen)
       state = order[cycleStart+int(rem)]
   }
   // output color of block
   bid := state / 8
   fmt.Printf("%c\n", blocks[bid].color)
}
