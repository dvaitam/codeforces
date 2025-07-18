#include <iostream>

#include <cmath>

#include <cstdio>

#include <cstring>

#include <queue>

#include <stack>

#include <algorithm>



using namespace std;



int a[110];

bool vis[110];



int main()

{

    int n;

    while (cin >> n){

        int sum = 0;

        for (int i = 1; i <= n; i++){

            cin >> a[i];

            sum += a[i];

        }

        memset(vis, false, sizeof(vis));

        int d = sum / (n / 2);

//        cout << d << endl;

        for (int i = 1; i <= n; i++){

            for (int j = i + 1; j <= n; j++){

                if (a[i] + a[j] == d && !vis[i] && !vis[j]){

                    vis[i] = true;

                    vis[j] = true;

                    cout << i << ' ' << j << endl;

                    break;

                }

//                cout << a[i] << ' ' << a[j] << endl;

            }

        }

    }

    return 0;

}