#include <bits/stdc++.h>
using namespace std;
// #include <ext/pb_ds/assoc_container.hpp>
// #include <ext/pb_ds/tree_policy.hpp>
// using namespace __gnu_pbds;
// template<class T> using pbset=tree<T, null_type, less<T>, rb_tree_tag,tree_order_statistics_node_update> ;
// template<class T> using pbmultiset=tree<T, null_type, less_equal<T>, rb_tree_tag,tree_order_statistics_node_update> ;
using ll = long long;
using ull = unsigned long long;
using lld = long double;
using pii = pair<int, int>;
using pll = pair<ll, ll>;
using vi = vector<int>;
using vl = vector<ll>;
using vii = vector<pii>;
using vll = vector<pll>;
using vvi = vector<vi>;
using vs = vector<string>;
using vb = vector<bool>;
using mii = map<int,int>;
using mll = map<ll,ll>;
using si = set<int>;
using sl = set<ll>;
#define tcT template<class T
#define f(i,x,n) for(int i = x; i < n; i++)
#define rf(i,x,n) for(int i = x; i >= n; i--)
#define sz(a) int((a).size())
#define bg(a) a.begin()
#define re return
#define pb push_back
#define mp make_pair
#define all(a) a.begin(), a.end()
#define rall(a) a.rbegin(), a.rend()
#define sqr(x) (1LL*(x)*(x))
#define gcd(a,b) __gcd(a,b)
#define lcm(a,b) (1LL*(a/gcd(a,b))*b)
#define fix(prec) {cout << setprecision(prec) << fixed;}
#define fi first
#define se second
#define CeilDiv(a,b) ((a+b-1)/b)
#define endl '\n'
#define yes cout<<"YES"<<endl
#define no cout<<"NO"<<endl
#define gg cout<<-1<<endl
#ifndef ONLINE_JUDGE
#define dbg(v) cout << "Line(" << __LINE__ << ") -> " << #v << " = " << (v) << endl;
#include <debugging.h>
#else
#define dbg(v)
#endif
#define lb lower_bound
#define ub upper_bound
tcT> using V = vector<T>;
tcT> int lwb(V<T>& a, const T& b) { return int(lb(all(a),b)-bg(a)); }
tcT> int upb(V<T>& a, const T& b) { return int(ub(all(a),b)-bg(a)); }
tcT> using pqmax = priority_queue<T>;
tcT> using pqmin = priority_queue<T,vector<T>,greater<T>>;
tcT> bool ckmin(T& a, const T& b) {
    return b < a ? a = b, 1 : 0; } // set a = min(a,b)
tcT> bool ckmax(T& a, const T& b) {
    return a < b ? a = b, 1 : 0; } // returns true if value changed
tcT> istream& operator>>(istream& is,  vector<T> &v){for (auto& i : v) is >> i; return is;}
tcT> ostream& operator<<(ostream& os,  vector<T> &v){for (auto& i : v) os << i << ' '; return os;}
#define tr(c, i) for (typeof (c).begin() i = c.begin(); i != c.end(); i++)
#define present(c, x) (c.find(x) != c.end())
#define cpresent(c, x) (find(all(c), x) != c.end())

///.........Bit_Manipulation...........///
#define MSB(mask) 63-__builtin_clzll(mask)  /// 0 -> -1
#define LSB(mask) __builtin_ctzll(mask)  /// 0 -> 64
#define SETBIT(mask) __builtin_popcountll(mask)
#define CHECKBIT(mask,bit) (mask&(1LL<<bit))
#define ONBIT(mask,bit) (mask|(1LL<<bit))
#define OFFBIT(mask,bit) (mask&~(1LL<<bit))
#define CHANGEBIT(mask,bit) (mask^(1LL<<bit))
const int inf = 2e9;
const ll mod = 1000000007;
// const ll mod = 998244353;

void solve(){
    int n;
    cin >> n;
    vi a(n-1);
    cin >> a;
    vi b(n);
    b[0] = 501;
    f(i,1,n){
        b[i] = b[i-1] + a[i-1];
    }
    cout << b << endl;
}

signed main(){
    ios_base::sync_with_stdio(0);
    cin.tie(0); cout.tie(0);
    int t = 1;
    cin >> t;
    for(int i = 1; i <= t; i++){
      //  cout << "Case #" << i << ": ";
        solve();
    }
    return 0;
}