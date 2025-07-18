#include<bits/stdc++.h>
using namespace std;
inline int read(){
 int x = 0,f = 1;char ch = getchar();
 while(ch < '0' || ch > '9'){if(ch == '-')f = -1;ch = getchar();}
 while(ch >= '0' && ch <= '9'){x = x*10+ch-'0';ch = getchar();}
 return x*f;
}

const int N = 100005;
int n;
int a[N];
char b[N];
int l,r;

inline int getMax(int k){
 int Max = -1e9;
 for(int i = 0;i <= 4;++i)
  Max = max(a[k-i],Max);
 return Max;
}

inline int getMin(int k){
 int Min = 1e9;
 for(int i = 0;i <= 4;++i)
  Min = min(a[k-i],Min);
 return Min;
}

int main(){
 n = read();
 for(int i = 1;i <= n;++i)
  a[i] = read();
 scanf("%s",b+1);
 l = -1e9;
 r = 1e9;
 int num = 0;
 for(int i = 5;i <= n;++i){
  if(b[i-1] == '1' && b[i] == '1'){
   l = max(getMin(i),l);
  }
  else if(b[i-1] == '1' && b[i] == '0'){
   r = min(getMin(i)-1,r);
   i += 3;
  }
  else if(b[i-1] == '0' && b[i] == '1'){
   l = max(getMax(i)+1,l);
   i += 3;
  }
  else r = min(getMax(i),r);
 }
 printf("%d %d\n", l,r);
 return 0;
}//13212213123