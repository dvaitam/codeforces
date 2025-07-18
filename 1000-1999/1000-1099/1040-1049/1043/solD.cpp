#include<iostream>
#include<cstdio>
#include<cstring>
#include<ctime>
#include<cstdlib>
#include<algorithm>
#include<cmath>
#include<string>
#include<queue>
#include<vector>
#include<map>
#include<set>
#include<utility>
#include<iomanip>
using namespace std;
int read(){
    int xx=0,ff=1;char ch=getchar();
    while(ch>'9'||ch<'0'){if(ch=='-')ff=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){xx=xx*10+ch-'0';ch=getchar();}
    return xx*ff;
}
long long READ(){
    long long xx=0,ff=1;char ch=getchar();
    while(ch>'9'||ch<'0'){if(ch=='-')ff=-1;ch=getchar();}
    while(ch>='0'&&ch<='9'){xx=xx*10+ch-'0';ch=getchar();}
    return xx*ff;
}
char one(){
	char ch=getchar();
	while(ch==' '||ch=='\n')
		ch=getchar();
	return ch;
}
const int maxn=100010;
int N,M,a[11][maxn];
int father[maxn],cnt[maxn];
inline int getfather(int x){
	if(x==father[x])return x;
	return father[x]=getfather(father[x]);
}
int main(){
	//freopen("in","r",stdin);
	N=read(),M=read();
	for(int i=1;i<=M;i++)
		for(int j=1;j<=N;j++)
			a[i][j]=read();
	for(int i=1;i<=N;i++)
		father[a[1][i]]=a[1][i+1];
	for(int i=2;i<=M;i++)
		for(int j=1;j<=N;j++)
			if(father[a[i][j]]!=a[i][j+1])
				father[a[i][j]]=0;
	for(int i=1;i<=N;i++)
		if(!father[i])
			father[i]=i;
	for(int i=1;i<=N;i++)
		if(getfather(i))
			cnt[father[i]]++;
	long long ans=N;
	for(int i=1;i<=N;i++)
		if(cnt[i]>1)
			ans+=1LL*cnt[i]*(cnt[i]-1)/2;
	cout<<ans<<endl;
	return 0;
}