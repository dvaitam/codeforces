#include <iostream>
using namespace std;

#define MAXN 128

struct event {
    int l, r, cur;
} n[MAXN];

int main() {
    int N, i, j;
    bool flag;

    cin >> N;
    for (i=0; i<N; i++) {
        cin >> n[i].l >> n[i].r; n[i].cur = n[i].l;
    }

    do {
        flag = false;
        for (i=0; i<N && !flag; i++) for (j=i+1; j<N && !flag; j++) if (n[i].cur == n[j].cur) {
            if (n[i].r < n[j].r) n[j].cur++;
            else n[i].cur++;
            flag = true;
        }
    } while (flag);
    for (i=0; i<N-1; i++) cout << n[i].cur << ' '; cout << n[N-1].cur << endl;

    return 0;
}