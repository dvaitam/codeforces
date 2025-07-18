//  Created by Sengxian on 2016/11/14.

//  Copyright (c) 2016年 Sengxian. All rights reserved.

#include <bits/stdc++.h>

using namespace std;



typedef long long ll;

const int MAX_N = 2000 + 3, MAX_M = 2000 + 3;

int n, m, a[MAX_N], cnt[MAX_M], id[MAX_M], tar[MAX_M];



inline bool cmp(const int &i, const int &j) {

	return cnt[i] > cnt[j];

}



void solve() {

	int per = n / m, more = n % m, ans = 0;

	for (int i = 0; i < n; ++i)

		if (a[i] < m) cnt[a[i]]++;

	for (int i = 0; i < m; ++i) id[i] = i;

	sort(id, id + m, cmp);



	for (int i = 0; i < more; ++i) if (cnt[id[i]] > per) tar[id[i]] = per + 1; else tar[id[i]] = per;

	for (int i = more; i < m; ++i) tar[id[i]] = per;



	for (int i = 0; i < m; ++i) if (cnt[i] > tar[i]) {

		int left = cnt[i] - tar[i];

		while (left--) {

			int idx = -1;

			for (int j = 0; j < n; ++j) if (a[j] == i) {idx = j; break;}

			for (int j = 0; j < m; ++j) if (cnt[j] < tar[j]) {

				a[idx] = j, cnt[j]++, cnt[i]--, ans++;

				break;

			}

		}

	}

	

	for (int i = 0; i < n; ++i) if (a[i] >= m)

		for (int j = 0; j < m; ++j) if (cnt[j] < tar[j]) {

			a[i] = j, cnt[j]++, ans++;

			break;

		}



	printf("%d %d\n", per, ans);

	for (int i = 0; i < n; ++i)

		printf("%d%c", a[i] + 1, i + 1 == n ? '\n' : ' ');

}



int main() {

#ifdef DEBUG

	freopen("C/test.in", "r", stdin);

#endif

	scanf("%d%d", &n, &m);

	for (int i = 0; i < n; ++i) scanf("%d", a + i), a[i]--;

	solve();

	return 0;

}