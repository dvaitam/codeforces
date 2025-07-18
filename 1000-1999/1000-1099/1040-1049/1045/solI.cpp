#include<cstdio>
#include<cstring>
#include<iostream>
#include<stdlib.h>
#include<ctime>
#include<string>
#include<cmath>
#include<algorithm>
#include<vector>
#include<set>
#include<map>
#include<queue>
#define pb push_back
#define mp make_pair
#define fi first
#define se second
#define LL long long
#define FOR(i,a,b) for (int i=a;i<=b;i++)
#define FORD(i,a,b) for (int i=a;i>=b;i--)
using namespace std;
void getint(int &v){
    char ch,fu=0;
    for(ch='*'; (ch<'0'||ch>'9')&&ch!='-'; ch=getchar());
    if(ch=='-') fu=1, ch=getchar();
    for(v=0; ch>='0'&&ch<='9'; ch=getchar()) v=v*10+ch-'0';
    if(fu) v=-v;
}
const int MO=100007;
char s[1000010];
int n,c[1000010],m,a[1000010],too[1000010],nedge,val[1000010],nxt[1000010],hed[1000010];
LL ans;
void ae(int x,int t){
	for (int i=hed[x];i;i=nxt[i]){
		int y=too[i];
		if (y==t){++c[i];return;}
	}
	nxt[++nedge]=hed[x];
	hed[x]=nedge;
	too[nedge]=t;
	c[nedge]=1;
}
int fnd(int x){
	for (int i=hed[x%MO];i;i=nxt[i]){
		int y=too[i];
		if (y==x) return c[i];
	}
	return 0;
}
int main(){
	cin>>n;
	FOR(i,1,n){
		scanf("%s",s+1);
		m=strlen(s+1);
		FOR(j,1,m) a[i]^=1<<(s[j]-'a');
	}
	FOR(i,1,n){
		ans+=fnd(a[i]);
		FOR(j,0,25) ans+=fnd(a[i]^(1<<j));
		ae(a[i]%MO,a[i]);
	}
	cout<<ans<<endl;
    return 0;
}