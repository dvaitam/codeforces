#include<bits/stdc++.h>
typedef int LL;
const LL maxn=3e6+9;
inline LL Read(){
	LL x(0),f(1); char c=getchar();
	while(c<'0' || c>'9'){
		if(c=='-') f=-1; c=getchar();
	}
	while(c>='0' && c<='9'){
		x=(x<<3)+(x<<1)+c-'0'; c=getchar();
	}
	return x*f;
}
LL n,tot;
LL prime[maxn],visit[maxn],sum[maxn],ys[maxn],a[maxn],b[maxn],ans[maxn],low[maxn];
void Fir(){
//	prime[++tot]=2;
	for(LL i=2;i<=3000000;++i){
		if(!visit[i]){
			prime[++tot]=i;
		}
		for(LL j=1;j<=tot && prime[j]*i<=3000000;++j){
			visit[prime[j]*i]=true;
			if(i%prime[j]==0) break;
		}
	}
}
int main(){
	Fir();
	for(LL i=1;i<=tot;++i) ys[prime[i]]=true;
//	printf("%d\n",prime[199999]);
//    for(LL i=1;i<=7;++i) printf("%d ",prime[i]);
	n=Read();
	for(LL i=1;i<=n*2;++i){
		LL val(Read());
//		++sum[val];
		if(ys[val]) a[++a[0]]=val;
		else b[++b[0]]=val;
	}
	std::sort(b+1,b+1+b[0]);
	for(LL i=b[0];i>=1;--i){
		if(low[b[i]]) low[b[i]]--;
		else{
			ans[++ans[0]]=b[i];
			for(LL j=1;j<=tot;++j){
				if(b[i]%prime[j]==0){
					++low[b[i]/prime[j]];
					break;
				}
			}
		}
	}
	std::sort(a+1,a+1+a[0]);
	for(LL i=1;i<=a[0];++i){
		if(low[a[i]]) low[a[i]]--;
		else{
			ans[++ans[0]]=a[i];
			low[prime[a[i]]]++;
		}
	}
//	for(LL i=1;i<=a[0];++i) printf("%d ",a[i]);
	for(LL i=1;i<=n;++i) printf("%d ",ans[i]);
	return 0;
}