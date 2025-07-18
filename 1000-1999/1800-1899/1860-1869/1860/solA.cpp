#include <bits/stdc++.h>
using namespace std;
int n, t;
char s[51];
int main()
{
	scanf("%d", &t);
	while (t--)
	{
		scanf("%s", s);
		n = 0;
		while (s[n])
			n++;
		bool type_1 = false;
		for (int i = 1; i < n; i++)
			if (s[i - 1] == ')' && s[i] == '(')
				type_1 = true;
		if (type_1)
		{
			puts("YES");
			for (int _ = n; _--; )
				putchar('(');
			for (int _ = n; _--; )
				putchar(')');
			puts("");
			continue;
		}
		bool type_2 = false;
		for (int i = 1; i < n; i++)
			if (s[i - 1] == s[i])
				type_2 = true;
		if (type_2)
		{
			puts("YES");
			for (int _ = n; _--; )
				putchar('('), putchar(')');
			puts("");
			continue;
		}
		puts("NO");
	}
}