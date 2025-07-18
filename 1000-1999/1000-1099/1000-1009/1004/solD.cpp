#include<bits/stdc++.h>
#define LL long long
#define ull unsigned long long
#define ULL ull
#define mp make_pair
#define pii pair<int,int>
#define piii pair<int, pii >
#define pll pair <ll,ll>
#define pb push_back
#define big 20160116
#define INF 2147483647
#define pq priority_queue
using namespace std;
inline int read(){
	int x=0,f=1;
	char ch=getchar();
	while (ch<'0'||ch>'9'){if(ch=='-') f=-1;ch=getchar();}
	while (ch>='0'&&ch<='9'){x=x*10+ch-'0';ch=getchar();}
	return x*f;
}
namespace Mymath{
	LL qp(LL x,LL p,LL mod){
		LL ans=1;
		while (p){
			if (p&1) ans=ans*x%mod;
			x=x*x%mod;
			p>>=1;
		}
		return ans;
	}
	LL inv(LL x,LL mod){
		return qp(x,mod-2,mod);
	}
	LL C(LL N,LL K,LL fact[],LL mod){
		return fact[N]*inv(fact[K],mod)%mod*inv(fact[N-K],mod)%mod;
	}
	template <typename Tp> Tp gcd(Tp A,Tp B){
		if (B==0) return A;
		return gcd(B,A%B);
	}
	template <typename Tp> Tp lcm(Tp A,Tp B){
		return A*B/gcd(A,B);
	}
};
namespace fwt{
	using namespace Mymath;
	void FWT(int a[],int n,LL mod)
	{
	    for(int d=1;d<n;d<<=1)
	        for(int m=d<<1,i=0;i<n;i+=m)
	            for(int j=0;j<d;j++)
	            {
	                int x=a[i+j],y=a[i+j+d];
	                a[i+j]=(x+y)%mod,a[i+j+d]=(x-y+mod)%mod;
	                //xor:a[i+j]=x+y,a[i+j+d]=x-y;
	                //and:a[i+j]=x+y;
	                //or:a[i+j+d]=x+y;
	            }
	}

	void UFWT(int a[],int n,LL mod)
	{
		LL rev=inv(2,mod);
	    for(int d=1;d<n;d<<=1)
	        for(int m=d<<1,i=0;i<n;i+=m)
	            for(int j=0;j<d;j++)
	            {
	                int x=a[i+j],y=a[i+j+d];
	                a[i+j]=1LL*(x+y)*rev%mod,a[i+j+d]=(1LL*(x-y)*rev%mod+mod)%mod;
	                //xor:a[i+j]=(x+y)/2,a[i+j+d]=(x-y)/2;
	                //and:a[i+j]=x-y;
	                //or:a[i+j+d]=y-x;
	            }
	}
	void solve(int a[],int b[],int n,LL mod)
	{
	    FWT(a,n,mod);
	    FWT(b,n,mod);
	    for(int i=0;i<n;i++) a[i]=1LL*a[i]*b[i]%mod;
	    UFWT(a,n,mod);
	}
};
const int Maxn=1e6+5;
int cnt[Maxn];
int c2[Maxn];
int t;
int Gbl(int x,int y,int X,int Y){
	return min(min(X-1,Y-1),min(x-X,y-Y));
} 
void Check(int x,int y,int X,int Y){
	memset(c2,0,sizeof(c2));
	for (int i=1;i<=x;i++){
		for (int j=1;j<=y;j++){
			c2[abs(X-i)+abs(Y-j)]++;
		}
	}
	for (int i=1;i<=t;i++){
		if (c2[i]!=cnt[i]) return;
	}
	printf("%d %d\n%d %d\n",x,y,X,Y);exit(0);
}
int main(){
	t=read();
	int Mx=0;
	for (int i=0;i<t;i++){
		int x=read();cnt[x]++;
		Mx=max(Mx,x);
	}
	int Bl=0;
	for (int i=1;i<=t;i++){
		if (cnt[i]!=i*4) {
			Bl=i-1;
			break;
		}
	}//cout<<Bl<<endl;
	for (int i=1;i*i<=t;i++){
		if (t%i==0){
			int n=i,m=t/i;
			if (n+m-2<Mx) continue;
		//	cout<<n<<m<<endl;
			for (int j=1;j<=n;j++){
				int k=Mx-j+2;
				//cout<<j<<k<<endl;
				if (k>m || k<=0) continue;
				//cout<<n<<m<<j<<k<<endl;
				if (Gbl(n,m,j,k)!=Bl) continue;
				
				Check(n,m,j,k);
			}
		}
	}
	puts("-1");
	return 0;
}