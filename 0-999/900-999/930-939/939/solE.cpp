/*"Nothing is true."*/
#include<bits/stdc++.h>
using namespace std;
inline int min(int _a, int _b) {return _a < _b ? _a : _b;}
inline int max(int _a, int _b) {return _a > _b ? _a : _b;}
const int inf = 0x3f3f3f3f;
const double eps = 1e-8;

template <class _T> inline void rd(_T &_a) {
 int _f=0,_ch=getchar();_a=0;
 while(_ch<'0' || _ch>'9'){if(_ch=='-')_f=1;_ch=getchar();}
 while(_ch>='0' && _ch<='9')_a=(_a<<1)+(_a<<3)+_ch-'0',_ch=getchar();
 if(_f)_a=-_a;
}
long long cnt, idx, sum; double aver = 10000000000;
long long a[1010000], Max;
int main() {
 //freopen("CF939E.in","r",stdin);
 //freopen("CF939E.out","w",stdout);
 int Q, op; rd(Q);
 while(Q--) {
  rd(op);
  if(op == 1) {
   rd(a[++idx]); Max = a[idx];
  }
  else {
   aver = (sum + a[idx]) * 1.0 / (cnt + 1);
   while((sum + a[idx] + a[cnt + 1]) * 1.0 / (cnt + 2) < aver && cnt < idx - 1) {
    aver = (sum + a[idx] + a[cnt + 1]) * 1.0 / (cnt + 2);
    sum += a[++cnt];
   }
   printf("%.8lf\n", Max - aver);
  }
 }
 return 0;
}