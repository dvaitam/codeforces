#include <bits/stdc++.h>
using namespace std;
const int MXN = 510;
int n, buck[MXN], s[MXN], calc[MXN];
char res[MXN];

int main()
{
    cin >> n;
    for (int i = 0; i < n; ++ i)
    {
        cin >> s[i];
        ++ buck[s[i]];
    }
    for (int i = 1; i < MXN; ++ i)
    {
        if (buck[i] <= 2){
            ++ calc[buck[i]];
        }
        else {
            ++ calc[3];
        }
    }
    if ((calc[1] & 1) == 1 && calc[3] == 0) {
        cout << "NO" << endl;
        return 0;
    }
    cout << "YES" << endl;
    int remains = (calc[1] & 1);
    int t = 0;
    for (int i = 0; i < n; ++ i) {
        if (buck[s[i]] == 1) {
            t += 1;
            if (t <= calc[1]/2) {
                res[i] = 'A';
            }
            else {
                res[i] = 'B';
            }
        }
        else if (buck[s[i]] == 2) {
            res[i] = 'A';
        }
        else if (buck[s[i]] >= 3) {
            if (remains == 1) {
                remains = 0;
                res[i] = 'A';
            }
            else {
                res[i] = 'B';
            }
        }
    }
    for (int i = 0; i < n; ++ i) {
        cout << res[i];
    }
    cout << endl;
    return 0;
}