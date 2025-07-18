#include <bits/stdc++.h>
using namespace std;
//solution starts at line 66

#define us unsigned
using ll = long long;
using db = double;
using str = string;

#define pr pair
using pi = pr<int,int>;
using pl = pr<ll,ll>;
using pd = pr<db,db>;
#define mp make_pair
#define fi first
#define se second

#define vc vector
using vi = vc<int>;
using vl = vc<ll>;
using vd = vc<db>;
using vs = vc<str>;
using vb = vc<bool>;
#define vpr vc<pr>
using vpi = vc<pi>;
using vpd = vc<pd>;
using vpl = vc<pl>;
#define pb push_back
#define all(x) x.begin(),x.end()
#define rall(x) x.rbegin(),x.rend()
#define sz(x) x.size()

template <typename T>
using mt = vc<vc<T> >;
using mti = mt<int>;
using mtl = mt<ll>;
using mtd = mt<db>;
using mtb = mt<bool>;

#define mxh priority_queue
template <typename T>
using mnh = mxh<T,vc<T>,greater<T> >;
#define ps push

constexpr int pct(int x){return __builtin_popcount(x);}
constexpr int clz(int x){return __builtin_clz(x);}
constexpr int ctz(int x){return __builtin_ctz(x);}
constexpr int pct(ll x){return __builtin_popcount(x);}
constexpr int clz(ll x){return __builtin_clz(x);}
constexpr int ctz(ll x){return __builtin_ctz(x);}

#define F_OR(i,f,t,a) for(int i = (f); i<(t); i += (a))
#define F_OR1(t) F_OR(i,0,t,1)
#define F_OR2(i,t) F_OR(i,0,t,1)
#define F_OR3(i,f,t) F_OR(i,f,t,1)
#define F_OR4(i,f,t,a) F_OR(i,f,t,a)
#define GET_MACRO(a,b,c,d,e,...) e
#define FOR_C(...) GET_MACRO(__VA_ARGS__, F_OR4, F_OR3, F_OR2, F_OR1)
#define FOR(...) FOR_C(__VA_ARGS__)(__VA_ARGS__)

#define read(x, n) FOR(obscurevariable,n) cin >> x[obscurevariable];
#define write(x,n,s) FOR(obscurevariable,n) cout << x[obscurevariable] << s;
#define yn(x) ((x)?"YES":"NO")

int dx[4] = {1,0,-1,0};
int dy[4] = {0,1,0,-1};

struct dtree{
  vector<vi> tr;
  int s = 0;
  int rt = 0;
  void resize(int n){
    s = n;
    tr.resize(s);
  }
  void add(int from, int to){
    tr[from].push_back(to);
  }
  void init(vi in){
    int n = sz(in);
    resize(n);
    FOR(n){
      if(in[i]-1==i)
        rt = i;
      else
        add(in[i]-1,i);
    }
  }
};

mti ans;

void dfs(dtree& t, int id, int pth){
  ans[pth].pb(id);
  int n = t.tr[id].size();
  if(n==0)
    return;
  dfs(t,t.tr[id][0],pth);
  FOR(i,1,n){
    int s = ans.size();
    ans.pb(vi());
    dfs(t,t.tr[id][i],s);
  }
}

void solve(){
  int n; cin >> n;
  vi p(n);
  read(p,n);
  dtree x;
  x.init(p);
  ans = mti(1);
  dfs(x,x.rt,0);
  int m = ans.size();
  cout << m << '\n';
  FOR(m){
    int k = ans[i].size();
    cout << k << '\n';
    FOR(j,k)
      cout << ans[i][j]+1 << ' ';
    cout << '\n';
  }
  cout << '\n';
}

int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);
    int tt = 1;
    cin >> tt;
    FOR(tt){
      solve();
    }
    return 0;
}