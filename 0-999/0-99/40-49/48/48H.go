package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, a, b, c int
   if _, err := fmt.Fscan(in, &n, &m, &a, &b, &c); err != nil {
       return
   }
   // Prepare output pixel grid 2n x 2m
   H := 2 * n
   W := 2 * m
   grid := make([][]byte, H)
   for i := 0; i < H; i++ {
       grid[i] = make([]byte, W)
   }
   // Define tile patterns
   type tile struct {
       pat [2][2]byte
       cnt *int
   }
   // Mutable counts
   cntB := a
   cntW := b
   cntM := c
   patterns := []tile{
       // pure black
       { [2][2]byte{{'B','B'},{'B','B'}}, &cntB },
       // pure white
       { [2][2]byte{{'W','W'},{'W','W'}}, &cntW },
       // mixed: top black, bottom white
       { [2][2]byte{{'B','B'},{'W','W'}}, &cntM },
       // mixed: top white, bottom black
       { [2][2]byte{{'W','W'},{'B','B'}}, &cntM },
       // mixed: left black, right white
       { [2][2]byte{{'B','W'},{'B','W'}}, &cntM },
       // mixed: left white, right black
       { [2][2]byte{{'W','B'},{'W','B'}}, &cntM },
   }
   // iterate cells row-major
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           // determine required top and left pixels
           var top [2]byte
           var left [2]byte
           hasTop := false
           hasLeft := false
           // top border
           if i > 0 {
               hasTop = true
               r := 2*i
               c0 := 2*j
               top[0] = grid[r-1][c0]
               top[1] = grid[r-1][c0+1]
           }
           if j > 0 {
               hasLeft = true
               r0 := 2*i
               c := 2*j
               left[0] = grid[r0][c-1]
               left[1] = grid[r0+1][c-1]
           }
           // pick a matching tile
           placed := false
           for _, t := range patterns {
               if *t.cnt <= 0 {
                   continue
               }
               ok := true
               if hasTop {
                   if t.pat[0][0] != top[0] || t.pat[0][1] != top[1] {
                       ok = false
                   }
               }
               if ok && hasLeft {
                   if t.pat[0][0] != left[0] || t.pat[1][0] != left[1] {
                       ok = false
                   }
               }
               if !ok {
                   continue
               }
               // place
               r0, c0 := 2*i, 2*j
               for di := 0; di < 2; di++ {
                   for dj := 0; dj < 2; dj++ {
                       grid[r0+di][c0+dj] = t.pat[di][dj]
                   }
               }
               *t.cnt--
               placed = true
               break
           }
           if !placed {
               // fallback: place any mixed
               for k := 2; k < len(patterns); k++ {
                   t := &patterns[k]
                   if *t.cnt > 0 {
                       r0, c0 := 2*i, 2*j
                       for di := 0; di < 2; di++ {
                           for dj := 0; dj < 2; dj++ {
                               grid[r0+di][c0+dj] = patterns[k].pat[di][dj]
                           }
                       }
                       *t.cnt--
                       placed = true
                       break
                   }
               }
           }
           if !placed {
               // no tile: error
               fmt.Fprintf(os.Stderr, "Error at cell %d,%d\n", i, j)
           }
       }
   }
   // output grid
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   for i := 0; i < H; i++ {
       out.Write(grid[i])
       out.WriteByte('\n')
   }
}
