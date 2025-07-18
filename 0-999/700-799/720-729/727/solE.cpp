#include <cassert>

#include <cstdint>

#include <cstdio>

#include <algorithm>

#include <vector>

using namespace std;

typedef int64_t i64;

typedef pair<i64, int> PLI;



static const i64 MOD = 92233720368547753LL;

static const int E = 29;



static const auto NMAX = 100000;

static const auto KMAX = 100000;

static const auto MMAX = 100000;



static int N;

static int K;

static int M;

static int NK;

static i64 EK_;

static char S[1000000 + 1];

static i64 H[1000000];

static PLI G[MMAX];



static inline int find(i64 x) {

	auto i = lower_bound(G, G + M, PLI(x, 0)) - G;

	return i != M && G[i].first == x ? G[i].second : 0;

}



static inline i64 SUB(i64 a, i64 b) {

	a -= b;

	return a < 0 ? a + MOD : a;

}



static void solve() {

	static int Z[NMAX];

	vector<int> A;

	for(auto o = 0; o < K; o++) {

		A.clear();

		auto s = o;

		for(auto i = 0; i < N; i++) {

			auto j = find(H[s]);

			if(j == 0) break;

			if(Z[j - 1] == o + 1) break;

			Z[j - 1] = o + 1;

			A.push_back(j);

			if((s += K) >= NK) s -= NK;

		}

		if(A.size() != N) continue;

		printf("YES\n%d", A[0]);

		for(auto i = 1; i < N; i++) printf(" %d", A[i]);

		putchar('\n');

		return;

	}

	puts("NO");

}



static void build() {

	NK = N * K;

	EK_ = 1;

	for(auto i = 1; i < K; i++) EK_ = EK_ * E % MOD;

	for(auto i = 0; i < K; i++) H[0] = (H[0] * E + (S[i] - 'a' + 1)) % MOD;

	for(auto i = 1; i < NK; i++) {

		H[i] = (SUB(H[i - 1], (S[i - 1] - 'a' + 1) * EK_ % MOD) * E + (S[(i + K - 1) % NK] - 'a' + 1)) % MOD;

	}

}



int main() {

	scanf("%d%d%s%d", &N, &K, S, &M);

	build();

	for(auto i = 0; i < M; i++) {

		scanf("%s", S);

		i64 h = 0;

		for(auto j = 0; j < K; j++) h = (h * E + (S[j] - 'a' + 1)) % MOD;

		G[i] = PLI(h, i + 1);

	}

	sort(G, G + M);

	solve();

	return 0;

}