package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func min(a, b int) int { if a < b { return a }; return b }

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // fast integer reader
   var readInt func() (int, error)
   readInt = func() (int, error) {
       // skip non-digit
       b, err := reader.ReadByte()
       if err != nil {
           return 0, err
       }
       for (b < '0' || b > '9') && b != '-' {
           b, err = reader.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       neg := false
       if b == '-' {
           neg = true
           b, err = reader.ReadByte()
           if err != nil {
               return 0, err
           }
       }
       v := 0
       for b >= '0' && b <= '9' {
           v = v*10 + int(b - '0')
           b, err = reader.ReadByte()
           if err != nil {
               break
           }
       }
       if neg {
           v = -v
       }
       return v, nil
   }
   // read n, m, k
   n64, _ := readInt()
   m64, _ := readInt()
   k64, _ := readInt()
   n, m, k := n64, m64, k64
   total := n * m
   grid := make([]int, total)
   colorPos := make([][]int, k)
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           x := readInt() - 1
           idx := i*m + j
           grid[idx] = x
           colorPos[x] = append(colorPos[x], idx)
       }
   }
   // distColors[c][idx]
   const INF = 1<<30
   distColors := make([][]int, k)
   // BFS helpers
   dirs := [4]int{-1, 1, 0, 0}
   dirc := [4]int{0, 0, -1, 1}
   for c := 0; c < k; c++ {
       dist := make([]int, total)
       for i := range dist {
           dist[i] = INF
       }
       bestColor := make([]bool, k)
       // queue
       q := make([]int, 0, total)
       // init
       for _, p := range colorPos[c] {
           dist[p] = 0
           q = append(q, p)
       }
       bestColor[c] = false
       // BFS
       for qi := 0; qi < len(q); qi++ {
           u := q[qi]
           du := dist[u]
           // teleport by color of u
           cu := grid[u]
           if !bestColor[cu] {
               for _, v := range colorPos[cu] {
                   if dist[v] > du+1 {
                       dist[v] = du + 1
                       q = append(q, v)
                   }
               }
               bestColor[cu] = true
           }
           // adjacent
           ui := u / m
           uj := u % m
           for d := 0; d < 4; d++ {
               ni := ui + dirs[d]
               nj := uj + dirc[d]
               if ni >= 0 && ni < n && nj >= 0 && nj < m {
                   v := ni*m + nj
                   if dist[v] > du+1 {
                       dist[v] = du + 1
                       q = append(q, v)
                   }
               }
           }
       }
       distColors[c] = dist
   }
   // queries
   qn := readInt()
   for qi := 0; qi < qn; qi++ {
       r1 := readInt() - 1
       c1 := readInt() - 1
       r2 := readInt() - 1
       c2 := readInt() - 1
       idx1 := r1*m + c1
       idx2 := r2*m + c2
       ans := abs(r1-r2) + abs(c1-c2)
       // via teleports
       for c := 0; c < k; c++ {
           d := distColors[c][idx1] + distColors[c][idx2] + 1
           if d < ans {
               ans = d
           }
       }
       fmt.Fprintln(writer, ans)
   }
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}
