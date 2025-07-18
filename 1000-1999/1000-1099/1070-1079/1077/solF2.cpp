#include <bits/stdc++.h>
#include <bits/stdc++.h>
#define ll long long
using namespace std;

inline int read(){
	int x=0,f=1;char cc=getchar();
	while(cc<'0' || cc>'9') {if(cc=='-') f=-1;cc=getchar();}
	while(cc>='0' && cc<='9') {x=x*10+cc-'0';cc=getchar();}
	return x*f;
}

int n,k,x,a[5010],dui[5010],t1,t2;
ll f[5010],g[5010],ans;

int main(){
	n=read();k=read();x=read();
	for(int i=1;i<=n;i++) a[i]=read();
	if(n/k>x){
		printf("-1");
		return 0;
	}
	for(int i=1;i<=k;i++) g[i]=a[i];
	for(int l=2;l<=x;l++){
		for(int i=1;i<=n;i++) f[i]=0;
		t1=1;t2=0;
		for(int i=l;i<=min(n,l*k);i++){
			while(t1<=t2) if(g[i-1]>=g[dui[t2]]) t2--;else break;
			t2++;dui[t2]=i-1;
			while(dui[t1]<i-k) t1++;
			f[i]=g[dui[t1]]+a[i];
		}
		for(int i=1;i<=n;i++) g[i]=f[i];
	}
	ans=0;
	for(int i=max(1,n-k+1);i<=n;i++) ans=max(ans,g[i]);
	printf("%lld",ans);
}