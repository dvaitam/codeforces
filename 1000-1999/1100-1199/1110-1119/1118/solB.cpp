#include<cstdio>
#include<iostream>
#include<algorithm>
#include<cstring>
#include<queue>
#include<map>
#include<set>
#include<cmath>
using namespace std;
#define in() freopen(".in","r",stdin)
#define out() freopen(".out","w",stdout)
#define I inline
#define R register
#define ls(x) (x<<1)
#define rs(x) (x<<1|1)
#define lb(x) (x&-x)
#define EL putchar(10)
#define SP putchar(32)
#define m(a) memset(a,0,sizeof(a))
#define gc() getchar()
typedef long long ll;
const int inf=2147483647;
const int ine=-2147483647;
//char ss[1<<17],*A=ss,*B=ss;I char gc(){return A==B&&(B=(A=ss)+fread(ss,1,1<<17,stdin),A==B)?EOF:*A++;}
template <typename _Tp> I void _(_Tp &x);
template <typename _Tq> I void wr(_Tq x);
I int mx(int a,int b){return a>b?a:b;}I int mn(int a,int b){return a<b?a:b;}
int n,a[200010],s[2][200010],sum[200010],ans;
int main(){
    //in();out();
    _(n);
    for(R int i=1;i<=n;++i) _(a[i]),sum[i]=sum[i-1]+a[i];
    for(R int i=2;i<=n;++i)
        s[0][i]=sum[i-1]-s[0][i-1];
    for(R int i=n-1;i>0;--i)
        s[1][i]=sum[n]-sum[i]-s[1][i+1];
    for(R int i=1;i<=n;++i){
        if(s[0][i]+s[1][i+1]==s[1][i]+s[0][i-1])
            ++ans;
    }
    wr(ans);
    return 0;
}
template <typename _Tp>
    I void _(_Tp &x){
        _Tp w=1;char c=0;x=0;
        while (c^45&&(c<48||c>57)) c=gc();
        if (c==45) w=-1, c=gc();
        while(c>=48&&c<=57) x=(x<<1)+(x<<3)+(c^48),c=gc();
        x*=w;
    }
template <typename _Tq>
    I void wr(_Tq x){
        if(x<0)
            putchar(45),x=-x;
        if(x<10){
            putchar(x+48);
            return;
        }
        wr(x/10);
        putchar(x%10+48);
    }