#include<bits/stdc++.h>
using namespace std;


int main()
{
//	freopen("/home/zz7/CFInput","r",stdin);
	int n,m;scanf("%d%d",&n,&m);
	int M=0,sum=0;
	for(int i=0;i<n;i++){
		int x;
		scanf("%d",&x);M=max(x,M);
		sum+=x;
	}
	int u=0;if((sum+m)%n!=0)u=1;
	cout<<max((sum+m)/n+u,M)<<' '<<M+m<<'\n';
}