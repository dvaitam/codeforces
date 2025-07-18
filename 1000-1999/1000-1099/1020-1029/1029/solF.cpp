#include <bits/stdc++.h>

#define llong long long

using namespace std;

vector<llong> f(llong x)
{
    vector<llong> m;
    for (llong i = 1; i * i <= x; i++)
    {
        if (x % i == 0)
        {
            m.push_back(i);
        }
    }
    return m;
}

int main(void)
{
    llong a = 0, b = 0, c = 0;
    cin >> a >> b;
    c = a + b;
    vector<llong> am = f(a);
    vector<llong> bm = f(b);
    vector<llong> cm = f(c);
    for (llong i = cm.size() - 1; i >= 0; i--)
    {
        for (llong j = 0; j < am.size(); j++)
        {
            if ((cm[i] >= am[j]) && ((c / cm[i]) >= (a / am[j])))
            {
                cout << (cm[i] + (c / cm[i])) * 2;
                return 0;
            }
        }
        for (llong j = 0; j < bm.size(); j++)
        {
            if ((cm[i] >= bm[j]) && ((c / cm[i]) >= (b / bm[j])))
            {
                cout << (cm[i] + (c / cm[i])) * 2;
                return 0;
            }
        }
    }
    return 0;
}