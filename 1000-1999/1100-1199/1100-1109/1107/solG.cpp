//ios::sync_with_stdio(false);
#include<bits/stdc++.h>
#define LL long long
#define F(x,y,z) for(int x=y;x<=z;++x)
#define D(x,y,z) for(int x=y;x>=z;--x)
using namespace std;
const int N=300010;
const LL Inf=100000000000000ll;
LL Min(LL x,LL y){return x<y?x:y;}
LL Max(LL x,LL y){return x>y?x:y;}
LL R(){
    LL ans=0,f=1;char c=' ';
    for(;c<'0'||c>'9';c=getchar()) if (c=='-') f=-1;
    for(;c>='0'&&c<='9';c=getchar()) ans=ans*10+c-'0';
    return ans*f;
}
int n,m;
LL dd[N],w[N],ans,d[N];
LL stk[N],tp,pla[N],a1[N],a2[N];
LL st1[20][N],st2[20][N];
LL Get1(int l,int r){
	int ci=0;
	while((1<<(ci+1))<(r-l+1))++ci;
	return Min(st1[ci][l],st1[ci][r-(1<<ci)+1]);
}
LL Get2(int l,int r){
	int ci=0;
	while((1<<(ci+1))<(r-l+1))++ci;
	return Max(st2[ci][l],st2[ci][r-(1<<ci)+1]);
}
int main(){
	n=R();m=R();
	F(i,1,n)dd[i]=R(),w[i]=m-R();
	F(i,1,n)ans=Max(ans,w[i]);
	F(i,1,n-1)d[i]=dd[i+1]-dd[i];
	stk[0]=Inf;
	F(i,1,n-1){
		while(stk[tp]<=d[i])--tp;
		a1[i]=pla[tp]+1;
		++tp;stk[tp]=d[i];pla[tp]=i;
	}
	tp=0;pla[tp]=n;
	D(i,n-1,1){
		while(stk[tp]<=d[i])--tp;
		a2[i]=pla[tp];
		++tp;stk[tp]=d[i];pla[tp]=i;
	}
	F(i,1,n)w[i]+=w[i-1];
	F(i,1,n)st1[0][i]=st2[0][i]=w[i];
	F(i,1,18)F(j,0,n)
		st1[i][j]=Min(st1[i-1][j],st1[i-1][j+(1<<(i-1))]);
	F(i,1,18)F(j,0,n)
		st2[i][j]=Max(st2[i-1][j],st2[i-1][j+(1<<(i-1))]);
	F(i,1,n-1)
		ans=Max(ans,Get2(i+1,a2[i])-Get1(a1[i]-1,i-1)-d[i]*d[i]);
	cout<<ans<<endl;
    return 0;
}