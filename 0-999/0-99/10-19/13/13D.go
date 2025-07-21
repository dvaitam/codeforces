package main

import (
   "bufio"
   "fmt"
   "os"
)

type Point struct{ x, y int64 }

func main() {
   in := bufio.NewReader(os.Stdin)
   var N, M int
   if _, err := fmt.Fscan(in, &N, &M); err != nil {
       return
   }
   reds := make([]Point, N)
   for i := 0; i < N; i++ {
       fmt.Fscan(in, &reds[i].x, &reds[i].y)
   }
   blues := make([]Point, M)
   for i := 0; i < M; i++ {
       fmt.Fscan(in, &blues[i].x, &blues[i].y)
   }
   if N < 3 {
       fmt.Println(0)
       return
   }
   // Special case: no blue points
   if M == 0 {
       // C(N,3)
       n := int64(N)
       fmt.Println(n * (n - 1) * (n - 2) / 6)
       return
   }
   // bitsets: bs[i][j] bitset of blues to left of red[i]->red[j]
   W := (M + 63) / 64
   bs := make([][][]uint64, N)
   for i := 0; i < N; i++ {
       bs[i] = make([][]uint64, N)
       for j := 0; j < N; j++ {
           if i != j {
               bs[i][j] = make([]uint64, W)
           }
       }
   }
   // fill bitsets
   for b := 0; b < M; b++ {
       bx, by := blues[b].x, blues[b].y
       wi := b >> 6
       mask := uint64(1) << (uint(b) & 63)
       for i := 0; i < N; i++ {
           xi, yi := reds[i].x, reds[i].y
           for j := 0; j < N; j++ {
               if i == j {
                   continue
               }
               xj, yj := reds[j].x, reds[j].y
               // cross of (j-i) and (b-i): >0 means left
               if (xj-xi)*(by-yi) - (yj-yi)*(bx-xi) > 0 {
                   bs[i][j][wi] |= mask
               }
           }
       }
   }
   // count triangles without blues inside
   var cnt int64
   // iterate triples
   for i := 0; i < N; i++ {
       for j := i + 1; j < N; j++ {
           xi, yi := reds[i].x, reds[i].y
           xj, yj := reds[j].x, reds[j].y
           for k := j + 1; k < N; k++ {
               xk, yk := reds[k].x, reds[k].y
               // orientation of i,j,k
               cross := (xj-xi)*(yk-yi) - (yj-yi)*(xk-xi)
               var a, b1, c int
               if cross > 0 {
                   a, b1, c = i, j, k
               } else {
                   a, b1, c = i, k, j
               }
               // check bitset intersection
               ok := true
               bsab := bs[a][b1]
               bsc := bs[b1][c]
               bsca := bs[c][a]
               for w := 0; w < W; w++ {
                   if bsab[w]&bsc[w]&bsca[w] != 0 {
                       ok = false
                       break
                   }
               }
               if ok {
                   cnt++
               }
           }
       }
   }
   fmt.Println(cnt)
}
