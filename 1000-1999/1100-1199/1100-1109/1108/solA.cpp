#include<bits/stdc++.h>
using namespace std;
int main()
{
	int t1,t2;
    int a,b,l1,l2,r1,r2;
	scanf("%d",&t1);
	for(t2=0;t2<t1;t2++)
	{
	    scanf("%d %d %d %d",&l1,&r1,&l2,&r2);
	    if(l1==r1)
	    {
	        a=l1;
	        b=l2;
	        if(a==b)
	        b++;
	        printf("%d %d\n",a,b);
	    }
	    if(l2==r2)
	    {
	        a=l1;
	        b=l2;
	        if(a==b)
	        a++;
	        printf("%d %d\n",a,b);
	    }
	    else
	    {
	        a=l1;
	        b=l2;
	        if(a==b)
	        b++;
	        printf("%d %d\n",a,b);
	    }
	    
	}
	return 0;
}