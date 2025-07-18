#include<cstdio>
#include<vector>
#include<algorithm>
using namespace std;
const int INF=0x3f3f3f3f;
const int N=200005;

char d[N];
bool vis[N];
int pos[11],st[11];

int main() {
	int T,n;
	scanf("%d",&T);
	for(int Cas=1;Cas<=T;Cas++) {
		scanf("%d",&n);
		scanf("%s",d);
		for(int i=0;i<=9;i++)
			pos[i]=-1,st[i]=INF;
		for(int i=0;i<n;i++) {
			pos[d[i]-'0']=i,vis[i]=true;
			if(st[d[i]-'0']==INF)
				st[d[i]-'0']=i;
		}
		for(int now=0,i=0,j;now<=9;now++) {
			for(j=i;j<=pos[now];j++)
				if(d[j]=='0'+now)
					vis[j]=false;
			if(st[now]<i) break;
			i=j;
		}
		int last=0;
		for(int i=0;i<n;i++)
			if(vis[i]) {
				if(d[i]-'0'>=last)
					last=d[i]-'0';
				else {
					last=233;
					break;
				}
			}
		if(last==233)
			printf("-");
		else for(int i=0;i<n;i++)
			putchar('1'+vis[i]);
		puts("");
	}
	return 0;
}