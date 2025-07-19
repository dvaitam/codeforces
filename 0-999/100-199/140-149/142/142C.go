package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, m   int
   ans     int
   rem     int
   p       [9][9]byte
   ansp    [9][9]byte
   a       = [3][12]byte{
       {'A','A','A','.','A','.','A','.','.','.','.','A'},
       {'.','A','.','.','A','.','A','A','A','A','A','A'},
       {'.','A','.','A','A','A','A','.','.','.','.','A'},
   }
)

func copyAns() {
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           ansp[i][j] = p[i][j]
       }
   }
}

func max(x, y int) int {
   if x > y {
       return x
   }
   return y
}

func saya(x, y, move int) {
   if x >= n-2 {
       ans = max(ans, move)
       if move == ans {
           copyAns()
       }
       return
   }
   if y >= m-2 {
       if p[x][y] == '.' {
           rem--
       }
       if p[x][y+1] == '.' {
           rem--
       }
       saya(x+1, 0, move)
       if p[x][y] == '.' {
           rem++
       }
       if p[x][y+1] == '.' {
           rem++
       }
       return
   }
   if rem/5 <= ans-move {
       return
   }
   if p[x][y] == '.' {
       rem--
   }
   // try place each of 4 orientations
   for d := 0; d < 12; d += 3 {
       flag := false
       for i := 0; i < 3 && !flag; i++ {
           for j := 0; j < 3; j++ {
               if a[i][d+j] == 'A' && p[x+i][y+j] != '.' {
                   flag = true
                   break
               }
           }
       }
       if !flag {
           for i := 0; i < 3; i++ {
               for j := 0; j < 3; j++ {
                   if a[i][d+j] == 'A' {
                       p[x+i][y+j] = 'A' + byte(move)
                   }
               }
           }
           rem -= 5
           saya(x, y+1, move+1)
           rem += 5
           for i := 0; i < 3; i++ {
               for j := 0; j < 3; j++ {
                   if a[i][d+j] == 'A' {
                       p[x+i][y+j] = '.'
                   }
               }
           }
       }
   }
   saya(x, y+1, move)
   if p[x][y] == '.' {
       rem++
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n, &m)
   // initialize
   for i := 0; i < n; i++ {
       for j := 0; j < m; j++ {
           p[i][j] = '.'
           ansp[i][j] = '.'
       }
   }
   rem = n * m
   saya(0, 0, 0)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, ans)
   for i := 0; i < n; i++ {
       writer.Write(ansp[i][:m])
       writer.WriteByte('\n')
   }
}
