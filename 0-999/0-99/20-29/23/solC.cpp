#include <iostream>

#include <cstdlib>

#include <cstdio>

using namespace std;

int a[200000], s[200000], pos[200000];

void qs(int q, int w) {

  int e = q, r = w, t = a[q + (rand() % (w - q + 1))];

  do {

    while (a[e] < t) e++;

    while (a[r] > t) r--;

    if (e <= r) {

      swap(a[e], a[r]);

      swap(s[e], s[r]);

      swap(pos[e], pos[r]);

      e++; r--;

    }

  } while (e <= r);

  if (q < r) qs(q, r);

  if (e < w) qs(e, w);

}

int main() {

  #ifdef Vlad_kv

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

  #endif // Vlad_kv

  int q, q1, w, numTest, test;

  long long sum[2];

  scanf("%d", &numTest);

  for (test = 0; test < numTest; test++) {

    scanf("%d", &q);

    q1 = q * 2 - 1;

    for (w = 0; w < q1; w++) {

      scanf("%d%d", &a[w], &s[w]);

      pos[w] = w + 1;

    }

    qs(0, q1 - 1);

    sum[0] = sum[1] = 0;

    for (w = 0; w < q1; w++) {

      sum[w & 1] += s[w];

    }

    printf("YES\n");

    if (sum[0] >= sum[1]) {

      for (w = 0; w < q1; w += 2) {

        printf("%d ", pos[w]);

      }

    } else {

      for (w = 1; w < q1; w += 2) {

        printf("%d ", pos[w]);

      }

      printf("%d ", pos[q1 - 1]);

    }

    printf("\n");

  }

  return 0;

}