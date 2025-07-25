#include <cstdio>
#include <algorithm>
#include <vector>
using namespace std;
int main() {
    int n, m, a[102];
    vector<int> pos, neg;
    bool zero = false;

    scanf("%d", &n);
    for (int i=0; i<n; i++) {
        scanf("%d", &a[i]);
        if (a[i] < 0) neg.push_back(a[i]);
        else if (a[i] == 0) zero = true;
        else pos.push_back(a[i]);
    }
    sort(neg.begin(), neg.end());
    bool out = false;
    for (size_t i=0; i<pos.size(); i++) { printf("%d ", pos[i]); out = true; }
    for (size_t i=0; i+1<neg.size(); i+=2) { printf("%d %d ", neg[i], neg[i+1]); out = true; }
    if (!out) {
        if (zero) puts("0");
        else printf("%d\n", neg[0]);
    }

    return 0;
}