#include <bits/stdc++.h>
#define endl '\n'
using namespace std;
const int INF = 1e9;
const int N = 1e5+10;
typedef double db;
typedef long long loint;
typedef unsigned long long unlo;
const loint mod = 998244353;

void solve(){
    int n;cin>>n;
    vector<int>a(n),b(n);
    int mi = INF,mx = -INF;
    for(int i = 0;i<n;i++){
        cin>>a[i];
    }
    for(int i = 0;i<n;i++){
        cin>>b[i];
        mi = min(a[i]-b[i], mi);
        mx = max(a[i]+b[i], mx);
    }
    cout<<0.5*(mi+mx)<<endl;
}

int main(){
    ios::sync_with_stdio(false);
    cin.tie(nullptr);
    cout<<fixed;
//    solve();
//    freopen("", "r", stdin);
//    freopen("", "w", stdout);
    int t;cin>>t;
    while(t--){
        solve();
    }
    return 0;
}