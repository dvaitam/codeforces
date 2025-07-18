#include <bits/stdc++.h>
#define fi first
#define se second
#define rep(i,n) for(int i = 0; i < (n); ++i)
#define rrep(i,n) for(int i = 1; i <= (n); ++i)
#define drep(i,n) for(int i = (n)-1; i >= 0; --i)
#define srep(i,s,t) for (int i = s; i < t; ++i)
#define rng(a) a.begin(),a.end()
#define maxs(x,y) (x = max(x,y))
#define mins(x,y) (x = min(x,y))
#define limit(x,l,r) max(l,min(x,r))
#define lims(x,l,r) (x = max(l,min(x,r)))
#define isin(x,l,r) ((l) <= (x) && (x) < (r))
#define pb push_back
#define sz(x) (int)(x).size()
#define pcnt __builtin_popcountll
#define uni(x) x.erase(unique(rng(x)),x.end())
#define snuke srand((unsigned)clock()+(unsigned)time(NULL));
#define show(x) cout<<#x<<" = "<<x<<endl;
#define PQ(T) priority_queue<T,v(T),greater<T> >
#define bn(x) ((1<<x)-1)
#define dup(x,y) (((x)+(y)-1)/(y))
#define newline puts("")
#define v(T) vector<T>
#define vv(T) v(v(T))
using namespace std;
typedef long long int ll;
typedef unsigned uint;
typedef unsigned long long ull;
typedef pair<int,int> P;
typedef vector<int> vi;
typedef vector<vi> vvi;
typedef vector<ll> vl;
typedef vector<P> vp;
inline int in() { int x; scanf("%d",&x); return x;}
template<typename T>inline istream& operator>>(istream&i,v(T)&v)
{rep(j,sz(v))i>>v[j];return i;}
template<typename T>string join(const v(T)&v)
{stringstream s;rep(i,sz(v))s<<' '<<v[i];return s.str().substr(1);}
template<typename T>inline ostream& operator<<(ostream&o,const v(T)&v)
{if(sz(v))o<<join(v);return o;}
template<typename T1,typename T2>inline istream& operator>>(istream&i,pair<T1,T2>&v)
{return i>>v.fi>>v.se;}
template<typename T1,typename T2>inline ostream& operator<<(ostream&o,const pair<T1,T2>&v)
{return o<<v.fi<<","<<v.se;}
template<typename T>inline ll suma(const v(T)& a) { ll res(0); for (auto&& x : a) res += x; return res;}
const double eps = 1e-10;
const ll LINF = 1001002003004005006ll;
const int INF = 1001001001;
#define dame { puts("-1"); return 0;}
#define yn {puts("YES");}else{puts("NO");}
const int MX = 200005;

// fast set
int bsr(ull x) { return 63^__builtin_clzll(x);}
int bsf(ull x) { return __builtin_ctzll(x);}
struct faset {
  int n, h;
  vv(ull) seg;
  faset() {}
  faset(int n):n(n) {
    while (n > 1) {
      n = (n+63)>>6;
      seg.pb(v(ull)(n));
    }
    h = sz(seg);
  }
  bool count(int x) const {
    int d = x>>6, r = x&63;
    return seg[0][d]>>r&1;
  }
  void insert(int x) {
    rep(i,h) {
      ull b = 1ull<<(x&63);
      x >>= 6;
      seg[i][x] |= b;
    }
  }
  void erase(int x) {
    rep(i,h) {
      ull b = 1ull<<(x&63);
      x >>= 6;
      seg[i][x] &= ~b;
      if (seg[i][x]) break;
    }
  }
  // x <= res
  int nxt(int x) {
    rep(i,h) {
      int d = x>>6, r = x&63;
      if (d == sz(seg[i])) break;
      ull s = seg[i][d]>>r;
      if (s) {
        x += bsf(s);
        while (i--) {
          x = x<<6|bsf(seg[i][x]);
        }
        return x;
      }
      x = (x>>6)+1;
    }
    return -1;
  }
  // res <= x
  int pre(int x) {
    rep(i,h) {
      if (x == -1) break;
      int d = x>>6, r = x&63;
      ull s = seg[i][d]<<(63^r);
      if (s) {
        x -= 63^bsr(s);
        while (i--) {
          x = x<<6|bsr(seg[i][x]);
        }
        return x;
      }
      x = (x>>6)-1;
    }
    return -1;
  }    
};
//
// Binary Indexed Tree
struct bit {
  int n; vector<ll> d;
  bit() {}
  bit(int mx) {
    n = 1;
    while (n < mx) n <<= 1;
    d = vl(n+1);
  }
  void add(int i, ll x) {
    for (++i;i<=n;i+=i&-i) d[i] += x;
  }
  ll sum(int i) {
    ll x = 0;
    for (++i;i;i-=i&-i) x += d[i];
    return x;
  }
  ll r;
  int nxt(ll x) {
    // cerr<<x<<": "<<d<<endl;
    r = x;
    int i = 0;
    int w = n>>1;
    while (w) {
      if (d[i|w] <= x) {
        i |= w;
        x -= d[i];
      }
      w >>= 1;
    }
    r -= x;
    return i;
  }
};
//

// coordinate compression
struct X {
  typedef int T;
  vector<T> d;
  X() {}
  // X(vector<T>& a): d(a) {
  //   init();
  //   for (T& na : a) na = (*this)(na);
  // }
  void add(T x) { d.pb(x);}
  void init() {
    sort(rng(d));
    d.erase(unique(rng(d)), d.end());
  }
  int size() const { return sz(d);}
  T operator[](int i) const { return d[i];}
  int operator()(T x) const { return upper_bound(rng(d),x)-d.begin()-1;}
};
//



int main() {
  int qn;
  scanf("%d",&qn);
  vi q;
  X xs;
  rep(qi,qn) {
    char c; int a;
    scanf(" %c%d",&c,&a);
    xs.add(a);
    if (c == '-') a = -a;
    q.pb(a);
  }
  xs.init();
  int n = sz(xs);
  bit d(n+2);

  vi cs(n+1);
  faset s(n+1);
  faset bs(n+1);
  vl w(n+1);

  auto check = [&](int i) {
    int j = s.pre(i-1);
    if (j == -1 || xs[j]*2 >= xs[i]) {
      bs.erase(i);
    } else {
      bs.insert(i);
      w[i] = d.sum(i-1);
    }
  };

  int cnt = 0;
  rep(qi,qn) {
    int a = q[qi];
    int i = xs(abs(a));
    if (a < 0) --cnt; else ++cnt;
    if (!cs[i]) {
      s.insert(i);
      check(i);
    }
    if (a < 0) --cs[i]; else ++cs[i];
    if (!cs[i]) {
      if (bs.count(i)) bs.erase(i);
      s.erase(i);
    }
    {
      int j = s.nxt(i+1);
      if (j != -1) check(j);
    }

    d.add(i,a);
    {
      int j = i;
      while (1) {
        j = bs.nxt(j+1);
        if (j == -1) break;
        w[j] += a;
      }
    }

    int ans = cnt-1;
    {
      int j = -1;
      while (1) {
        j = bs.nxt(j+1);
        if (j == -1) break;
        
        // cerr<<j<<" "<<xs[j]<<" "<<w[j]<<endl;
        if (w[j]*2 < xs[j]) ans--;
      }
    }

    // cerr<<"ans:"<<" "<<ans<<endl;
    printf("%d\n", max(ans,0));
  }
  return 0;
}