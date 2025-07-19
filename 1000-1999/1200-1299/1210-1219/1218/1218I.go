package main

import (
   "bufio"
   "fmt"
   "os"
)

var n int
var ini, fin, tmp [][]bool
var mag []bool
var c1 int
type Op struct{ typ, idx int }
var oper []Op

func column(i int) {
   for j := 0; j < n; j++ {
       ini[j][i] = ini[j][i] != mag[j]
   }
}

func row(i int) {
   for j := 0; j < n; j++ {
       ini[i][j] = ini[i][j] != mag[j]
   }
}

func check() {
   for i := 0; i < n; i++ {
       if ini[i][c1] != fin[i][c1] {
           row(i)
           oper = append(oper, Op{0, i})
       }
   }
   for j := 0; j < n; j++ {
       if j == c1 {
           continue
       }
       cnt, cnt2 := 0, 0
       for i := 0; i < n; i++ {
           if (ini[i][j] != mag[i]) == fin[i][j] {
               cnt++
           }
           if ini[i][j] == fin[i][j] {
               cnt2++
           }
       }
       if cnt == n {
           oper = append(oper, Op{1, j})
       } else if cnt2 != n {
           return
       }
   }
   fmt.Println(len(oper))
   for _, op := range oper {
       if op.typ == 1 {
           fmt.Printf("col %d\n", op.idx)
       } else {
           fmt.Printf("row %d\n", op.idx)
       }
   }
   os.Exit(0)
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   fmt.Fscan(reader, &n)
   ini = make([][]bool, n)
   tmp = make([][]bool, n)
   fin = make([][]bool, n)
   mag = make([]bool, n)
   for i := 0; i < n; i++ {
       ini[i] = make([]bool, n)
       tmp[i] = make([]bool, n)
       fin[i] = make([]bool, n)
   }
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j := 0; j < n; j++ {
           ini[i][j] = s[j] == '1'
           tmp[i][j] = ini[i][j]
       }
   }
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       for j := 0; j < n; j++ {
           fin[i][j] = s[j] == '1'
       }
   }
   var s string
   fmt.Fscan(reader, &s)
   for j := 0; j < n; j++ {
       mag[j] = s[j] == '1'
       if mag[j] {
           c1 = j
       }
   }
   check()
   // try flipping column c1 first
   oper = nil
   for i := 0; i < n; i++ {
       copyRow := tmp[i]
       for j := 0; j < n; j++ {
           ini[i][j] = copyRow[j]
       }
   }
   column(c1)
   oper = append(oper, Op{1, c1})
   check()
   fmt.Println(-1)
}
