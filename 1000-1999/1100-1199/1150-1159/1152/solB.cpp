#pragma gcc optimze("Ofast")
#pragma loop_opt(on)
#include <bits/stdc++.h>
#define debug(x) 1&&(cout<<#x<<':'<<x<<'\n')
#define mkp make_pair
#define pb emplace_back
#define ff first
#define ss second
#define siz(v) (ll(v.size()))
#define all(v) begin(v),end(v)
#define rall(v) rbegin(v),rend(v)
#define REP(i,l,r) for(int i=(l);i<(r);i++)
#define PER(i,l,r) for(int i=(r)-1;i>=(l);i--)
#define mid (l+(r-l>>1))

using namespace std;
typedef long long ll;
typedef long double ld;
typedef pair<ll,ll> pll;
constexpr long double PI = acos(-1),eps = 1e-8;
constexpr ll N = 1<<20, INF = 1e18, MOD = 1000003;

ll n,t,s;
vector<ll> ans;
bool AC(ll x){return __builtin_popcount(x+1)==1;}
signed main(){
    cin >> n;
    while(1){
        if(AC(n)) break;
        t++;
        for(s = 1;s <= N;s<<=1)if(n&s)break;
        s = __builtin_ctz(s);
        //if(s == 0) s++;
        ans.pb(s);
        n ^= (1LL<<s)-1;
        //debug(n);
        if(AC(n)) break;
        t++;
        n++;
    }
    cout << t << '\n';
    for(auto i:ans) cout << i << ' ';
}