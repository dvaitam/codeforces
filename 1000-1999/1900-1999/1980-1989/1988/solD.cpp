#include<bits/stdc++.h>
using namespace std; 
vector<int> v[300005];
int T,n;
long long a[300005],dp[300005][20];
void dfs(int x,int fa){
	for(int i=1;i<20;++i) dp[x][i]=a[x]*i;
	for(int u : v[x]){
		if(u==fa)continue;
		dfs(u,x);
		for(int j=1;j<20;++j){
			long long sum=0x3f3f3f3f3f3f3f3f;
			for(int k=1;k<20;++k)
				if(j!=k) sum=min(sum,dp[u][k]);
			dp[x][j]+=sum;
		}
	}
}
int main(){
	ios::sync_with_stdio(0),cin.tie(0),cout.tie(0); 
	cin>>T;
	while(T--){
		cin>>n;
		for(int i=1;i<=n;++i) cin>>a[i],v[i].clear();
		for(int i=1;i<n;++i){
			int x,y;
			cin>>x>>y;
			v[x].push_back(y),v[y].push_back(x);
		}
		dfs(1,0);
		cout<<*std::min_element(dp[1]+1,dp[1]+20)<<endl;
	}
	return 0;
}