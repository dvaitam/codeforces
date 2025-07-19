package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var N, M int
   if _, err := fmt.Fscan(reader, &N, &M); err != nil {
       return
   }
   table := make([][]byte, N)
   for i := 0; i < N; i++ {
       var s string
       fmt.Fscan(reader, &s)
       table[i] = []byte(s)
   }

   dx := [4]int{1, -1, 0, 0}
   dy := [4]int{0, 0, -1, 1}

   for i := 0; i < N; i++ {
       for j := 0; j < M; j++ {
           if table[i][j] != '.' {
               continue
           }
           found := false
           bk := -1
           for k := 0; k < 4; k++ {
               ni := i + dy[k]
               nj := j + dx[k]
               if ni < 0 || ni >= N || nj < 0 || nj >= M {
                   continue
               }
               if table[ni][nj] == '#' {
                   continue
               }
               if table[ni][nj] != '.' {
                   bk = k
               } else {
                   dig := (i%3)*3 + j%3
                   table[i][j] = byte('0' + dig)
                   table[ni][nj] = byte('0' + dig)
                   found = true
                   break
               }
           }
           if !found && bk != -1 {
               ni := i + dy[bk]
               nj := j + dx[bk]
               table[i][j] = table[ni][nj]
               found = true
           }
           if !found {
               writer.WriteString("-1")
               return
           }
       }
   }

   for i := 0; i < N; i++ {
       writer.Write(table[i])
       writer.WriteByte('\n')
   }
}
