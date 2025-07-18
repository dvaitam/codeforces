#pragma GCC optimize ("Ofast")
#include<bits/stdc++.h>
using namespace std;
inline int my_getchar(){
  static char buf[1048576];
  static int s = 1048576;
  static int e = 1048576;
  if(s == e && e == 1048576){
    e = fread(buf, 1, 1048576, stdin);
    s = 0;
  }
  if(s == e){
    return EOF;
  }
  return buf[s++];
}
inline void rd(int &x){
  int k;
  int m=0;
  x=0;
  for(;;){
    k = my_getchar();
    if(k=='-'){
      m=1;
      break;
    }
    if('0'<=k&&k<='9'){
      x=k-'0';
      break;
    }
  }
  for(;;){
    k = my_getchar();
    if(k<'0'||k>'9'){
      break;
    }
    x=x*10+k-'0';
  }
  if(m){
    x=-x;
  }
}
inline int rd_int(void){
  int x;
  rd(x);
  return x;
}
struct MY_WRITER{
  char buf[1048576];
  int s;
  int e;
  MY_WRITER(){
    s = 0;
    e = 1048576;
  }
  ~MY_WRITER(){
    if(s){
      fwrite(buf, 1, s, stdout);
    }
  }
}
;
MY_WRITER MY_WRITER_VAR;
void my_putchar(int a){
  if(MY_WRITER_VAR.s == MY_WRITER_VAR.e){
    fwrite(MY_WRITER_VAR.buf, 1, MY_WRITER_VAR.s, stdout);
    MY_WRITER_VAR.s = 0;
  }
  MY_WRITER_VAR.buf[MY_WRITER_VAR.s++] = a;
}
inline void wt_L(char a){
  my_putchar(a);
}
inline void wt_L(int x){
  int s=0;
  int m=0;
  char f[10];
  if(x<0){
    m=1;
    x=-x;
  }
  while(x){
    f[s++]=x%10;
    x/=10;
  }
  if(!s){
    f[s++]=0;
  }
  if(m){
    my_putchar('-');
  }
  while(s--){
    my_putchar(f[s]+'0');
  }
}
const int dx[] = {1, -1, 0, 0};
const int dy[] = { 0, 0, 1, -1};
int n;
int a[110][110];
bool within(int x, int y){
  return 0 <= x  &&  x < n && 0 <= y  &&  y < n;
}
bool check(){
  int i;
  for(i=(0);i<(n);i++){
    int j;
    for(j=(0);j<(n);j++){
      int k;
      for(k=(0);k<(4);k++){
        int ni = i + dx[k];
        int nj = j + dy[k];
        if(!within(ni,nj)){
          continue;
        }
        if(abs(a[ni][nj] - a[i][j]) == 1){
          return false;
        }
      }
    }
  }
  return true;
}
int main(){
  int rbGYC8dh;
  int NN9wFlop = rd_int();
  for(rbGYC8dh=(0);rbGYC8dh<(NN9wFlop);rbGYC8dh++){
    int i;
    rd(n);
    for(i=(0);i<(n+2);i++){
      int j;
      for(j=(0);j<(n+2);j++){
        a[i][j] = 0;
      }
    }
    int left = 2;
    int right = n*n-1;
    for(i=(0);i<(2 * n);i++){
      int j;
      if(i == 0 || i == 2*(n-1)){
        continue;
      }
      for(j=(0);j<(n);j++){
        if(!within(i-j, j) ){
          continue;
        }
        a[i-j][j] = left++;
      }
    }
    a[0][0] = 0;
    if(n >= 2){
      a[n-1][n-1] = 1;
    }
    for(i=(0);i<(n);i++){
      int j;
      for(j=(0);j<(n);j++){
        ++a[i][j];
      }
    }
    if(!check()){
      wt_L(-1);
      wt_L('\n');
      continue;
    }
    for(i=(0);i<(n);i++){
      {
        int a3PBjchb;
        if(n==0){
          wt_L('\n');
        }
        else{
          for(a3PBjchb=(0);a3PBjchb<(n-1);a3PBjchb++){
            wt_L(a[i][a3PBjchb]);
            wt_L(' ');
          }
          wt_L(a[i][a3PBjchb]);
          wt_L('\n');
        }
      }
    }
  }
  return 0;
}
// cLay version 20210405-1

// --- original code ---
// //no-unlocked
// const int dx[] = {1, -1, 0, 0};
// const int dy[] = { 0, 0, 1, -1};
// int n;
// int a[110][110];
// bool within(int x, int y){
//     return 0 <= x < n && 0 <= y < n;
// }
// bool check(){
//     rep(i, n) rep(j, n){
//         rep(k, 4){
//             int ni = i + dx[k];
//             int nj = j + dy[k];
//             if(!within(ni,nj)) continue;
//             if(abs(a[ni][nj] - a[i][j]) == 1) return false;
//         }
//     }
//     return true;
// }
// 
// { 
//     REP(rd_int()){
//         rd(n);
//         rep(i, n+2) rep(j, n+2) a[i][j] = 0;
//         int left = 2, right = n*n-1;
//         rep(i, 2 * n){
//             if(i == 0 || i == 2*(n-1)) continue;
//             rep(j, n){
//                 if(!within(i-j, j) ) continue;
//                 a[i-j][j] = left++;
//             }
//         }
//         a[0][0] = 0;
//         if(n >= 2)
//         a[n-1][n-1] = 1;
//         rep(i, n) rep(j, n) ++a[i][j];
//         if(!check()) wt(-1), continue;
//         rep(i, n) wt(a[i](n));
//     }
// }