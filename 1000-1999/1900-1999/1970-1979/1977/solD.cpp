#include<bits/stdc++.h>
using namespace std;
typedef long long LL;
const int N=3e5+5,P=13331,M1=1e9+7,M2=998244353;
string s[N];
int p1[N],p2[N];
void solve(){
			int n,m;
			cin>>n>>m;
			for(int i=0;i<n;++i){
				cin>>s[i];
			}
			p1[0]=p2[0]=1;
			for(int i=1;i<n;++i){
				p1[i]=(1ll*p1[i-1]*P)%M1;
				p2[i]=(1ll*p2[i-1]*P)%M2;
			}
			map<array<int,2>,int> mp;
			for(int j=0;j<m;++j){
				int h1=0,h2=0;
				for(int i=0;i<n;++i){
					h1=(1ll*h1*P+s[i][j])%M1;
					h2=(1ll*h2*P+s[i][j])%M2;
				}
				for(int i=0;i<n;++i){
					array<int,2> t;
					t[0]=(h1+1ll*((s[i][j]^1)-s[i][j]+M1)*p1[n-1-i])%M1;
					t[1]=(h2+1ll*((s[i][j]^1)-s[i][j]+M2)*p2[n-1-i])%M2;
					mp[t]++;
				}
			}
			int ret=0;
			for(auto &it:mp){
				ret=max(ret,it.second);
			}
			cout<<ret<<'\n';
			for(int j=0;j<m;++j){
				int h1=0,h2=0;
				for(int i=0;i<n;++i){
					h1=(1ll*h1*P+s[i][j])%M1;
					h2=(1ll*h2*P+s[i][j])%M2;
				}
				for(int i=0;i<n;++i){
					array<int,2> t;
					t[0]=(h1+1ll*((s[i][j]^1)-s[i][j]+M1)*p1[n-1-i])%M1;
					t[1]=(h2+1ll*((s[i][j]^1)-s[i][j]+M2)*p2[n-i-1])%M2;
					if(mp[t]==ret){
						s[i][j]^=1;
						for(int i=0;i<n;++i){
							cout<<s[i][j];
						}
						cout<<'\n';
						return;
					}
				}
			}
}
int main(){
	ios::sync_with_stdio(false);
	cin.tie(nullptr);
	int num;
	cin>>num;
	while(num--){
		solve();
	}
	return 0;
}