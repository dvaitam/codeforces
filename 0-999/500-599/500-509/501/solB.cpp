#include<stdio.h>
#include<string.h>
int main() 
{
	int t,i,j,m;char s1[1010][30],s2[1010][30];
	scanf("%d",&t);
	for(i=1;i<=t;i++)
	{
		scanf("%s%s",s1[i],s2[i]);
		for(j=i-1;j>=1;j--)
		if(strcmp(s1[i],s2[j])==0)
		{
    	strcpy(s2[j],s2[i]);
		i=i-1;
		t=t-1;
		break;
		}
	}
	printf("%d\n",i-1);
	m=i-1;
	for(i=1;i<=m;i++)
	printf("%s %s\n",s1[i],s2[i]);
	return 0;
}