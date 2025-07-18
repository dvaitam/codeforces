#include <bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef unsigned int uint;
typedef unsigned long long ull;
const ll mod=998244353;
const ll inf=2e9;
const int N=2e5+5;
const int M=2e5+5;
const int V=2e6+5;
int n,x,c[105];
double dp[105][10005],C[105][105];
void solve(int Ca){
	cin>>n>>x;
	for(int i=1;i<=n;i++) cin>>c[i];
	for(int i=0;i<=n;i++){
		C[i][0]=C[i][i]=1;
		for(int j=1;j<i;j++) C[i][j]=C[i-1][j]+C[i-1][j-1];
	}
	dp[0][0]=1;
	int now=0;
	for(int i=1;i<=n;i++){
		now+=c[i];
		for(int j=i;j>=1;j--) for(int p=now;p>=c[i];p--) dp[j][p]+=dp[j-1][p-c[i]];
	}
	double ans=0;
	for(int i=1;i<=n;i++){
		for(int j=1;j<=now;j++) ans+=dp[i][j]/C[n][i]*min((n*1.0/i+1)*x/2,j*1.0/i);
	}
	cout<<fixed<<setprecision(10)<<ans<<"\n";
}
int main(){
	#ifdef ONLINE_JUDGE
	ios::sync_with_stdio(false); cin.tie(nullptr); cout.tie(nullptr);
	#endif
	#ifndef ONLINE_JUDGE
	freopen("test.in","r",stdin);
	freopen("test.out","w",stdout);
	#endif
	
	int Ca=1;
//	cin>>Ca;
	for(int i=1;i<=Ca;i++){
		solve(i);
	}
	return 0;
}