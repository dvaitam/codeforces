#include<cstdio>
#include<queue>
#include<algorithm>
#include<string>

using namespace std;

void print(const string & S)
{
	bool isfirst = 1;
	int l = S.size();
	for(int i=0;i<l;++i)
	{
		if(S[i]=='#')
		{
			putchar('}');
			isfirst=1;
		}
		else
		{
		if(!isfirst)
			putchar(',');
		else
		{
			if(i!=0)
				putchar(',');
			putchar('{');
		}
		isfirst = 0;
		if(S[i]<='9')
			putchar(S[i]);
		else
		{
			putchar('1');
			putchar('0');
		}
		}
	}
	printf("\n");
}

int main()
{
	int n;
	queue <string> Old, New;
	scanf("%d",&n);
	Old.push("1#");
	n = n + '0';
	string Sold, Snew;
	int pos;
	int start;
	for(char c = '2'; c<=n; ++c)
	{
		int dir = 1;
		while(!Old.empty())
		{
			Sold=Old.front();
			Old.pop();
			if(dir==1)
			{
				start = 0;
			}
			else
			{
				Snew = Sold;
				Snew += c;
				Snew += '#';
				New.push(Snew);
				start=Sold.size()-1;
			}
			pos = start;
			while(pos >=0 && pos < Sold.size())
			{
				if(Sold[pos]=='#')
				{
					Snew = Sold;
					Snew.insert(pos,1, c);
					New.push(Snew);
				}
				pos += dir;
			}
			if(dir == 1)
			{
				Snew = Sold;
				Snew += c;
				Snew += '#';
				New.push(Snew);
			}
			dir = -dir;
		}
		while(!New.empty())
		{
			Old.push(New.front());
			New.pop();
		}
	}
	printf("%d\n",Old.size());
	while(!Old.empty())
	{
		print(Old.front());
		Old.pop();
	}
	return 0;
}