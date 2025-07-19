package main

import (
   "bufio"
   "fmt"
   "os"
)

type Pair struct { y, x int }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var tt int
   fmt.Fscan(reader, &tt)
   for tt > 0 {
       tt--
       var H, W int
       fmt.Fscan(reader, &H, &W)
       grid := make([][]byte, H)
       for i := 0; i < H; i++ {
           var s string
           fmt.Fscan(reader, &s)
           grid[i] = []byte(s)
       }
       var ops [][]Pair

       // handle small regions
       handleit := func(c1 [2]Pair, c2 [2]Pair) {
           tl, bl := c1[0], c1[1]
           tr, br := c2[0], c2[1]
           cnt := 0
           if grid[tl.y][tl.x] == '1' {
               cnt++
           }
           if grid[bl.y][bl.x] == '1' {
               cnt++
           }
           if grid[tr.y][tr.x] == '1' {
               cnt++
           }
           if grid[br.y][br.x] == '1' {
               cnt++
           }
           switch cnt {
           case 1:
               if grid[tl.y][tl.x] == '0' && grid[bl.y][bl.x] == '0' {
                   return
               }
               op := []Pair{tr, br}
               if grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
                   grid[tl.y][tl.x] = '0'
               }
               if grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
                   grid[bl.y][bl.x] = '0'
               }
               grid[tr.y][tr.x] = '1'
               grid[br.y][br.x] = '1'
               ops = append(ops, op)
           case 2:
               // skip perfect col2
               if grid[tr.y][tr.x] == '1' && grid[br.y][br.x] == '1' {
                   return
               }
               var op []Pair
               if grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
               }
               if grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
               }
               if grid[tr.y][tr.x] == '1' {
                   op = append(op, tr)
               }
               if grid[br.y][br.x] == '1' {
                   op = append(op, br)
               }
               found := false
               if !found && grid[tr.y][tr.x] == '0' {
                   op = append(op, tr)
                   found = true
               }
               if !found && grid[br.y][br.x] == '0' {
                   op = append(op, br)
                   found = true
               }
               if !found && grid[tl.y][tl.x] == '0' {
                   op = append(op, tl)
                   found = true
               }
               if !found && grid[bl.y][bl.x] == '0' {
                   op = append(op, bl)
                   found = true
               }
               for _, p := range op {
                   if grid[p.y][p.x] == '1' {
                       grid[p.y][p.x] = '0'
                   } else {
                       grid[p.y][p.x] = '1'
                   }
               }
               ops = append(ops, op)
           case 3:
               op := make([]Pair, 0, 3)
               if grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
                   grid[tl.y][tl.x] = '0'
               }
               if grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
                   grid[bl.y][bl.x] = '0'
               }
               if grid[tr.y][tr.x] == '1' {
                   op = append(op, tr)
                   grid[tr.y][tr.x] = '0'
               }
               if grid[br.y][br.x] == '1' {
                   op = append(op, br)
                   grid[br.y][br.x] = '0'
               }
               ops = append(ops, op)
           case 4:
               op := []Pair{tl, bl, tr}
               grid[tl.y][tl.x] = '0'
               grid[bl.y][bl.x] = '0'
               grid[tr.y][tr.x] = '0'
               ops = append(ops, op)
           }
       }

       var handlefull func(c1 [2]Pair, c2 [2]Pair)
       handlefull = func(c1 [2]Pair, c2 [2]Pair) {
           tl, bl := c1[0], c1[1]
           tr, br := c2[0], c2[1]
           cnt := 0
           if grid[tl.y][tl.x] == '1' {
               cnt++
           }
           if grid[bl.y][bl.x] == '1' {
               cnt++
           }
           if grid[tr.y][tr.x] == '1' {
               cnt++
           }
           if grid[br.y][br.x] == '1' {
               cnt++
           }
           switch cnt {
           case 1:
               op := make([]Pair, 0, 3)
               if grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
               }
               if grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
               }
               if grid[br.y][br.x] == '1' {
                   op = append(op, br)
               }
               if grid[tr.y][tr.x] == '1' {
                   op = append(op, tr)
               }
               cnt0 := 0
               if cnt0 < 2 && grid[tl.y][tl.x] == '0' {
                   op = append(op, tl)
                   cnt0++
               }
               if cnt0 < 2 && grid[bl.y][bl.x] == '0' {
                   op = append(op, bl)
                   cnt0++
               }
               if cnt0 < 2 && grid[tr.y][tr.x] == '0' {
                   op = append(op, tr)
                   cnt0++
               }
               if cnt0 < 2 && grid[br.y][br.x] == '0' {
                   op = append(op, br)
                   cnt0++
               }
               ops = append(ops, op)
               for _, p := range op {
                   if grid[p.y][p.x] == '1' {
                       grid[p.y][p.x] = '0'
                   } else {
                       grid[p.y][p.x] = '1'
                   }
               }
               handlefull(c1, c2)
           case 2:
               op := make([]Pair, 0, 3)
               cnt1 := 0
               if cnt1 < 1 && grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
                   cnt1++
               }
               if cnt1 < 1 && grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
                   cnt1++
               }
               if cnt1 < 1 && grid[br.y][br.x] == '1' {
                   op = append(op, br)
                   cnt1++
               }
               if cnt1 < 1 && grid[tr.y][tr.x] == '1' {
                   op = append(op, tr)
                   cnt1++
               }
               cnt0 := 0
               if cnt0 < 2 && grid[tl.y][tl.x] == '0' {
                   op = append(op, tl)
                   cnt0++
               }
               if cnt0 < 2 && grid[bl.y][bl.x] == '0' {
                   op = append(op, bl)
                   cnt0++
               }
               if cnt0 < 2 && grid[tr.y][tr.x] == '0' {
                   op = append(op, tr)
                   cnt0++
               }
               if cnt0 < 2 && grid[br.y][br.x] == '0' {
                   op = append(op, br)
                   cnt0++
               }
               ops = append(ops, op)
               for _, p := range op {
                   if grid[p.y][p.x] == '1' {
                       grid[p.y][p.x] = '0'
                   } else {
                       grid[p.y][p.x] = '1'
                   }
               }
               handlefull(c1, c2)
           case 3:
               op := make([]Pair, 0, 3)
               if grid[tl.y][tl.x] == '1' {
                   op = append(op, tl)
                   grid[tl.y][tl.x] = '0'
               }
               if grid[bl.y][bl.x] == '1' {
                   op = append(op, bl)
                   grid[bl.y][bl.x] = '0'
               }
               if grid[tr.y][tr.x] == '1' {
                   op = append(op, tr)
                   grid[tr.y][tr.x] = '0'
               }
               if grid[br.y][br.x] == '1' {
                   op = append(op, br)
                   grid[br.y][br.x] = '0'
               }
               ops = append(ops, op)
           case 4:
               op := []Pair{tl, bl, tr}
               grid[tl.y][tl.x] = '0'
               grid[bl.y][bl.x] = '0'
               grid[tr.y][tr.x] = '0'
               ops = append(ops, op)
               handlefull(c1, c2)
           }
       }

       // cover grid except last 2x2
       if H > 2 || W > 2 {
           for j := 0; j+1 < H; j += 2 {
               for i := 0; i < W-2; i++ {
                   c1 := [2]Pair{{j, i}, {j + 1, i}}
                   c2 := [2]Pair{{j, i + 1}, {j + 1, i + 1}}
                   handleit(c1, c2)
               }
           }
           if H%2 == 1 {
               j := H - 2
               for i := 0; i < W-2; i++ {
                   c1 := [2]Pair{{j, i}, {j + 1, i}}
                   c2 := [2]Pair{{j, i + 1}, {j + 1, i + 1}}
                   handleit(c1, c2)
               }
           }
           if H > 2 {
               for j := 0; j+1 < H; j++ {
                   i := W - 1
                   c1 := [2]Pair{{j, i}, {j, i - 1}}
                   c2 := [2]Pair{{j + 1, i}, {j + 1, i - 1}}
                   handleit(c1, c2)
               }
           }
       }
       // last 2x2
       i := W - 2
       j := H - 2
       c1 := [2]Pair{{j, i}, {j + 1, i}}
       c2 := [2]Pair{{j, i + 1}, {j + 1, i + 1}}
       handlefull(c1, c2)

       // output
       fmt.Fprintln(writer, len(ops))
       for _, op := range ops {
           for _, p := range op {
               fmt.Fprintf(writer, "%d %d ", p.y+1, p.x+1)
           }
           fmt.Fprintln(writer)
       }
   }
}
