#include<bits/stdc++.h>
using namespace std;
template<class T1,class T2> bool cmax(T1& a,const T2& b){return a<b?a=b,1:0;}
template<class T1,class T2> bool cmin(T1& a,const T2& b){return b<a?a=b,1:0;}
#define fir(i,x,y,...) for(int i(x),##__VA_ARGS__;i<=(y);i++)
#define firr(i,x,y,...) for(int i(x),##__VA_ARGS__;i>=(y);i--)
#define fird(i,x,y,d,...) for(int i(x),##__VA_ARGS__;i<=(y);i+=(d))
#define firrd(i,x,y,d,...) for(int i(x),##__VA_ARGS__;i>=(y);i-=(d))
#define Yes cout<<"Yes\n"
#define No cout<<"No\n"
#define YES cout<<"YES\n"
#define NO cout<<"NO\n"
#define endl "\n"
#define ll long long
#define ull unsigned long long
int T;
bool _mul=1;
const int N=6e5+5;
int n;
ll k;
int a[N],b[N],c[N],p[N];
ll gcd(ll x,ll y){if(!y) return x;return gcd(y,x%y);}
ll power(ll b,ll p,ll mod){ll res=1;while(p){if(p&1) res=res*b%mod;b=b*b%mod;p>>=1;}return res;}
int t[N];
void add(int x){for(int i=x;i<=n;i+=(i&-i)) t[i]++;}
int ask(int x){int res=0;for(int i=x;i;i-=(i&-i)) res+=t[i];return res;}
void solve()
{
    cin>>n>>k;
    // ll K=k;
    fir(i,1,n,x) cin>>x,p[i]=x,a[x]=i,t[i]=b[i]=0;
    ll sum=1ll*n*(n-1)/2;
    fir(i,1,n)
    {
        sum-=ask(a[i]-1);
        add(a[i]);
    }
    // cout<<sum<<endl;
    k-=sum;
    if(k%2!=0||k/2<0||k/2>(1ll*n*(n-1)/2-sum)){NO;return;}
    k/=2;
    // cout<<k<<endl;
    fir(i,1,n) t[i]=0;
    fir(i,1,n)
    {
        int x=ask(a[i]-1);
        if(k>x) k-=x;
        else if(k&&k<=x)
        {
            int id=0;
            fir(j,1,i)
            {
                if(a[j]<a[i]&&k)
                {
                    id=j;
                    k--;
                }
            }
            fir(j,1,id) b[j]=i+1-j;
            fir(j,id+1,i-1) b[j]=i-j;
            b[i]=i-id;
        }
        else b[i]=i;
        add(a[i]);
    }
    YES;
    fir(i,1,n) cout<<b[i]<<" ";
    cout<<endl;
    // assert(cnt==K);
    return;
}
int main()
{
    //freopen(".in","r",stdin);
    //freopen(".out","w",stdout);
    cin.tie(nullptr)->sync_with_stdio(false);
    if(_mul) cin>>T;
    else T=1;
    while(T--) solve();
    return 0;
}
/*
7
3 0
2 1 3
3 1
2 1 3
3 2
2 1 3
3 3
2 1 3
3 4
2 1 3
3 5
2 1 3
3 6
2 1 3
*/