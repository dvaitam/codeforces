#include<cstdio>
#include<cmath>
#include<algorithm>
#include<cstring>
#include<queue>
#include<set>
using namespace std;
char *p1,*p2,buf[10010];
#define gc() (p1==p2&&(p2=(p1=buf)+fread(buf,1,10000,stdin),p1==p2)?EOF:*p1++)
int read(){
	int x=0;
	char c=gc();
	while(c<48)c=gc();
	while(c>47)x=(x<<3)+(x<<1)+(c&15),c=gc();
	return x;
}
bool vs[8510];
int a[8510],nxt[8510],ans[8510];
int main(){
	set<int>s;
	set<int>::iterator it;
	int t=read(),n,col;
	while(t--){
		a[0]=(n=read())+1;
		for(int i=1;i<=n;++i)a[i]=read();
		reverse(a+1,a+n+1);
		for(int i=1;i<=n;++i)nxt[a[i-1]]=a[i];
		nxt[a[n]]=col=0;
		while(nxt[n+1]){
			++col;
			for(int i=n+1;nxt[i];i=nxt[i]){
				vs[nxt[i]]=0;
				s.insert(nxt[i]);
			}
			for(int i=n+1;nxt[i];){
				if(!vs[nxt[i]]){
					it=s.find(nxt[i]);
					if(it!=s.begin()){
						vs[*--it]=1;
						++it;
					}
					if(++it!=s.end())vs[*it]=1;
					ans[nxt[i]]=col;
					nxt[i]=nxt[nxt[i]];
					s.erase(--it);
				}else s.erase(i=nxt[i]);
			}
		}
		for(int i=1;i<=n;++i)printf("%d ",ans[i]);
		putchar('\n');
	}
	return 0;
}