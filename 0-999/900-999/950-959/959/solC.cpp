#include<cstdio>
#include<algorithm>
#include<cstdlib>
#include<cstring>
#include<memory.h>
#include<map>
#include<set>
#include<queue>
using namespace std;
int main()
{
	int n;
	scanf("%d",&n);
	if(n<=5)puts("-1");else
	{
		int tmp=n-2;
		printf("1 2\n");
		int tot=tmp/2;
		for(int i=1;i<=tot;i++)printf("1 %d\n",i+2);
		for(int i=1;i<=tmp-tot;i++)printf("2 %d\n",i+tot+2);
	}
	for(int i=2;i<=n;i++)printf("1 %d\n",i);
}