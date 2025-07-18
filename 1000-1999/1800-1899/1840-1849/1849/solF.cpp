#include<bits/stdc++.h>
using namespace std;

#define rep(i, a, n) for(int i=(a); i<(n); ++i)
#define per(i, a, n) for(int i=(a); i>(n); --i)
#define pb emplace_back
#define mp make_pair
#define clr(a, b) memset(a, b, sizeof(a))
#define all(x) (x).begin(),(x).end()
#define lowbit(x) (x & -x)
#define fi first
#define se second
#define lson o<<1
#define rson o<<1|1
#define gmid l[o]+r[o]>>1
 
using ll = long long;
using ull = unsigned long long;
using pii = pair<int,int>;
using pll = pair<ll, ll>;
using ui = unsigned int;
 
constexpr int mod = 998244353;
constexpr int inf = 0x3f3f3f3f;
constexpr double EPS = 1e-8;
const double PI = acos(-1.0);

constexpr int N = 2e5 + 10;

int n, a[N];
pii b[N];
char s[N];

void solve(){
  int mask = 0;
  vector<vector<int>> ss;
  per(i, 29, -1){
    int pre = -1;
    int cnt = 0;
    bool flag = 1;
    rep(j, 0, n){
      if((b[j].fi & mask) != pre){
        pre = b[j].fi & mask;
        cnt = 0;
      }
      ++cnt;
      if(cnt > 4){
        flag = 0;
        break;
      }
    }
    if(!flag){
      mask |= (1 << i);
      continue;
    }
    pre = -1;
    vector<int> cur;
    rep(j, 0, n){
      if((b[j].fi & mask) != pre){
        pre = b[j].fi & mask;
        ss.pb(cur);
        cur.clear();
      }
      cur.pb(b[j].se);
    }
    ss.pb(cur);
    break;
  }

  for(auto &p : ss){
    int len = p.size();
    if(len <= 2){
      rep(i, 0, len){
        s[p[i]] = '0' + (i & 1);
      }
    } else {
      int k0 = 0, k1 = 1;
      if(len == 3){
        if((a[p[0]] ^ a[p[2]]) > (a[p[k0]] ^ a[p[k1]])){
          k0 = 0;
          k1 = 2;
        }
        if((a[p[1]] ^ a[p[2]]) > (a[p[k0]] ^ a[p[k1]])){
          k0 = 1;
          k1 = 2;
        }
        s[p[k0]] = s[p[k1]] = '0';
      } else {
        int sum = 0;
        for(int x : p)  sum ^= a[x];
        int val = min(a[p[0]] ^ a[p[1]], sum ^ a[p[0]] ^ a[p[1]]);
        rep(i, 0, 4){
          rep(j, 1, 4){
            int cur = a[p[i]] ^ a[p[j]];
            if(min(cur, sum ^ cur) > val){
              val = min(cur, sum ^ cur);
              k0 = i;
              k1 = j;
            }
          }
        }
        s[p[k0]] = s[p[k1]] = '0';
      }
      rep(i, 0, len){
        if(i != k0 && i != k1){
          s[p[i]] = '1';
        }
      }
    }
  }
}

void _main(){
  cin >> n;
  rep(i, 0, n){
    cin >> a[i];
    b[i].fi = a[i];
    b[i].se = i;
  }
  sort(b, b+n);

  solve();

  s[n] = '\0';
  cout << s << '\n';
}

int main(){
  ios::sync_with_stdio(0);
  cin.tie(0); cout.tie(0);
  _main();
  return 0;
}