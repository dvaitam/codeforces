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

   var n, m, x int
   if _, err := fmt.Fscan(reader, &n, &m, &x); err != nil {
       return
   }
   // positions of letters and shifts
   type pos struct{ r, c int }
   letters := make([][]pos, 26)
   var shifts []pos
   // read keyboard
   for i := 0; i < n; i++ {
       var row string
       fmt.Fscan(reader, &row)
       for j := 0; j < m; j++ {
           ch := row[j]
           if ch == 'S' {
               shifts = append(shifts, pos{i, j})
           } else if ch >= 'a' && ch <= 'z' {
               letters[ch-'a'] = append(letters[ch-'a'], pos{i, j})
           }
       }
   }
   // precompute for each letter if uppercase can be typed one-handed
   good := make([]bool, 26)
   maxd2 := x * x
   if len(shifts) > 0 {
       for i := 0; i < 26; i++ {
           if len(letters[i]) == 0 {
               continue
           }
           ok := false
           for _, p := range letters[i] {
               for _, s := range shifts {
                   dr := p.r - s.r
                   dc := p.c - s.c
                   if dr*dr+dc*dc <= maxd2 {
                       ok = true
                       break
                   }
               }
               if ok {
                   break
               }
           }
           good[i] = ok
       }
   }
   // read text
   var q int
   fmt.Fscan(reader, &q)
   var text string
   fmt.Fscan(reader, &text)

   count := 0
   for i := 0; i < len(text); i++ {
       ch := text[i]
       if ch >= 'a' && ch <= 'z' {
           idx := ch - 'a'
           if len(letters[idx]) == 0 {
               fmt.Fprintln(writer, -1)
               return
           }
       } else if ch >= 'A' && ch <= 'Z' {
           idx := ch - 'A'
           // lowercase existence
           if len(letters[idx]) == 0 || len(shifts) == 0 {
               fmt.Fprintln(writer, -1)
               return
           }
           if !good[idx] {
               count++
           }
       } else {
           // other chars not expected
       }
   }
   fmt.Fprintln(writer, count)
}
