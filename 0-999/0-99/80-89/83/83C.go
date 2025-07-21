package main

import (
   "bufio"
   "fmt"
   "os"
   "unicode"
   "math/bits"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m, k int
   if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
       return
   }
   board := make([][]rune, n)
   var sr, sc, tr, tc int
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       board[i] = []rune(line)
       for j, ch := range board[i] {
           if ch == 'S' {
               sr, sc = i, j
           } else if ch == 'T' {
               tr, tc = i, j
           }
       }
   }
   // BFS with state (r,c,mask)
   type void struct{}
   dist := make(map[int64]int)
   // encode state as mask<<12 | r<<6 | c
   encode := func(r, c, mask int) int64 {
       return (int64(mask) << 12) | (int64(r) << 6) | int64(c)
   }
   decode := func(key int64) (r, c, mask int) {
       c = int(key & 0x3F)
       r = int((key >> 6) & 0x3F)
       mask = int(key >> 12)
       return
   }
   startKey := encode(sr, sc, 0)
   dist[startKey] = 0
   queue := []int64{startKey}
   head := 0
   L := -1
   // BFS level by level
   depth := 0
   found := false
   for head < len(queue) && !found {
       sz := len(queue) - head
       for i := 0; i < sz; i++ {
           curr := queue[head]
           head++
           r, c, mask := decode(curr)
           if r == tr && c == tc {
               L = depth
               found = true
               continue
           }
           // expand if not yet found
           for _, d := range [4][2]int{{-1,0},{1,0},{0,-1},{0,1}} {
               nr, nc := r + d[0], c + d[1]
               if nr < 0 || nr >= n || nc < 0 || nc >= m {
                   continue
               }
               ch := board[nr][nc]
               // cannot revisit start
               if ch == 'S' {
                   continue
               }
               newMask := mask
               if unicode.IsLower(ch) {
                   bit := 1 << (ch - 'a')
                   newMask = mask | bit
                   if bits.OnesCount(uint(newMask)) > k {
                       continue
                   }
               }
               key2 := encode(nr, nc, newMask)
               if _, ok := dist[key2]; ok {
                   continue
               }
               dist[key2] = depth + 1
               queue = append(queue, key2)
           }
       }
       if !found {
           depth++
       }
   }
   if L < 0 {
       fmt.Println(-1)
       return
   }
   // reconstruct lexicographically minimal path
   // initial frontier
   frontier := []int64{startKey}
   result := make([]rune, 0, L)
   for step := 1; step <= L; step++ {
       var best rune = 0
       nextMap := make(map[int64]void)
       for _, curr := range frontier {
           r, c, mask := decode(curr)
           for _, d := range [4][2]int{{-1,0},{1,0},{0,-1},{0,1}} {
               nr, nc := r + d[0], c + d[1]
               if nr < 0 || nr >= n || nc < 0 || nc >= m {
                   continue
               }
               ch := board[nr][nc]
               newMask := mask
               if unicode.IsLower(ch) {
                   bit := 1 << (ch - 'a')
                   newMask = mask | bit
                   if bits.OnesCount(uint(newMask)) > k {
                       continue
                   }
               }
               key2 := encode(nr, nc, newMask)
               if d2, ok := dist[key2]; !ok || d2 != step {
                   continue
               }
               // determine if valid move for this step
               if step < L {
                   if nr == tr && nc == tc {
                       continue
                   }
                   if !unicode.IsLower(ch) {
                       continue
                   }
                   if best == 0 || ch < best {
                       best = ch
                       nextMap = make(map[int64]void)
                       nextMap[key2] = void{}
                   } else if ch == best {
                       nextMap[key2] = void{}
                   }
               } else {
                   // last step, must go to T
                   if nr == tr && nc == tc {
                       nextMap[key2] = void{}
                   }
               }
           }
       }
       if step < L {
           result = append(result, best)
       }
       // update frontier
       frontier = frontier[:0]
       for key2 := range nextMap {
           frontier = append(frontier, key2)
       }
   }
   // print result
   if len(result) > 0 {
       fmt.Println(string(result))
   } else {
       fmt.Println()
   }
}
