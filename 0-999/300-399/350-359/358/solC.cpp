#include <bits/stdc++.h>
using namespace std;
int main()
{
	int i,n;
	scanf("%d",&n);
	vector<int>a;
	int p;
	while(n--)
	{
		scanf("%d",&p);
		if(p==0)
		{
			if(a.size()<3)
			{
				if(a.size()==2) 
					printf("pushStack\npushQueue\n2 popStack popQueue\n");
				else if(a.size()==1)
					printf("pushStack\n1 popStack\n");
				else printf("0\n");
				a.clear();
				continue;
			}
			int m1=a.size(),m2=m1,m3=m1;
			a.push_back(-1);
			for(i=0;i<a.size();i++)
				if(a[m1]<a[i]) m1=i;
			for(i=0;i<a.size();i++)
				if(m1!=i and a[m2]<a[i]) m2=i;
			for(i=0;i<a.size();i++)
				if(m1!=i and m2!=i and a[m3]<a[i]) m3=i;

			for(i=0;i<a.size()-1;i++)
			{
				if(i==m1) printf("pushStack\n");
				else if(i==m2) printf("pushQueue\n");
				else if(i==m3) printf("pushFront\n");
				else printf("pushBack\n");
			}
			printf("3 popStack popQueue popFront\n");
			a.clear();
		}
		else a.push_back(p);
	}
	for(i=0;i<a.size();i++) printf("pushStack\n");
	return 0;
}