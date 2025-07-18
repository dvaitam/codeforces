#include<bits/stdc++.h>
using namespace std;
#define pii pair<int,int>
#define ff first
#define ss second
#define sz(x) ((int)(x.size()))
const int maxn=205;
const int B=100;
int O;
int a[maxn];
pii s[maxn];
bool gz(pii x,pii y){
	return x.ff%O<y.ff%O;
}
int main(){
	int n,k;
	srand(time(0));
	cin>>n>>k;
	for(int i=1;i<=n;i++){
		cin>>a[i];
	}
	vector<bool>vis(n+5);
	for(int i=1;i<=k;i++){
		cin>>s[i].ff;
		s[i].ss=i;
	}
	sort(s+1,s+1+k);
	vector<vector<int> >g(k+5);
	for(int i=1;i<k;i++){
		int id=s[i].ss;
		vector<pii>_;
		_.push_back({-1,0});
		O=s[i].ff;
		for(int j=1;j<=n;j++){
			if(!vis[j]){
				_.push_back({a[j],j});
			}
		}
		int m=sz(_)-1,c=0;
		bool ok=0;
		while(!ok){
			c++;
			if(c==1){
				sort(_.begin(),_.end(),gz);
			}else{
				for(int j=1;j<=20;j++){
					int x=rand()%m+1,y=rand()%m+1;
					swap(_[x],_[y]);
				}
			}
			vector<int>pre(m+5);
			for(int j=1;j<=m;j++){
				pre[j]=pre[j-1]+_[j].ff;
			}
			for(int j=s[i].ff;j<=m;j++){
				if((pre[j]-pre[j-s[i].ff])%s[i].ff==0){
					ok=1;
					for(int k=j-s[i].ff+1;k<=j;k++){
						g[id].push_back(_[k].ss);
						vis[_[k].ss]=1;
					}
					break;
				}
			}
		}
	}
	int sum=0;
	for(int i=1;i<=n;i++){
		if(!vis[i]){
			sum+=a[i];
			g[s[k].ss].push_back(i);
		}
	}
	int ljdlovesxyq=s[k].ff-sum%s[k].ff;
	g[s[k].ss].push_back(n+1);
	a[n+1]=ljdlovesxyq;
	cout<<ljdlovesxyq<<endl;
	for(int i=1;i<=k;i++){
		for(int t:g[i]){
			cout<<a[t]<<" ";
		}
		cout<<endl;
	}
	return 0;
}