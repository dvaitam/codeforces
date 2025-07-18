#include<bits/stdc++.h>
using namespace std;
//#pragma GCC optimize(2)
#ifdef ONLINE_JUDGE
#define filename(_ST_) \
		freopen(_ST_".in","r",stdin); \
		freopen(_ST_".out","w",stdout);
#else
#define filename(_ST_)
#endif
char outp[64];
template<typename T>
inline void read(T &x){
	char ch=0; x=0;
	bool sign=false;
	while(ch<'0'||'9'<ch) sign|=ch=='-',ch=getchar();
	while('0'<=ch&&ch<='9') x=(x<<3)+(x<<1)+(ch^48),ch=getchar();
	x=sign?-x:x;
}
template<typename T>
inline void print(T x){
	if(!x){ putchar('0'); return; }
	if(x<0) putchar('-'),x=-x;
	int tot=0;
	while(x) outp[tot++]=x%10+'0',x/=10;
	while(tot) putchar(outp[--tot]);
}
typedef long long ll;
typedef pair<ll,ll> pll;
typedef pair<int,int> pii;
#define fi first
#define se second
#define pb push_back
#define mp make_pair
#define rt register int
#define read2(_a,_b) read(_a),read(_b)
#define read3(_a,_b,_c) read(_a),read(_b),read(_c)
#define For1(f,e,item) for(rt item=f;item!=e;++item)
#define Rep1(f,e,item) for(rt item=f;item!=e;--item)
#define For2(f,e,item) for(rt item=f;item<e;++item)
#define Rep2(f,e,item) for(rt item=f;item>e;--item)
#define For3(f,e,item) for(rt item=f;item<=e;++item)
#define Rep3(f,e,item) for(rt item=f;item>=e;--item)
typedef long long ll;
const int MAXN=2e5+128;
array<ll,MAXN> a;
array<ll,MAXN/2> b;
int main(){
	//filename("");
	//ios::sync_with_stdio(false);
	int n;
	read(n);
	For3(1,n>>1,i){
		read(b[i]);
	}
	a[1]=0;
	a[n]=b[1];
	For3(2,n>>1,i){
		int j=n-i+1;
		a[i]=a[i-1];
		a[j]=b[i]-a[i];
		if(a[j]>a[j+1]){
			a[i]+=a[j]-a[j+1];
			a[j]=a[j+1];
		}
	}
	For3(1,n,i){
		print(a[i]);
		putchar(i==n?'\n':' ');
	}
	return 0;
}