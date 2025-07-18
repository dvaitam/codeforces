#include <iostream>

#include <algorithm>

using namespace std;



int vis[150]; char edges[150]; string ans;



void dfs(char u) {

	vis[u] = 3;

	if (vis[edges[u]] != 3 && edges[u] >= 'a') dfs(edges[u]);

	ans.push_back(u);

}



int main() {

	int n, i;

	string s;

	cin >> n;



	while (n--) {

		cin >> s;

		for (i = 0; i + 1 < s.size(); ++i) edges[s[i]] = s[i + 1], vis[s[i + 1]] = 2;

		if (vis[s[0]] != 2) vis[s[0]] = 1;

	}



	for (i = 'a'; i <= 'z'; ++i) if (vis[i] == 1) dfs(i);

	reverse(ans.begin(), ans.end());

	cout << ans << endl;

	return 0;

}