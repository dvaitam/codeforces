#include<cstdio>
#include<cstring>
#include<algorithm>
namespace Qza{
	#define MAXN 100005
	int in[MAXN];
	struct E {
		int u,v,ans;
	}e[MAXN];
	int read()
	{
		char ch=getchar();int x=0;
		while(ch<'0' || ch>'9') ch=getchar();
		while(ch>='0'&&ch<='9') x=(x*10)+(ch^48),ch=getchar();
		return x;
	}
	void main()
	{
		int n,node,cnt;
		n=read();
		for(int i=1;i<n;++i) {
			e[i].u=read();
			e[i].v=read();
			e[i].ans=-1;
			++in[e[i].u];++in[e[i].v];
		}
		node=0;
		for(int i=1;i<=n;++i) {
			if(in[i]>=3) {
				node=i;break;
			}
		}
		if(!node) {
			for(int i=1;i<n;++i) printf("%d\n",i-1);
		}
		else {
			cnt=0;
			for(int i=1;i<n;++i) {
				if(e[i].u==node || e[i].v==node) e[i].ans=cnt++;	/* 不要写成两个 u !!! */
			}
			for(int i=1;i<n;++i) {
				if(e[i].ans==-1) e[i].ans=cnt++;
			}
			for(int i=1;i<n;++i) printf("%d\n",e[i].ans);
		}
	}
}
int main()
{
	Qza::main();
	return 0;
}