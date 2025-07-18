#include <bits/stdc++.h>
using namespace std;

int main() {
	cin.tie(0)->sync_with_stdio(0);

	int n;
	cin >> n;

	vector<int> a(n);
	for (auto &ai : a) cin >> ai, ai--;

	vector<int> ideg(n);
	for (auto ai : a) ideg[ai] += 1;

	vector<int> type(n);
	queue<pair<int, int>> q;
	auto Push = [&](int i, int t) -> void {
		if (type[i] != 0) {
			if (type[i] != t) {
				cout << -1;
				exit(0);
			}
			return;
		}
		type[i] = t;
		q.push({i, t});
	};
	for (int i = 0; i < n; i++) {
		if (ideg[i] == 0) {
			Push(i, 1);
		}
	}
	while (q.size() > 0) {
		auto [i, t] = q.front();
		q.pop();

		if (t == 1) {
			Push(a[i], 2);
		}
		if (t == 2) {
			ideg[a[i]] -= 1;
			if (ideg[a[i]] == 0) {
				Push(a[i], 1);
			}
		}
	}

	for (int i = 0; i < n; i++) {
		if (type[i] != 0) continue;

		auto Dfs = [&](auto Down, int j, int t) -> void {
			if (type[j] != 0) {
				if (type[j] != t) {
					cout << -1;
					exit(0);
				}
				return;
			}
			type[j] = t;

			Down(Down, a[j], 3 - t);
		};
		Dfs(Dfs, i, 1);
	}

	vector<int> ans;
	for (int i = 0; i < n; i++) {
		if (type[i] == 1) {
			ans.push_back(a[i]);
		}
	}
	
	cout << ans.size() << '\n';
	for (auto i : ans) cout << i + 1 << ' ';
}
/*
1. Think big picture!

2. Don't waste time coding before thinking it through!

3. Think fast! Don't be afraid to churn through ideas!

4. Try modelling the problem with as few moving parts as possible!
	- What does the problem look like?
	- What might the solution look like?

5. Try solving an easier version of the problem!
	- What if the problem didn't have this constraint?
	- How would I solve a more general version of the problem?

6. Ask questions!
	- What has similar behavior?
	- I've assumed this, what now?

Uncircling an element implies that whatever it's pointing to is
circled

If nothing forces it, it must not be circled
*/