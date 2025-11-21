#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
int main(){
  ios::sync_with_stdio(false);
  cin.tie(nullptr);
  int t; if(!(cin>>t)) return 0;
  while(t--){
    int n; ll k; cin>>n>>k;
    vector<ll>a(n),b(n);
    for(int i=0;i<n;i++) cin>>a[i];
    for(int i=0;i<n;i++) cin>>b[i];
    ll best=min(2LL*k,n);
    ll sum=0,ans=0;
    for(int i=0;i<best;i++){
      sum+=a[i];
      ans=max(ans,sum+(best-i-1)*b[i]);
    }
    cout<<ans<<"\n";
  }
  return 0;
}
