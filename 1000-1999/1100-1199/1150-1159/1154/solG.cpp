/*
Created at: 2019 April 16 19:24:04
Created by: metamorphosis00
*/
#include <bits/stdc++.h>
#define int long long

using namespace std;

const int N = 2e6;

int readInt()
{
    bool minus = false;
    int result = 0;
    char ch;
    ch = getchar();
    while (true) {
        if (ch == '-')
            break;
        if (ch >= '0' && ch <= '9')
            break;
        ch = getchar();
    }
    if (ch == '-')
        minus = true;
    else
        result = ch - '0';
    while (true) {
        ch = getchar();
        if (ch < '0' || ch > '9')
            break;
        result = result * 10 + (ch - '0');
    }
    if (minus)
        return -result;
    else
        return result;
}
int a[N];
int cnt[10000001];
int ans = LLONG_MAX;
int lid, rid;
void check(int res, int id1, int id2)
{
    if (res < ans) {
        ans = res;
        lid = id1;
        rid = id2;
    }
    return;
}
main()
{
    int n = readInt();
    for (int i = 1; i <= n; i++) {
        a[i] = readInt();
        if (cnt[a[i]]) {
            check(a[i], cnt[a[i]], i);
        }
        cnt[a[i]] = i;
    }
    int t = 1e7;
    for (int i = 1; i <= t; i++) {
        int l = 0, r = 0, lid1, rid1;
        for (int j = i; j <= t; j += i) {
            if (cnt[j]) {
                if (!l) {
                    l = j;
                    lid1 = cnt[j];
                }
                else {
                    r = j;
                    rid1 = cnt[j];
                }
            }
            if (r)
                break;
        }
        if (r) {
            check(1ll * l * r / i, lid1, rid1);
        }
    }
    if(lid > rid)
        swap(lid, rid);
    cout << lid << " " << rid;
}