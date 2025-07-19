#include<bits/stdc++.h>
#define ll long long
#define M(x) (((x)+2*n)%(2*n))
using namespace std;
const int INF = 2e9+1;
void solve(){
	int n;cin>>n;vector<int>v(2*n);
	for(int &e:v)cin>>e;
	if(n%2==0){
		int maxx=0,minn=INF;
		for(int i=0;i<n/2;i++){
			int s[4]={v[2*i]+v[2*i+1],v[2*i]+v[2*i+n+1],v[2*i+n]+v[2*i+n+1],v[2*i+n]+v[2*i+1]};
			sort(s,s+4);
			maxx=max(maxx,s[2]);
			minn=min(minn,s[1]);
		}
		cout<<maxx-minn<<'\n';
		return;
	}
	if(n==1){
		puts("0");
		return;
	}
	vector<int>r;
	int cnt=0;
	for(int i=0;i<n;i++){
		r.push_back(v[cnt]);
		r.push_back(v[cnt^=1]);
		cnt=M(cnt+n);
	}
	int ans=INF;
	for(int id=0;id<n;id++){
		for(int m1=0;m1<2;m1++){
			for(int m2=0;m2<2;m2++){
				int minn=r[M(2*id-m1)]+r[M(2*id+m2+1)];
				vector<int>dp[2]{vector<int>(n,INF),vector<int>(n,INF)};
				dp[m2][id]=minn;
				for(int j=1;j<n;j++){
					int d2=(id+j)%n,d1=(id+j-1)%n;
					for(int c1=0;c1<2;c1++){
						for(int c2=0;c2<2;c2++){
							if(dp[c1][d1]!=INF&&r[M(2*d2-c1)]+r[M(2*d2+c2+1)]>=minn){
								dp[c2][d2]=min(dp[c2][d2],max(dp[c1][d1],r[M(2*d2-c1)]+r[M(2*d2+c2+1)]));
							}
						}
					}
				}
				int p=(id+n-1)%n;
				if(dp[m1][p]!=INF)ans=min(ans,dp[m1][p]-minn);
			}
		}
	}
	cout<<ans<<'\n';
}
int main(){
	int T;cin>>T;
	while(T--)solve();
	return 0;
}