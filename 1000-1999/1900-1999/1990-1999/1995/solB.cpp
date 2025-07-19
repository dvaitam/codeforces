#include <bits/stdc++.h>
#include <limits.h>
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
using namespace std;
using namespace __gnu_pbds;
#define ll long long

template<typename T>
using ordered_set = tree<T, null_type, less_equal<T>, rb_tree_tag, tree_order_statistics_node_update>;

void fastIO() {
    ios::sync_with_stdio(false);
    cin.tie(NULL);
    cout.tie(NULL);
}

int main() {
    fastIO();
ll t; cin>>t;
while(t--){
ll n,m; cin>>n>>m;
vector<ll>v(n);
for(ll i=0 ; i<n ; i++){
  cin>>v[i];
}
sort(v.begin(),v.end());
ll l=0, r=1,sum=v[0],ans=0;
if(v.front()>m){
  cout<<0<<'\n';
  continue;
}
while(l<n){
if(r<n && v[r]-v[l]<=1 && sum+v[r]<=m){
  sum+=v[r];
  r++;
}
else{
  sum-=v[l];
  l++;
}
ans=max(ans,sum);
}
cout<<max(ans,v.front())<<'\n';
}
return 0;
}