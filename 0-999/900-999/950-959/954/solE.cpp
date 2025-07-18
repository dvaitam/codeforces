#include<bits/stdc++.h>
#define fp(i,a,b) for(int i=a,I=b+1;i<I;++i)
#define fd(i,a,b) for(int i=a,I=b-1;i>I;--i)
#define go(u) for(register int i=fi[u],v=e[i].to;i;v=e[i=e[i].nx].to)
#define file(s) freopen(s".in","r",stdin),freopen(s".out","w",stdout)
using namespace std;
char ss[1<<17],*A=ss,*B=ss;
inline char gc(){if(A==B){B=(A=ss)+fread(ss,1,1<<17,stdin);if(A==B)return EOF;}return*A++;}
template<class T>inline void sd(T&x){
    char c;T y=1;while(c=gc(),c<48||57<c)if(c=='-')y=-1;x=c^48;
    while(c=gc(),47<c&&c<58)x=(x<<1)+(x<<3)+(c^48);x*=y;
}
const int N=2e5+5;
const double eps=1e-8;
typedef int arr[N];
int n,T;double tp,ans;
struct da{int a,t;double s;}a[N];
inline bool cmp(da a,da b){return a.t<b.t;}
int main(){
    #ifndef ONLINE_JUDGE
        file("s");
    #endif
    sd(n);sd(T);
    fp(i,1,n)sd(a[i].a);fp(i,1,n)sd(a[i].t),a[i].t-=T;
    sort(a+1,a+n+1,cmp);
    if(a[1].t>0||a[n].t<0)return puts("0"),0;
    fp(i,1,n)a[i].s=(double)a[i].a*a[i].t,tp+=a[i].s;
    if(tp<0){fp(i,1,n)a[i].t=-a[i].t,a[i].s=-a[i].s;reverse(a+1,a+n+1);}
    tp=0;
    fp(i,1,n){
     if(tp+a[i].s>0){ans+=(-tp)/a[i].t;break;}
     tp+=a[i].s,ans+=a[i].a;
 }
 printf("%.6lf",ans);
return 0;
}