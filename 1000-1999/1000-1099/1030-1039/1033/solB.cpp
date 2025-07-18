#include <bits/stdc++.h>
using namespace std;
#define int long long
#define ii pair<int, int>
#define vi vector<int>
#define pb push_back
#define sz(x) (int)x.size()
#define all(v) v.begin(), v.end()
#define x first
#define y second
#define rep(i, j, k) for(i=j; i<k; i++)
#define sep(i, j, k) for(i=j; i>k; i--)
#define ios ios_base::sync_with_stdio(false);cin.tie(0);cout.tie(0);
const int inf = 1e9+7;
// const int N = 2e5+5;
int n, k, ans=0;
// int A[N];

bool isPrime(int N) {
	if(N<2)return false;
	if(N<4)return true;
	if((N&1)==0)return false;
	if(N%3==0)return false;
	int curr=5;
	while (curr<=sqrt(N)){
		if(N%curr==0)return false;
		curr+=2;
		if(N%curr==0)return false;
		curr+=4;
	}
	return true;
}

void solve()
{
	int i, j, a, b, x, y;
	cin>>a>>b;
	if(isPrime(a+b) && a-b==1){
		cout<<"YES\n";
	}
	else cout<<"NO\n";
}

signed main()
{
	ios
	int i, t=1, j;
	cin>>t;
	while(t--) 
		solve();
}