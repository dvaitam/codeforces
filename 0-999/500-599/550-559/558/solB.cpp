#include <iostream>
#include <vector>
#include <string>
#include <cstdio>
#include <memory.h>
#include <queue>
#include <map>
#include <algorithm>

using namespace std;

const int MAX_N = 100000;

pair<int, int> a[MAX_N];


int main()
{
//    freopen("infile.txt", "r", stdin);

    int n, m, t;
    scanf("%d", &n);

    for (int i = 0; i < n; ++i) {
        scanf("%d", &a[i].first);
        a[i].second = i;
    }

    sort(a, a+n);

    int maxbty = 1;
    int l = 0, r = 0;
    int retl = 0, retr = 0;

    while (r < n) {
        if (a[l].first == a[r].first) {
            ++r;
        }
        else {
            if (r-l > maxbty) {
                maxbty = r-l;
                retl = l;
                retr = r-1;
            }
            else if (r-l == maxbty && a[r-1].second-a[l].second < a[retr].second-a[retl].second) {
                retl = l;
                retr = r-1;
            }
            l = r;
        }
    }
    if (r-l > maxbty) {
        maxbty = r-l;
        retl = l;
        retr = r-1;
    }
    else if (r-l == maxbty && a[r-1].second-a[l].second < a[retr].second-a[retl].second) {
        retl = l;
        retr = r-1;
    }
    l = r;

    printf("%d %d\n", a[retl].second+1, a[retr].second+1);

    return 0;
}