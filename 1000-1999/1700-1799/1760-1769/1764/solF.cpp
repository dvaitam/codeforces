#include <bits/stdc++.h>



using namespace std;



const int INF = 1.01e9;



int main() {

	ios_base::sync_with_stdio(false); cin.tie(nullptr);



	int N; cin >> N;

	vector<vector<int64_t>> A(N, vector<int64_t>(N));



	for (int i = 0; i < N; ++i) {

		for (int j = 0; j < N; ++j) {

			if (i >= j) {

				cin >> A[i][j];

				A[j][i] = A[i][j];

			}

		}

	}	



	vector<int> sz(N, 1);

	vector<int> alive(N, true);

	for (int z = 0; z < N - 1; ++z) {

		int v = -1;

		for (int i = 0; i < N; ++i) if (alive[i]) {

			if (v == -1 || A[i][i] > A[v][v]) {

				v = i;

			}

		}



		alive[v] = false;



		int u = -1;

		for (int i = 0; i < N; ++i) if (alive[i]) {

			if (u == -1 || A[v][i] > A[v][u]) {

				u = i;

			}

		}



		int w = (A[u][u] - A[v][u]) / sz[v];

		sz[u] += sz[v];

		cout << v + 1 << ' ' << u + 1 << ' ' << w << '\n';

	}

	return 0;

}