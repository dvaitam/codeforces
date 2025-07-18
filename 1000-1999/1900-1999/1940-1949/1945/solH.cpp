#include<bits/stdc++.h>
using namespace std;
typedef long long ll;

#define N 400000
const int inf=0x7fffffff;

int i,j,k,n,m,t,a[N+50],li,sb,tot;
int pos[N+50][2],num[N+50];

basic_string<int> v1,v2;

bool fuck(){
	int i,j,k,tot=100;
	for(i=li;i>=1;i--){
		basic_string<int> q;
		for(j=i;j<=li;j+=i){
			if(pos[j][0])q+=pos[j][0];
			if(pos[j][1])q+=pos[j][1];
		}
		for(auto x:q)for(auto y:q)if(x<y){
			tot--;
			k=inf; for(j=1;j<=n;j++)if(j!=x&&j!=y)k&=a[j]; k+=m;
			if(i>k){
				v1={a[x],a[y]};
				for(i=1;i<=n;i++)if(i!=x&&i!=y)v2+=a[i];
				return 1;
			}
			if(tot<=0)return 0;
		}
	}
	return 0;
}

int main(){
	ios::sync_with_stdio(0); cin.tie(0);
	cin>>t;
	while(t--){
		cin>>n>>m; v1=v2={}; sb=inf;
		for(i=1;i<=n;i++)cin>>a[i];
		sort(a+1,a+n+1); li=a[n];
		for(i=0;i<=li;i++){
			pos[i][0]=0; pos[i][1]=0; num[i]=0;
		}
		for(i=1;i<=n;i++){
			num[a[i]]++;
			if(!pos[a[i]][0])pos[a[i]][0]=i;
			else pos[a[i]][1]=i;
			sb&=a[i];
		}
		if(!fuck())cout<<"NO\n";
		else{
			cout<<"YES\n";
			cout<<v1.size()<<' '; for(auto i:v1){cout<<i<<' ';} cout<<'\n';
			cout<<v2.size()<<' '; for(auto i:v2){cout<<i<<' ';} cout<<'\n';
		}
	}
}