#include <iostream>
#include <algorithm>
#include <sstream>
#include <string>
#include <queue>
#include <cstdio>
#include <map>
#include <set>
#include <utility>
#include <stack>
#include <cstring>
#include <cmath>
#include <vector>
#include <ctime>
#include <bitset>
#include <assert.h>
using namespace std;
#define pb push_back
#define sd(n) scanf("%d",&n)
#define sdd(n,m) scanf("%d%d",&n,&m)
#define sddd(n,m,k) scanf("%d%d%d",&n,&m,&k)
#define sld(n) scanf("%lld",&n)
#define sldd(n,m) scanf("%lld%lld",&n,&m)
#define slddd(n,m,k) scanf("%lld%lld%lld",&n,&m,&k)
#define sf(n) scanf("%lf",&n)
#define sff(n,m) scanf("%lf%lf",&n,&m)
#define sfff(n,m,k) scanf("%lf%lf%lf",&n,&m,&k)
#define ss(str) scanf("%s",str)
#define ansn() printf("%d\n",ans)
#define lansn() printf("%lld\n",ans)
#define r0(i,n) for(int i=0;i<(n);++i)
#define r1(i,e) for(int i=1;i<=e;++i)
#define rn(i,e) for(int i=e;i>=1;--i)
#define mst(abc,bca) memset(abc,bca,sizeof abc)
#define lowbit(a) (a&(-a))
#define all(a) a.begin(),a.end()
#define pii pair<int,int>
#define pll pair<long long,long long>
#define mp(aa,bb) make_pair(aa,bb)
#define lrt rt<<1
#define rrt rt<<1|1
#define X first
#define Y second
#define PI (acos(-1.0))
double pi = acos(-1.0);
typedef long long ll;
typedef unsigned long long ull;
typedef long double ld;
const ll mod = 1000000007;
const double eps=1e-12;
const int inf=0x3f3f3f3f;
//const ll infl = 100000000000000000;//1e17
const int maxn=  2e6+20;
const int maxm = 5e3+20;
//muv[i]=(p-(p/i))*muv[p%i]%p;
inline int in(int &ret) {
    char c;
    int sgn ;
    if(c=getchar(),c==EOF)return -1;
    while(c!='-'&&(c<'0'||c>'9'))c=getchar();
    sgn = (c=='-')?-1:1;
    ret = (c=='-')?0:(c-'0');
    while(c=getchar(),c>='0'&&c<='9')ret = ret*10+(c-'0');
    ret *=sgn;
    return 1;
}
ll dp[33][33];
int a[33];
ll dfs(int cur,bool limit,int had)
{
    if(cur>19)return 1;
    if(!limit&&dp[cur][had]!=-1)return dp[cur][had];
    if(had==3)return dp[cur][3] = 1;
    ll res = 0;
    int mx = 9;
    if(limit)mx = a[cur];
    for(int i=0;i<=mx;++i)
    {
        if((a[cur]==i)&&limit)
            res += dfs(cur+1,1,had + ( i!=0));
        else
            res += dfs(cur+1,0,had + ( i!=0));
    }
    if(!limit)return dp[cur][had] = res;
    return res;
}
ll solve(ll x)
{
    if(x<0)return 0;
//    cout<<x<<' ';
    r1(i,19)
    {
        a[19-i+1] = x%10;
        x/=10;
    }
    ll ret = dfs(1,1,0);
//    cout<<ret<<endl;
    return ret;
}
int main() {
#ifdef LOCAL
    freopen("input.txt","r",stdin);
//    freopen("output.txt","w",stdout);
#endif // LOCAL
    mst(dp,-1);
    int t;
    sd(t);
    for(; t--;) {
        ll l,r;
        sldd(l,r);
        ll ans = solve(r) - solve(l-1);
        lansn();
    }

    return 0;
}