//In the name of God
#include <iostream>
#include <vector>
using namespace std;

const int N = 1e5 + 2;

int n, m;
vector<int> adj[N], push;
int pt[N], a[N];

void dfs(int v) {
	pt[v]++, push.push_back(v);
	for (int i = 0; i < adj[v].size(); i++)
		pt[adj[v][i]]++;
	for (int i = 0; i < adj[v].size(); i++)
		if (pt[adj[v][i]] == a[adj[v][i]])
			dfs(adj[v][i]);
}
int main() {
	ios::sync_with_stdio(false);
	cin >> n >> m;
	for (int i = 0; i < m; i++) {
		int u, v;
		cin >> v >> u;
		v--, u--;
		adj[v].push_back(u);
		adj[u].push_back(v);
	}
	for (int i = 0; i < n; i++)
		cin >> a[i];
	for (int i = 0; i < n; i++)
		if (pt[i] == a[i])
			dfs(i);
	cout << push.size() << '\n';
	for (int i = 0; i < push.size(); i++)
		cout << push[i] + 1 << ' ';
	cout << '\n';
	return 0;
}