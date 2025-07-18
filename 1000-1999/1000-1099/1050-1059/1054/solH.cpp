// becoder Submission #undefined @ 1714646675162
#include<bits/stdc++.h>
#define MAXN (2000005)
#define ll long long
#define inv(x) qpow(x,mod-2)
#define PI (acos(-1.0))
using namespace std;
const int mod=490019,B=32768;
const double eps=1e-6,hf=0.5;
void File()
{
    freopen("in.txt","r",stdin);
    freopen("out.txt","w",stdout);
}
template<typename type>
void read(type &x)
{
    bool f=0;char ch=0;x=0;
    while(ch<'0'||ch>'9'){f|=!(ch^'-');ch=getchar();}
    while(ch>='0'&&ch<='9'){x=(x<<3)+(x<<1)+(ch^48);ch=getchar();}
    x=f?-x:x;
}
template<typename type,typename... Args>
void read(type &x,Args &... args)
{
    read(x);
    read(args...);
}
int qpow(int x,int y)
{
    int res=1;
    while(y)
    {
        if(y&1) res=1ll*res*x%mod;
        x=1ll*x*x%mod,y>>=1;
    }
    return res;
}
struct comp{
    double Re,Im;
    comp(){}comp(double nx,double ny):Re(nx),Im(ny){}
    comp operator+(comp rsh){return (comp){Re+rsh.Re,Im+rsh.Im};}
    comp operator-(comp rsh){return (comp){Re-rsh.Re,Im-rsh.Im};}
    comp operator*(comp rsh){return (comp){Re*rsh.Re-Im*rsh.Im,Re*rsh.Im+Im*rsh.Re};}
    comp operator*(double rsh){return (comp){Re*rsh,Im*rsh};}
    comp operator/(double rsh){return (comp){Re/rsh,Im/rsh};}
}f[2][MAXN],g[MAXN],w[MAXN];
int n,m,c,ans;
int a[MAXN],b[MAXN],rev[MAXN];
void getrev(int lim){for(int i=0;i<lim;i++) rev[i]=((rev[i>>1]>>1)|((i&1)?(lim>>1):0));}
int getlim(int x){return 1<<(32-__builtin_clz(x-1));}
void DFT(comp *arr,int typ,int lim)
{
    w[0]={1.0,0.0};
    for(int i=0;i<lim;i++) if(i<rev[i]) swap(arr[i],arr[rev[i]]);
    for(int len=1;len<lim;len<<=1)
    {
        comp t={cos(PI/len),sin(PI/len)*typ};
        for(int i=len-2;i>=0;i-=2) w[i+1]=(w[i]=w[i>>1])*t;
        for(int k=0;k<lim;k+=len<<1)
        {
            for(int p=0;p<len;p++)
            {
                comp cv=w[p]*arr[k|len|p];
                arr[k|len|p]=arr[k|p]-cv;
                arr[k|p]=arr[k|p]+cv;
            }
        }
    }
}
void tFFT(int *F,int *G,int len)
{
    int lim=getlim(len);
    getrev(lim);
    for(int i=0;i<lim;i++)
    {
        int x=F[i]>>15,y=F[i]&(B-1),z=G[i]>>15,w=G[i]&(B-1);
        f[0][i]={(double)x,(double)y};
        f[1][i]={(double)x,-(double)y};
        g[i]={(double)z,(double)w};
    }
    DFT(f[0],1,lim),DFT(f[1],1,lim),DFT(g,1,lim);
    for(int i=0;i<lim;i++)
    {
        g[i].Re/=lim,g[i].Im/=lim;
        f[0][i]=f[0][i]*g[i],f[1][i]=f[1][i]*g[i];
    }
    DFT(f[0],-1,lim),DFT(f[1],-1,lim);
    for(int i=0;i<len;i++)
    {
        ll x=(ll)floor((f[0][i].Re+f[1][i].Re)*hf+hf)%mod;
        ll y=(ll)floor((f[0][i].Im+f[1][i].Im)*hf+hf)%mod;
        ll z=((ll)floor(f[0][i].Im+hf)-y)%mod;
        ll w=((ll)floor(f[1][i].Re+hf)-x)%mod;
        F[i]=((x*B%mod*B%mod+(y+z)*B%mod+w)%mod+mod)%mod;
    }
}
int main()
{
    // File();
    // fdsf
    read(n,m,c);
    for(int i=0,x,p;i<n;i++)
    {
        read(x);
        p=1ll*i*i%(mod-1);
        a[p]=(a[p]+1ll*x*inv(qpow(c,((1ll*p*(p-1))>>1)%(mod-1)))%mod)%mod;
    }
    for(int i=0,x,p;i<m;i++)
    {
        read(x);
        p=1ll*i*i*i%(mod-1);
        b[p]=(b[p]+1ll*x*inv(qpow(c,((1ll*p*(p-1))>>1)%(mod-1)))%mod)%mod;
    }
    tFFT(a,b,mod<<1);
    for(int i=0;i<mod<<1;i++) ans=(ans+1ll*a[i]*qpow(c,((1ll*i*(i-1))>>1)%(mod-1))%mod)%mod;
    printf("%d",ans);
}