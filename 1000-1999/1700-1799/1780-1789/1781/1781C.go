package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type fairy struct {
   pos, val int
}

func solve(in *bufio.Reader, out *bufio.Writer) {
   var t int
   fmt.Fscan(in, &t)
   for t > 0 {
       t--
       var n int
       fmt.Fscan(in, &n)
       var s string
       fmt.Fscan(in, &s)
       sum := make([]int, 26)
       for i := 0; i < n; i++ {
           sum[s[i]-'a']++
       }
       minn := n
       var ans int
       // find optimal group size
       for num := 1; num <= n; num++ {
           if n%num != 0 || n/num > 26 {
               continue
           }
           grp := n / num
           temp := make([]fairy, 0, 26)
           for i := 0; i < 26; i++ {
               if sum[i] > 0 {
                   temp = append(temp, fairy{pos: i, val: sum[i]})
               }
           }
           sort.Slice(temp, func(i, j int) bool { return temp[i].val > temp[j].val })
           cnt := 0
           if len(temp) < grp {
               for _, f := range temp {
                   if f.val > num {
                       cnt += f.val - num
                   } else {
                       break
                   }
               }
           } else {
               for i, f := range temp {
                   if i < grp && f.val > num {
                       cnt += f.val - num
                   }
                   if i >= grp {
                       cnt += f.val
                   }
               }
           }
           if cnt < minn {
               minn = cnt
               ans = num
           }
       }
       fmt.Fprintln(out, minn)
       grp := n / ans
       temp := make([]fairy, 0, 26)
       for i := 0; i < 26; i++ {
           if sum[i] > 0 {
               temp = append(temp, fairy{pos: i, val: sum[i]})
           }
       }
       sort.Slice(temp, func(i, j int) bool { return temp[i].val > temp[j].val })
       good := make([]int, 26)
       for i, f := range temp {
           if i < grp && f.val <= ans {
               good[f.pos] = 1
           } else {
               good[f.pos] = 2
           }
       }
       bs := []byte(s)
       // first pass: replace overrepresented or inappropriate letters
       for i := 0; i < n; i++ {
           ch := bs[i] - 'a'
           if good[ch] == 2 && ((sum[ch] > 0 && sum[ch] < ans) || sum[ch] > ans) {
               for j := 0; j < 26; j++ {
                   if byte(j) == ch {
                       continue
                   }
                   if good[j] == 1 && sum[j] < ans {
                       sum[j]++
                       sum[ch]--
                       bs[i] = byte('a' + j)
                       break
                   }
               }
           }
       }
       // second pass: fill gaps from unused letters
       for i := 0; i < n; i++ {
           ch := bs[i] - 'a'
           if good[ch] == 2 && sum[ch] > ans {
               for j := 0; j < 26; j++ {
                   if good[j] == 0 && sum[j] < ans {
                       sum[j]++
                       sum[ch]--
                       bs[i] = byte('a' + j)
                       break
                   }
               }
           }
       }
       out.Write(bs)
       out.WriteByte('\n')
   }
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   solve(in, out)
}
