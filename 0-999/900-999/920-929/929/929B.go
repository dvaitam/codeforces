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

   var n, k int
   fmt.Fscan(reader, &n, &k)
   s := make([][]byte, n)
   for i := 0; i < n; i++ {
       var line string
       fmt.Fscan(reader, &line)
       s[i] = []byte(line)
   }
   kLeft := k

   // First pass: avoid seating next to status passengers
   for i := 0; i < n; i++ {
       if kLeft == 0 {
           break
       }
       if s[i][0] == '.' && s[i][1] != 'S' {
           s[i][0] = 'x'
           kLeft--
       }
       if kLeft == 0 {
           break
       }
       if s[i][11] == '.' && s[i][10] != 'S' {
           s[i][11] = 'x'
           kLeft--
       }
       if kLeft == 0 {
           break
       }
       for j := 1; j < 11 && kLeft > 0; j++ {
           if s[i][j] == '.' && s[i][j-1] != 'S' && s[i][j+1] != 'S' {
               s[i][j] = 'x'
               kLeft--
           }
       }
   }

   // Second pass: avoid creating new worst-case adjacencies
   for i := 0; i < n; i++ {
       if kLeft == 0 {
           break
       }
       if s[i][0] == '.' {
           s[i][0] = 'x'
           kLeft--
       }
       if kLeft == 0 {
           break
       }
       if s[i][11] == '.' {
           s[i][11] = 'x'
           kLeft--
       }
       if kLeft == 0 {
           break
       }
       for j := 1; j < 11 && kLeft > 0; j++ {
           if s[i][j] == '.' && !(s[i][j-1] == 'S' && s[i][j+1] == 'S') {
               s[i][j] = 'x'
               kLeft--
           }
       }
   }

   // Third pass: fill any remaining free seats
   for i := 0; i < n && kLeft > 0; i++ {
       for j := 0; j < 12 && kLeft > 0; j++ {
           if s[i][j] == '.' {
               s[i][j] = 'x'
               kLeft--
           }
       }
   }

   // Count adjacent neighbors for status passengers
   ans := 0
   for i := 0; i < n; i++ {
       if s[i][0] == 'S' {
           if s[i][1] != '.' && s[i][1] != '-' {
               ans++
           }
       }
       if s[i][11] == 'S' {
           if s[i][10] != '.' && s[i][10] != '-' {
               ans++
           }
       }
       for j := 1; j < 11; j++ {
           if s[i][j] == 'S' {
               if s[i][j-1] != '.' && s[i][j-1] != '-' {
                   ans++
               }
               if s[i][j+1] != '.' && s[i][j+1] != '-' {
                   ans++
               }
           }
       }
   }

   fmt.Fprintln(writer, ans)
   for i := 0; i < n; i++ {
       writer.Write(s[i])
       writer.WriteByte('\n')
   }
}
