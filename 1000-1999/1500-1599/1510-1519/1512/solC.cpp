#include <cstdio>
#include <cstring>
using namespace std;
int T, A, B, len;
char str[200010];
int OK()
{
	for(int i=0;i<len;i++)
	{
		int x = len-i-1;
		if(str[i]!=str[x])
		    return 0;
	}
	return 1;
}
int main()
{
	scanf("%d", &T);
	while(T--)
	{
		scanf("%d%d", &A, &B);
		scanf("%s", str);
		len = strlen(str);
		int wenhao = 0;
		for(int i=0;i<len;i++)
		{
			if(str[i]=='?')
			    wenhao++;
			if(str[i]=='1')
			    B--;
			if(str[i]=='0')
			    A--;
		}
		for(int i=0;i<len;i++)
		{
			int x = len-1-i;
			if(str[i]!='?'&&str[x]=='?')
			{
				if(str[i]=='0')
				{
					str[x] = '0';
					A--;
				}
				else
				{
					str[x] = '1';
					B--;
				}
				wenhao--;
			}
		}
		for(int i=len/2;i<len;i++)
		{
			if(str[i]!='?')
			    continue;
			int x = len-i-1;
			if(i==x)
			{
				if(A%2==1)
				{
					str[i] = '0';
					A--;
				}
				else if(B%2==1)
				{
					str[i] = '1';
					B--;
				}
				wenhao--;
				continue;
			}
			if(A)
			{
				A -= 2;
				str[i] = '0';
				str[x] = '0';
				wenhao -= 2;
			}
			else if(B)
			{
				B -= 2;
				str[i] = '1';
				str[x] = '1';
				wenhao -= 2;
			}
		}
		
		if(!OK()||(A!=0||B!=0||wenhao!=0))
		    printf("-1\n");
		else
		{
			printf("%s\n", str);
		}
	}
	return 0;
}