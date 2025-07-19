#include<bits/stdc++.h>
#define ll long long
using namespace std;

void solve() {
    ll n,x,ans=0,temp=0;
    cin>>n>>x;
    n--; 
	while(n--) {
		ll y;
		cin>>y;
		if(y==1 && x!=1) {
			ans=-1;
		}
		if(ans==-1) {
			continue;
		}
		ans+=temp=max((ll)0,temp+(ll)ceil(log2(log(x)/log(y))));
		x=y;
	}
    cout<<ans<<"\n";
}

int main() {
    ios::sync_with_stdio(false);
    cin.tie(0), cout.tie(0);
    ll t;
    cin >> t;
    while (t--) {
        solve();
    }
    
    return 0;
}