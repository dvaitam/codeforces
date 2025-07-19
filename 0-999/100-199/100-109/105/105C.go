package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var n int
   fmt.Fscan(reader, &n)
   name := make([]string, n)
   kind := make([]string, n)
   a := make([]int, n)
   b := make([]int, n)
   c := make([]int, n)
   s := make([]int, n)
   totalSlots := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &name[i], &kind[i], &a[i], &b[i], &c[i], &s[i])
       totalSlots += s[i]
   }
   var k int
   fmt.Fscan(reader, &k)
   if totalSlots == k {
       // Exact assignment case
       p := make([]int, n)
       names := make([]string, k)
       belongs := make([]string, k)
       kinds := make([]string, k)
       vals := make([]int, k)
       for i := 0; i < k; i++ {
           fmt.Fscan(reader, &names[i], &kinds[i], &vals[i], &belongs[i])
           for j := 0; j < n; j++ {
               if name[j] == belongs[i] {
                   p[j]++
                   switch kinds[i][0] {
                   case 'g':
                       a[j] += vals[i]
                   case 's':
                       b[j] += vals[i]
                   case 'p':
                       c[j] += vals[i]
                   }
                   break
               }
           }
       }
       // Find winners
       A, B, C := -1, -1, -1
       AA, BB, CC := 0, 0, 0
       for i := 0; i < n; i++ {
           switch kind[i][0] {
           case 'w':
               if a[i] > A { A, AA = a[i], i }
           case 'a':
               if b[i] > B { B, BB = b[i], i }
           case 'o':
               if c[i] > C { C, CC = c[i], i }
           }
       }
       // Output for each
       fmt.Fprint(writer, name[AA], " ", p[AA])
       for i := 0; i < k; i++ {
           if belongs[i] == name[AA] {
               fmt.Fprint(writer, " ", names[i])
           }
       }
       fmt.Fprintln(writer)
       fmt.Fprint(writer, name[BB], " ", p[BB])
       for i := 0; i < k; i++ {
           if belongs[i] == name[BB] {
               fmt.Fprint(writer, " ", names[i])
           }
       }
       fmt.Fprintln(writer)
       fmt.Fprint(writer, name[CC], " ", p[CC])
       for i := 0; i < k; i++ {
           if belongs[i] == name[CC] {
               fmt.Fprint(writer, " ", names[i])
           }
       }
       fmt.Fprintln(writer)
   } else {
       // Redistribution case
       type pair struct { first int; second string }
       var gl, se, ph []pair
       for i := 0; i < k; i++ {
           var nm, kd, bl string
           var v int
           fmt.Fscan(reader, &nm, &kd, &v, &bl)
           switch kd[0] {
           case 'g': gl = append(gl, pair{v, nm})
           case 's': se = append(se, pair{v, nm})
           case 'p': ph = append(ph, pair{v, nm})
           }
       }
       sort.Slice(gl, func(i, j int) bool { return gl[i].first < gl[j].first })
       sort.Slice(se, func(i, j int) bool { return se[i].first < se[j].first })
       sort.Slice(ph, func(i, j int) bool { return ph[i].first < ph[j].first })
       for i := 0; i < n; i++ {
           switch kind[i][0] {
           case 'w':
               for j := 0; j < s[i] && j < len(gl); j++ {
                   a[i] += gl[len(gl)-1-j].first
               }
           case 'a':
               for j := 0; j < s[i] && j < len(se); j++ {
                   b[i] += se[len(se)-1-j].first
               }
           case 'o':
               for j := 0; j < s[i] && j < len(ph); j++ {
                   c[i] += ph[len(ph)-1-j].first
               }
           }
       }
       A, B, C := -1, -1, -1
       AA, BB, CC := 0, 0, 0
       for i := 0; i < n; i++ {
           switch kind[i][0] {
           case 'w': if a[i] > A { A, AA = a[i], i }
           case 'a': if b[i] > B { B, BB = b[i], i }
           case 'o': if c[i] > C { C, CC = c[i], i }
           }
       }
       var As, Bs, Cs []string
       for j := 0; j < s[AA] && j < len(gl); j++ {
           As = append(As, gl[len(gl)-1-j].second)
       }
       for j := 0; j < s[BB] && j < len(se); j++ {
           Bs = append(Bs, se[len(se)-1-j].second)
       }
       for j := 0; j < s[CC] && j < len(ph); j++ {
           Cs = append(Cs, ph[len(ph)-1-j].second)
       }
       for j := 0; j < len(gl)-s[AA]; j++ {
           if len(Bs) < s[BB] {
               Bs = append(Bs, gl[j].second)
           } else if len(Cs) < s[CC] {
               Cs = append(Cs, gl[j].second)
           }
       }
       for j := 0; j < len(se)-s[BB]; j++ {
           if len(As) < s[AA] {
               As = append(As, se[j].second)
           } else if len(Cs) < s[CC] {
               Cs = append(Cs, se[j].second)
           }
       }
       for j := 0; j < len(ph)-s[CC]; j++ {
           if len(Bs) < s[BB] {
               Bs = append(Bs, ph[j].second)
           } else if len(As) < s[AA] {
               As = append(As, ph[j].second)
           }
       }
       fmt.Fprint(writer, name[AA], " ", len(As))
       for _, v := range As { fmt.Fprint(writer, " ", v) }
       fmt.Fprintln(writer)
       fmt.Fprint(writer, name[BB], " ", len(Bs))
       for _, v := range Bs { fmt.Fprint(writer, " ", v) }
       fmt.Fprintln(writer)
       fmt.Fprint(writer, name[CC], " ", len(Cs))
       for _, v := range Cs { fmt.Fprint(writer, " ", v) }
       fmt.Fprintln(writer)
   }
}
