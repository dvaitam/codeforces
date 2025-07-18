#include <cstdio>

int a[100005];

int main() {
    int n, a0 = 1;
    scanf("%d", &n);
    for (int i = 0; i < n; ++i) {
        scanf("%d", a + i);
        if (a[i] >= 0) a[i] = ~a[i];
        if (~a[i]) a0 = 0;
    }
    if (a0) {
        if (n & 1)
            for (int i = 0; i < n; ++i) a[i] = 0;
    } else {
        if (n & 1) {
            int mv = a[0], mi = 0;
            for (int i = 1; i < n; ++i) if (a[i] < mv) {
                mv = a[i];
                mi = i;
            }
            a[mi] = ~a[mi];
        }
    }
    printf("%d", a[0]);
    for (int i = 1; i < n; ++i) printf(" %d", a[i]);
    return 0;
}