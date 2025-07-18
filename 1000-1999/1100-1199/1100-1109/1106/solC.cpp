#include <bits/stdc++.h>
#define jh(x,y) x^=y^=x^=y
#define lowbit(x) (x&-x)
#define loc(x,y) (x-1)*n+y
#define rg register
#define inl inline
typedef long long ll;
const int N = 3e5 + 5;
using namespace std;
namespace fast_IO {
	inl ll read()
	{
		rg ll num = 0;
		rg char ch;
		rg bool flag = false;
		while ((ch = getchar()) == ' ' || ch == '\n' || ch == '\r');
		if (ch == EOF)return ch; if (ch == '-')flag = true; else num = ch ^ 48;
		while ((ch = getchar()) != ' '&&ch != '\n'&&ch != '\r'&&~ch)
			num = (num << 1) + (num << 3) + (ch ^ 48);
		if (flag)return -num; return num;
	}
	inl ll max(rg ll a, rg ll b) { if (a > b)return a; return b; }
	inl ll min(rg ll a, rg ll b) { if (a < b)return a; return b; }
	void write(rg long long x) { if (x < 0)putchar('-'), x = -x; if (x >= 10)write(x / 10); putchar(x % 10 ^ 48); }
};
int a[N];

int main(void)
{
	rg ll ans = 0;
	rg int n = fast_IO::read();
	for (rg int i = 1; i <= n; ++i)a[i] = fast_IO::read();
	sort(a + 1, a + n + 1);  
	rg int l = 1, r = n;
	for (rg int i = 1; i <= n / 2; ++i)
	{
		if (r - l + 1 != 3)
		{
			ans += (a[l] + a[r])*(a[l] + a[r]);
			++l, --r;
		}
		else
		{
			ans += (a[l] + a[l + 1] + a[l + 2])*(a[l] + a[l + 1] + a[l + 2]);
			break;
		}
	}
	fast_IO::write(ans);
	return 0;
}