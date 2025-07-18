#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef long double ld;
#define rep(x, y, z) for(int x=y; x < (z); x++)
#define all(x) x.begin(), x.end()
#define st first
#define nd second
#define pb push_back
inline int in(){int x; cin>>x; return x;}
const int N=3e5+5, MOD=1e9+7;
void solve(){
  int n=in(), k=in();
  array<ll, 3> a[n];
  map<int, int> M;
  rep(i, 0, n){
    a[i][0]=in(), a[i][1]=in();
    a[i][1]=a[i][0]+a[i][1];
    a[i][2]=i;
    M[a[i][0]]=M[a[i][1]]=M[a[i][1]-1]=1;
  }
  int ind=0;
  for(auto& [X, Y] : M) Y=ind++;
  sort(a, a+n, [&](array<ll, 3> X, array<ll, 3> Y){
    if(X[1] != Y[1]) return X[1] < Y[1];
    return X[0] > Y[0];
  });
  vector<int> ans;
  auto check=[&](int mid, bool save){
    multiset<int> bus;
    rep(i, 0, mid) bus.insert(-1); 
    int cnt=0;
    rep(i, 0, n){
      auto it=bus.upper_bound(a[i][0]);
      if(it != bus.begin()){
        it--;
        bus.erase(it);
        bus.insert(a[i][1]);
        cnt+= 1;
        if(save && ans.size() < k) ans.pb(a[i][2]);
      }
    }
    return (cnt >= k);
  };
  int l=1, r=n;
  while(r-l > 1){
    int mid=(r+l)/2;
    if(check(mid, 0)) r=mid;
    else l=mid;
  }
  if(check(l, 0)) r=l;
  assert(check(r, 1));
  cout<<r<<"\n";
  for(int i : ans) cout<<i+1<<" ";
  cout<<"\n";
}
int main(){
  ios::sync_with_stdio(false);
  cin.tie(0);
  int T=1;
  // T=in();
  while(T--) solve();
  cerr<<double(clock())/CLOCKS_PER_SEC<<"\n";
}