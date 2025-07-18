#include<iostream>
#include<stdio.h>
#include<string.h>
#include<stdlib.h>
#include<algorithm>
using namespace std;
const int N=100008;
char a[N];
int len;
int main(void)
{
	int i,j,k,p1,p2,p,t1,t2,t12,t22,s1,s2,q;
	scanf("%s",a+1);len=strlen(a+1);
	p=p1=p2=0;
	for(i=1;i<=len;i++){p=(p+a[i]-48)%3;if((a[i]-48)%3==1)p1++;else if((a[i]-48)%3==2)p2++;}
	if(p==0){printf("%s",a+1);return 0;}
	if(p==1&&(p1||p2>=2))
	{
		if(len==1&&p1==1){printf("-1");return 0;}
		if(len==2&&p2==2){printf("-1");return 0;}
		s1=s2=N;
		if(p1)
		{
			for(t1=len;t1>=1&&(a[t1]-48)%3!=1;t1--){;}j=t1;
			if(t1>1){for(q=1;q<t1;q++)printf("%c",a[q]);for(q=t1+1;q<=len;q++)printf("%c",a[q]);return 0;}
			else{for(j=2;a[j]=='0'&&j<len;j++);j--;s1=j+1;}
		}
		if(p2>=2)
		{
			for(t22=len;t22>=1&&(a[t22]-48)%3!=2;t22--);
			for(t2=t22-1;t2>=1&&(a[t2]-48)%3!=2;t2--);
			s2=2;k=t2;
			if(t2==1){for(k=2;(k<len-((a[len]-48)%3==2))&&(a[k]=='0'||k==t22);k++);k--;s2=k-t2+1+(t22>k);}
		}
		if(s1<=s2){printf("%s",a+j+1);return 0;}
		else
		{
			for(q=1;q<t2;q++)printf("%c",a[q]);
			if(k>=t22)for(q=k+1;q<=len;q++)printf("%c",a[q]);
			else
			{
				for(q=k+1;q<t22;q++)printf("%c",a[q]);
				for(q=t22+1;q<=len;q++)printf("%c",a[q]);
			}
		}
	}
	else if(p==2&&(p1>=2||p2))
	{
		if(len==1&&p2==1){printf("-1");return 0;}
		if(len==2&&p1==2){printf("-1");return 0;}
		s1=s2=N;
		if(p2)
		{
			for(t2=len;t2>=1&&(a[t2]-48)%3!=2;t2--){;}j=t2;
			if(t2>1){for(q=1;q<t2;q++)printf("%c",a[q]);for(q=t2+1;q<=len;q++)printf("%c",a[q]);return 0;}
			else{for(j=2;a[j]=='0'&&j<len;j++);j--;s2=j+1;}
		}
		if(p1>=2)
		{
			for(t12=len;t12>=1&&(a[t12]-48)%3!=1;t12--);
			for(t1=t12-1;t1>=1&&(a[t1]-48)%3!=1;t1--);
			s1=2;k=t1;
			if(t1==1){for(k=2;(k<len-((a[len]-48)%3==1))&&(a[k]=='0'||k==t12);k++);k--;s1=k-t1+1+(t12>k);}
		}
		if(s2<=s1){printf("%s",a+j+1);return 0;}
		else
		{
			for(q=1;q<t1;q++)printf("%c",a[q]);
			if(k>=t12)for(q=k+1;q<=len;q++)printf("%c",a[q]);
			else
			{
				for(q=k+1;q<t12;q++)printf("%c",a[q]);
				for(q=t12+1;q<=len;q++)printf("%c",a[q]);
			}
		}
	}
	else{printf("-1");return 0;}
	return 0;
}