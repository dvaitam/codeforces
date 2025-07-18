// Duet of Dusk Embers--XrkArul

#include <bits/stdc++.h>

#include <ext/rope>

using namespace std;

using namespace __gnu_cxx;

#define ll long long

#define lll __int128_t

#define ull unsigned long long

#define vi vector<int>

#define vll vector<ll>

#define endl '\n'

#define ednl '\n'

#define pb push_back

#define fi first

#define se second

#define ls (p<<1)

#define rs (p<<1|1)

#define fix setprecision

#define all(v) (v).begin(),(v).end()

#define pii pair<int, int>

#define pll pair<ll, ll>

#define rep(i, a, b) for (int i = a; i <= b; ++i)

#define pq priority_queue<int, vector<int>>

#define PQ priority_queue<int, vector<int>, greater<int>>

const int inf = 2147483647;

const ll mod=1e9+7;

ll powmod(ll a,ll b){ll s=1;a%=mod;while(b){if(b&1)s=s*a%mod;b>>=1;a=a*a%mod;}return s%mod;}



//#define int long long

void solve(){

    vector<pii> p(3+1);

    rep(i,1,3)cin>>p[i].fi>>p[i].se;

    auto cal=[&](double a,double b,double x,double y){

        return sqrt(abs(a-x)*abs(a-x)+abs(b-y)*abs(b-y));

    };

    double ans=0;

    if(p[1].se==p[2].se&&p[1].se>p[3].se){

        ans=cal(p[1].fi,p[1].se,p[2].fi,p[2].se);

    }else if(p[2].se==p[3].se&&p[2].se>p[1].se){

        ans=cal(p[3].fi,p[3].se,p[2].fi,p[2].se);

    }else if(p[1].se==p[3].se&&p[3].se>p[2].se){

        ans=cal(p[1].fi,p[1].se,p[3].fi,p[3].se);

    }

    cout<<fixed<<setprecision(10)<<ans<<endl;

}

signed main(){

    ios::sync_with_stdio(false);

    cin.tie(nullptr);

    int T=1;

    cin>>T;

    while(T--){

        solve();

    }

    return 0;

}

/*



*/