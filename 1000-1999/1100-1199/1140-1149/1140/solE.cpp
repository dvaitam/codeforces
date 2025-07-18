#include<cstdio>
#include<cstring>
#include<algorithm>
#include<queue>
#include<cctype>
#define qmin(x,y) (x=min(x,y))
#define qmax(x,y) (x=max(x,y))
#define vi vector<int>
#define vit vector<int>::iterator
#define pir pair<int,int>
#define fr first
#define sc second
#define mp(x,y) make_pair(x,y)
#define rsort(x,y) sort(x,y),reverse(x,y)
using namespace std;
 
inline char gc() {
//  static char buf[100000],*p1,*p2;
//  return (p1==p2)&&(p2=(p1=buf)+fread(buf,1,100000,stdin),p1==p2)?EOF:*p1++;
    return getchar();
}
 
template<class T>
int read(T &ans) {
    ans=0;char ch=gc();T f=1;
    while(!isdigit(ch)) {
        if(ch==EOF) return -1;
        if(ch=='-') f=-1;
        ch=gc();
    }
    while(isdigit(ch))
        ans=ans*10+ch-'0',ch=gc();
    ans*=f;return 1;
}
 
template<class T1,class T2>
int read(T1 &a,T2 &b) {
    return read(a)!=EOF&&read(b)!=EOF?2:EOF;
}
 
template<class T1,class T2,class T3>
int read(T1 &a,T2 &b,T3 &c) {
    return read(a,b)!=EOF&&read(c)!=EOF?3:EOF;
}
 
typedef long long ll;
const int Maxn=1100000;
const int inf=0x3f3f3f3f;
const int mod=998244353;

int n,k,a[Maxn],b[Maxn],num;
ll ans=1,f[Maxn],g[Maxn];

void work() {
	int temp=1;
	for(int i=1;i<=num;i++) {
		if(~b[i]) {
			if(temp==i) {
				temp++;
				continue;
			}
			if(temp==1) ans=ans*(f[i-temp-1]+g[i-temp-1]*(k-1)%mod)%mod;
			else if(b[i]==b[temp-1]) ans=ans*f[i-temp]%mod;
			else ans=ans*g[i-temp]%mod;
			temp=i+1;
		}
	}
	if(temp!=num+1) {
		if(temp!=1) ans=ans*(f[num-temp]+g[num-temp]*(k-1)%mod)%mod;
		else if(num==1) ans=ans*k%mod;
		else ans=ans*(k*f[num-2]%mod+1ll*k*(k-1)%mod*g[num-2]%mod)%mod;
	}
}

signed main() {
//	freopen("test.in","r",stdin);
    read(n,k);
    if(k==1) {
    	if(n<=2) return 0*puts("1");
    	else return 0*puts("0");
	}
    f[1]=k-1,g[1]=k-2;g[0]=1;
    for(int i=2;i<=n;i++) {
    	f[i]=(k-1)*g[i-1]%mod;
    	g[i]=((k-2)*g[i-1]+f[i-1])%mod;
	}
	for(int i=1;i<=n;i++) read(a[i]);
	for(int i=1;i<=n;i+=2) b[++num]=a[i];
	for(int i=1;i<num;i++) if(~b[i]&&b[i]==b[i+1]) {
		return 0*puts("0");
	}
	work();num=0;
	for(int i=2;i<=n;i+=2) b[++num]=a[i];
	for(int i=1;i<num;i++) if(~b[i]&&b[i]==b[i+1]) {
		return 0*puts("0");
	}
	work();
	printf("%I64d\n",ans);
    return 0;
}