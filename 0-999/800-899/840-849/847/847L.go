package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "strconv"
   "strings"
)

type pair struct{first, second int}

var (
   n    int
   v    [][][]int
   cur  []int
   stmp int
   ans  []pair
)

func goRec(vv []int, prt int) {
   stmp++
   for _, son := range vv {
       if son >= 0 && son < len(cur) {
           cur[son] = stmp
       }
   }
   for _, son := range vv {
       ff := -1
       for i, nbrList := range v[son] {
           foundPrt := false
           same := false
           cc := 0
           for _, j := range nbrList {
               if j >= 0 && j < len(cur) && cur[j] == stmp {
                   same = true
                   break
               } else {
                   cc++
               }
               if j == prt {
                   foundPrt = true
               }
           }
           if !same && foundPrt && cc+len(vv) == n {
               ff = i
               break
           }
       }
       if ff == -1 {
           fmt.Println(-1)
           os.Exit(0)
       }
       ans = append(ans, pair{son, prt})
       for i, nbrList := range v[son] {
           if i == ff {
               continue
           }
           goRec(nbrList, son)
       }
       return
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read n
   if _, err := fmt.Fscanln(reader, &n); err != nil {
       return
   }
   v = make([][][]int, n+1)
   cur = make([]int, n+1)
   ans = make([]pair, 0)
   // read each line
   for i := 1; i <= n; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && err != io.EOF {
           fmt.Fprintln(os.Stderr, "read error:", err)
           os.Exit(1)
       }
       line = strings.TrimRight(line, "\r\n")
       ptr := 0
       L := len(line)
       for ptr < L {
           // parse tot
           j := ptr
           for j < L && line[j] != ':' {
               j++
           }
           tot, _ := strconv.Atoi(line[ptr:j])
           ptr = j + 1
           tmp := make([]int, 0, tot)
           // parse tot values
           for k := 0; k < tot; k++ {
               j = ptr
               for j < L && line[j] != ',' && line[j] != '-' {
                   j++
               }
               val, _ := strconv.Atoi(line[ptr:j])
               tmp = append(tmp, val)
               ptr = j + 1
           }
           v[i] = append(v[i], tmp)
       }
   }
   // solve
   if len(v) > 1 {
       for _, firstList := range v[1] {
           goRec(firstList, 1)
       }
   }
   // output answer
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fprintln(out, len(ans))
   for _, p := range ans {
       fmt.Fprintln(out, p.first, p.second)
   }
}
