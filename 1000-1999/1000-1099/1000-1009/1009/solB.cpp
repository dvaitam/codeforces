#include<bits/stdc++.h>
#define INF	0x3f3f3f3f
#define	LL	long long
#define	MAXN	100010
using namespace std;
char s[MAXN], t[MAXN];
int size[3], last;

template <typename T> void chkmin(T &x, T y){x = min(x, y);}
template <typename T> void chkmax(T &x, T y){x = max(x, y);}
template <typename T> void read(T &x){
	x = 0; int f = 1; char ch = getchar();
	while (!isdigit(ch)) {if (ch == '-') f = -1; ch = getchar();}
	while (isdigit(ch)) {x = x * 10 + ch - '0'; ch = getchar();}
	x *= f;
}

int main(){
	scanf("%s", s + 1);
	int n = strlen(s + 1);
	last = n;
	for (int i = n; i >= 1; --i){
		if (s[i] == '0') ++size[0];
		else if (s[i] == '1') ++size[1];
		else {
			while (size[0] > 0){
				t[last] = '0', --size[0];
				--last;
			}
			t[last] = '2';
			--last;
		}
	}
	while (size[0] + size[1]){
		if (size[1] > 0) t[last] = '1', --size[1];
		else if (size[0] > 0) t[last] = '0', --size[0];
		--last;
	}
	printf("%s", t + 1);
	return 0;
}