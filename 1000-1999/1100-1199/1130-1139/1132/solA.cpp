#include<bits/stdc++.h>
using namespace std;
int n1,n2,n3,n4;
int main()
{
	scanf("%d%d%d%d",&n1,&n2,&n3,&n4);
	if(n1!=n4) puts("0");
	else
	{
		if(n1==0)
		{
			if(n3!=0) puts("0");
			else puts("1");
		}
		else puts("1");
	}
	
	return 0;
}