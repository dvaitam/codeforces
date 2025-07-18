#include<cstdio>
#include<cstring>
#include<vector>
namespace Qza{
	#define MAXN 100005
	std::vector<int> pri[MAXN];
	bool used[MAXN];
	int read()
	{
		char ch=getchar(); int x=0;
		while(ch<'0' || ch>'9') ch=getchar();
		while(ch>='0'&&ch<='9') x=(x*10)+(ch^48),ch=getchar();
		return x;
	}
	void main()
	{
		int t,n,m,prince;
		bool flag; 
		t=read(); 
		while(t--) {
			prince=0;
			n=read();
			memset(used,0,sizeof(int)*(n+5));
			for(int i=1;i<=n;++i) {
				flag=false;
				m=read();
				pri[i].clear(); 
				for(int j=1;j<=m;++j) {
						pri[i].push_back(read());
				}
				for(int j=0;j<m;++j) {
					if(!used[pri[i][j]]) {
						used[pri[i][j]]=1; flag=true; break;
					}
				}
				if(!flag) prince=i;
			}
			if(!prince) printf("OPTIMAL\n");
			else {
				for(int i=1;i<=n;++i) if(!used[i]) {
					printf("IMPROVE\n%d %d\n",prince,i); break;
				}
			}
		}
	}
}
int main()
{
	Qza::main();
	return 0;
}