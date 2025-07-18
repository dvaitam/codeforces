#include <bits/stdc++.h>
#define debug(x) 1&&(cout << #x << ':' << x << '\n')
#define mkp make_pair
#define ff first
#define ss second
#define REP(i,l,r) for(int i = l;i < r;i++)
#define mid (l+(r-l>>1))
#define N 200000
#define EL cout<<'\n'

using namespace std;
typedef long long ll;
typedef pair<ll,ll> pll;
struct piii{
    int ff,ss,th;
    piii(int f=0,int s=0,int t=0):ff(f),ss(s),th(t){}
    bool operator<(piii b){return ff!=b.ff?ff<b.ff:ss!=b.ss?ss<b.ss:th<b.th;}
};
const long double PI = acos(-1);
const ll INF = 1e18, MOD = 1000000007;

//#define __fast_IO
#ifdef __fast_IO
#define EL out('\n')
template<typename T> bool input(T &x){
    x = 0;
    char c = getchar(),f = 0;
    while(c!=EOF&&c!='-'&&!isdigit(c))c = getchar();
    if(c == '-')f = 1,c = getchar();
    while(c!=EOF&&isdigit(c))x=x*10+c-'0',c = getchar();
    if(f)x = -x;
    return c!=EOF;
}
const int __buf_max = 10000; int __buf_p = 0; char __buf[__buf_max+40] = {};
void out(){__buf[__buf_p++] = '\0', puts(__buf), __buf_p = 0;}
void out(char c){__buf[__buf_p++] = c;if(__buf_p >= __buf_max)out();}
template<typename T> void out(T x){
    if(x<0)__buf[__buf_p++] = '-', x = -x;
    if(x==0)return out('0'),void();
    char tmp[40] = {},p = 0;
    while(x)tmp[p++]=x%10+'0',x/=10;
    while(p--)__buf[__buf_p++] = tmp[p];
    if(__buf_p >= __buf_max)out();
}
#else
#endif // __fast_IO
//#define __number_theorem
#ifdef __number_theorem
ll invmod(ll a,ll m = MOD){
    ll mod = m, x = 1, y = 0, t;
    while(m)t=x-(a/m)*y,x=y,y=t, t=a%m,a=m,m=t;
    x%=mod, x<0&&(x+=mod);
    return x;
}
ll modadd(ll&a,ll b,ll m = MOD){ return a += b, (a>=m)&&a%=m, (a<0)&&a+=m, a;}
ll modpow(ll e,ll p,ll m = MOD){ ll r = 1; while(p) (p&1)&&(r=r*e%m), e=e*e%m, p>>=1; return r;}
ll fracMod[N+1] = {1};
ll mu[N+1] = {},phi[N+1] = {};
vector<ll> primes;
bitset<N+1> notP;
void init(){//O(N) builds phi, mu, primes, fracMod
    for(int i = 1;i <= N;i++)fracMod[i] = fracMod[i-1]*i;
    phi[1] = mu[1] = 1;
    for(ll i = 2;i <= N;i++){
        if(!notP[i]){
            primes.push_back(i);
            phi[i] = i-1, mu[i] = -1;
        }
        for(auto p:primes){
            if(i*p>N)break;
            notP[i*p] = true;
            if(i%p){
                phi[i*p] = phi[i]*(p-1), mu[i*p] = -mu[i];
            }else{
                phi[i*p] = phi[i]*p, mu[i*p] = 0;
                break;
            }
        }
    }
}
bool is_p(int x){}
ll C(ll n,ll m,ll p){
    ll res = 1;
    while(m){
        ll a = n%p, b = m%p;
        if(b<0||a-b<0)return 0;
        res = res*fracMod[a]%p*invmod(fracMod[a-b],p)%p*invmod(fracMod[b],p)%p;
        n/=p, m/=p;
    }
    return res;
}
#endif // __number_theorem


signed main(){
    ios::sync_with_stdio(0), cin.tie(0);
    //init();
    ll n, m, a[101] = {},tot=0,res=0;
    cin >> n >> m;
    REP(i,0,n)cin >> a[i],tot+=a[i];
    if(tot<m)return cout<<-1<<'\n',0;
    sort(a,a+n);
    for(int day = 1;day <= m;day++){
        int tot = 0;
        for(int i = n-1,t=0,cnt=0;i >= 0;i--){
            if(a[i]<cnt)break;
            tot += a[i]-cnt;
            if(++t==day)cnt++,t=0;
            //debug(a[i]),debug(cnt),debug(t);
        }
        if(tot>=m)return cout<<day<<'\n',0;
    }
    cout<<-1<<'\n';
}