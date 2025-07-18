#include<cstdio>
#include<iostream>
#include<algorithm>
#include<vector>
#include<set>
#include<map>
#include<queue>
#include<cmath>
#include<cstdlib>
#include<string>
#include<cstring>
#include<bitset>
#define LL long long
#define mod 1e9+7
#define INF 0x3f3f3f3f
using namespace std;

namespace FastIO {
	template<typename tp> inline void read(tp &x) {
		x=0; register char c=getchar(); register bool f=0;
		for(;c<'0'||c>'9';f|=(c=='-'),c = getchar());
		for(;c>='0'&&c<='9';x=(x<<3)+(x<<1)+c-'0',c = getchar());
		if(f) x=-x;
	}
	template<typename tp> inline void write(tp x) {
		if (x==0) return (void) (putchar('0'));
		if (x<0) putchar('-'),x=-x;
		int pr[20]; register int cnt=0;
		for (;x;x/=10) pr[++cnt]=x%10;
		while (cnt) putchar(pr[cnt--]+'0');
	}
	template<typename tp> inline void writeln(tp x) {
		write(x);
		putchar('\n');
	}
}
using namespace FastIO;
int n,m;
int a[100010],b[100010];
LL total=0,tot=0;
int main(){
	read(n),read(m);
	for(int i=1;i<=n;++i) read(a[i]),read(b[i]),total+=b[i],tot+=a[i],b[i]=a[i]-b[i];
	if(total>m) return puts("-1"),0;
	sort(b+1,b+n+1);
	int ans=0,biao=n;
	while(tot>m){
		tot-=b[biao--];
		ans++;
	}
	writeln(ans);
	return 0;
}