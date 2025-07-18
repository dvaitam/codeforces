#include<bits/stdc++.h>
using namespace std;
template <typename T>
inline void read(T &x) {
	x = 0;
	T f = 1;
	char c = getchar();
	for (; !isdigit(c); c = getchar()) if (c == '-') f = -1;
	for (; isdigit(c); c = getchar()) x = (x << 1) + (x << 3) + (c ^ 48);
	x *= f;
}
inline void d_read(double &x) {
	x = 0.0;
	int f = 1;
	char c = getchar();
	for (; !isdigit(c); c = getchar()) if (c == '-') f = -1;
	for (; isdigit(c); c = getchar()) x = x * 10 + (c ^ 48);
	if (c == '.'){
		double num = 1.0;
		c = getchar();
		for (; isdigit(c); c = getchar()) x = x + (num /= 10) * (c ^ 48);
	}
	x *= f;
}

template <typename T>
inline void w(T x) {
	if (x > 9) w(x / 10);
	putchar(x % 10 + 48);
}
template <typename T>
inline void write(T x, char c) {
	if (x < 0){
		putchar('-');
		x = -x;
	}
	w(x);
	putchar(c);
}
typedef pair <int, int> pii;
typedef long long ll;
typedef unsigned long long ull;
int cnt[200005], a[200005][2];
int siz[200005];
queue <int> q;
bool vis[200005], in[2000005];
vector <pii> v[200005];
stack <int> s;
int main(){
	int n, m;
	read(n); read(m);
	for (int i = 1; i <= n; i ++) {
		read(cnt[i]);
	}
	for (int i = 1; i <= m; i ++) {
		read(a[i][0]);
		read(a[i][1]);
		v[a[i][0]].push_back(make_pair(i, a[i][1]));
		v[a[i][1]].push_back(make_pair(i, a[i][0]));
	}
	for (int i = 1; i <= n; i ++) {
		siz[i] = v[i].size();
		if (cnt[i] >= siz[i]) {
			q.push(i);
			vis[i] = true;
		}
	}
	while(!q.empty()) {
		int d = q.front();
		q.pop();
		for (int i = 0; i < v[d].size(); i ++) {
			int u = v[d][i].first, val = v[d][i].second;
			if (!in[u]) {
				s.push(u);
				in[u] = true;
			}
			if (vis[val]) continue;
			siz[val] --;
			if (cnt[val] >= siz[val]) {
				q.push(val);
				vis[val] = true;
			}
		}
	}
	if (s.size() != m) {
		puts("DEAD");
	}
	else {
		puts("ALIVE");
		while(s.size() > 1) {
			write(s.top(), ' ');
			s.pop();
		}
		write(s.top(), '\n');
	}
	return 0;
}