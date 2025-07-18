#include <stdio.h>

#include <stdlib.h>

#include <string.h>

#include <sys/time.h>



#define N	300000

#define M	300000



int min(int a, int b) { return a < b ? a : b; }



unsigned int X;



void srand_() {

	struct timeval tv;



	gettimeofday(&tv, NULL);

	X = tv.tv_sec ^ tv.tv_usec | 1;

}



int rand_() {

	return (X *= 3) >> 1;

}



int ii[M], jj[M], ds[M]; char tree[M];

int *eh[N], eo[N], pp[N], ta[N], tb[N], prev[N], next[N], xx[N], _; char in[N];



void append(int i, int h) {

	int o = eo[i]++;



	if (o >= 2 && (o & o - 1) == 0)

		eh[i] = (int *) realloc(eh[i], o * 2 * sizeof *eh[i]);

	eh[i][o] = h;

}



int find(int i) {

	return ds[i] < 0 ? i : (ds[i] = find(ds[i]));

}



void join(int i, int j) {

	i = find(i);

	j = find(j);

	if (i == j)

		return;

	if (ds[i] > ds[j])

		ds[i] = j;

	else {

		if (ds[i] == ds[j])

			ds[i]--;

		ds[j] = i;

	}

}



void dfs(int p, int f, int i) {

	int o;



	pp[i] = p, ta[i] = tb[i] = ++_;

	for (o = eo[i]; o--; ) {

		int h = eh[i][o], j = i ^ ii[h] ^ jj[h];



		if (j != p) {

			if (!ta[j]) {

				tree[h] = 1;

				dfs(i, h, j);

				tb[i] = min(tb[i], tb[j]);

				if (tb[j] < ta[i])

					join(h, f);

			} else if (ta[j] < ta[i]) {

				tb[i] = min(tb[i], ta[j]);

				join(h, f);

			}

		}

	}

}



int compare_comp(int h1, int h2) {

	int h1_ = find(h1);

	int h2_ = find(h2);



	return h1_ != h2_ ? h1_ - h2_ : ta[ii[h1]] - ta[ii[h2]];

}



int compare_lr(int h1, int h2) {

	return xx[ii[h1]] != xx[ii[h2]] ? xx[ii[h1]] - xx[ii[h2]] : xx[jj[h2]] - xx[jj[h1]];

}



int (*compare)(int, int);



void sort(int *hh, int l, int r) {

	while (l < r) {

		int i = l, j = l, k = r, h = hh[l + rand_() % (r - l)], tmp;



		while (j < k) {

			int c = compare(hh[j], h);



			if (c == 0)

				j++;

			else if (c < 0) {

				tmp = hh[i], hh[i] = hh[j], hh[j] = tmp;

				i++, j++;

			} else {

				k--;

				tmp = hh[j], hh[j] = hh[k], hh[k] = tmp;

			}

		}

		sort(hh, l, i);

		l = k;

	}

}



int solve(int *hh, int m) {

	static int qu[M];

	int h, h_, i, j, k, cnt, tmp;



	h_ = hh[0], i = ii[h_], j = jj[h_];

	if (m == 1) {

		eh[i][eo[i]++] = h_, eh[j][eo[j]++] = h_;

		return 1;

	}

	cnt = 0;

	in[i] = 1, qu[cnt++] = i;

	next[j] = i, prev[i] = j;

	while (j != i) {

		in[j] = 1, qu[cnt++] = j;

		next[pp[j]] = j, prev[j] = pp[j];

		j = pp[j];

	}

	for (h = 1; h < m; h++) {

		h_ = hh[h], i = ii[h_], j = jj[h_];

		if (tree[h_] || in[j])

			continue;

		while (!in[j])

			j = pp[j];

		if (j == prev[i]) {

			j = jj[h_], k = i;

			while (!in[j]) {

				in[j] = 1, qu[cnt++] = j;

				next[prev[k]] = j, prev[j] = prev[k];

				next[j] = k, prev[k] = j;

				k = j, j = pp[j];

			}

		} else if (j == next[i]) {

			j = jj[h_], k = i;

			while (!in[j]) {

				in[j] = 1, qu[cnt++] = j;

				next[j] = next[k], prev[next[k]] = j;

				next[k] = j, prev[j] = k;

				k = j, j = pp[j];

			}

		} else {

			while (cnt--)

				in[qu[cnt]] = 0;

			return 0;

		}

	}

	j = i = ii[hh[0]];

	cnt = 0;

	do

		qu[cnt++] = j, j = next[j];

	while (j != i);

	for (h = 0; h < cnt; h++)

		in[qu[h]] = 0, xx[qu[h]] = h;

	for (h = 0; h < m; h++) {

		h_ = hh[h];

		if (xx[ii[h_]] > xx[jj[h_]])

			tmp = ii[h_], ii[h_] = jj[h_], jj[h_] = tmp;

	}

	compare = compare_lr, sort(hh, 0, m);

	cnt = 0;

	for (h = 0; h < m; h++) {

		h_ = hh[h];

		while (cnt && xx[jj[qu[cnt - 1]]] <= xx[ii[h_]])

			cnt--;

		if (cnt && xx[jj[h_]] > xx[jj[qu[cnt - 1]]])

			return 0;

		qu[cnt++] = h_;

	}

	for (h = m - 1; h >= 0; h--) {

		h_ = hh[h], i = ii[h_];

		eh[i][eo[i]++] = h_;

	}

	for (h = 0; h < m; h++) {

		h_ = hh[h], j = jj[h_];

		eh[j][eo[j]++] = h_;

	}

	return 1;

}



int main() {

	int t;



	srand_();

	scanf("%d", &t);

	while (t--) {

		static int hh[M];

		int n, m, h, h_, h1, i, j, o, tmp, yes;



		scanf("%d%d", &n, &m);

		for (i = 0; i < n; i++)

			eh[i] = (int *) malloc(2 * sizeof *eh[i]), eo[i] = 0;

		for (h = 0; h < m; h++) {

			scanf("%d%d", &i, &j);

			ii[h] = i, jj[h] = j;

			append(i, h), append(j, h);

		}

		memset(ta, 0, n * sizeof *ta), memset(tb, 0, n * sizeof *tb);

		memset(ds, -1, m * sizeof *ds), memset(tree, 0, m * sizeof *tree);

		_ = 0, dfs(-1, -1, 0);

		for (h = 0; h < m; h++) {

			if (ta[ii[h]] > ta[jj[h]])

				tmp = ii[h], ii[h] = jj[h], jj[h] = tmp;

			hh[h] = h;

		}

		compare = compare_comp, sort(hh, 0, m);

		memset(eo, 0, n * sizeof *eo);

		memset(in, 0, n * sizeof *in), memset(prev, -1, n * sizeof *prev), memset(next, -1, n * sizeof *next);

		yes = 1;

		for (h = 0; h < m; h = h_) {

			h1 = find(hh[h]), h_ = h + 1;

			while (h_ < m && find(hh[h_]) == h1)

				h_++;

			if (!solve(hh + h, h_ - h)) {

				yes = 0;

				break;

			}

		}

		printf(yes ? "YES\n" : "NO\n");

		if (yes) {

			for (i = 0; i < n; i++) {

				for (o = 0; o < eo[i]; o++) {

					h = eh[i][o];

					printf("%d ", i ^ ii[h] ^ jj[h]);

				}

				printf("\n");

			}

		}

		for (i = 0; i < n; i++)

			free(eh[i]);

	}

	return 0;

}