# include <stdio.h>
# include <algorithm>

using namespace std;

int n,k,q=0;

int n1[100050];

int main()
{
	scanf("%d",&n);
	
	for(int i=0;i<n;i++)
	{
		scanf("%d",&k);
		
		n1[q++]=k;
	}
	
	sort(n1,n1+q);
	
	int q1=0;
	
	for(int i=0;i<q-1;i++)
	if(n1[i]*2 > n1[i+1] && n1[i]!=n1[i+1]) {printf("YES\n");q1=1;break;}
	
	if(q1==0) printf("NO\n");
}