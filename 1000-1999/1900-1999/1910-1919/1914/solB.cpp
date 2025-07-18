#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef pair<ll,ll> pi;
void solve(){
	int n,k; cin>>n>>k;
	k++;
	for(int i=n+1-k; i<=n; i++) cout<<i<<" ";
	for(int i=n-k; i>=1; i--) cout<<i<<" ";
	cout<<'\n';
}
int main(){
	ios_base::sync_with_stdio(false); cin.tie(NULL);
	int t = 1;
	cin >> t;
	while(t--) solve();
}