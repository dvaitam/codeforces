#include<bits/stdc++.h>
using namespace std;
int a[500005];
vector<int>g[500005];
long long S[500005];
void add(int l,int r,long long x){
	S[l]+=x;S[r+1]-=x;
}
int main(){
	int T;
	scanf("%d",&T);
	while(T--)
	{
		int n;
		scanf("%d",&n);
		S[0]=0;
		for(int i=1;i<=n;++i) g[i].clear(),S[i]=0;
		for(int i=1;i<=n;++i) scanf("%d",&a[i]),g[a[i]].emplace_back(i);
		set<int>se;
		se.emplace(0);se.emplace(n+1);
		for(int i=1;i<=n;++i)
		{
			for(auto x:g[i])
			{
				auto it=se.emplace(x).first;
				auto t1=it,t2=it;--t1;++t2;
				int L=*t1,R=*t2;
				add(0,L-1,1ll*(x-L)*(R-x)*a[x]);
				add(R+1,n,1ll*(x-L)*(R-x)*a[x]);
				add(L+1,x-1,1ll*(x-L-1)*(R-x)*a[x]);
				add(x+1,R-1,1ll*(x-L)*(R-x-1)*a[x]);
				if(L){
					--t1;int LL=*t1;
					add(L,L,1ll*(x-LL-1)*(R-x)*a[x]);
				}
				if(R!=n+1){
					++t2;int RR=*t2;
					add(R,R,1ll*(x-L)*(RR-x-1)*a[x]);
				}
			}
		}	
		for(int i=1;i<=n;++i)printf("%lld ",S[i]+=S[i-1]);
		printf("\n");
	}
	return 0;
}