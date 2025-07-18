#include <iostream>

#include <vector>

#include <algorithm>

using namespace std;



int main()

{

    int t;

    cin >> t;

    for(int tt = 0; tt < t; tt++)

    {

        int n, m;

        cin >> n >> m;

        if(n == 2 && m == 2)

        {

            cout << "WB\nBB\n";

            continue;

        }

        vector<vector<bool>> field(n);

        for(int i = 0; i < n; i++)

        {

            for(int j = 0; j < m; j++)

            {

                if(j % 2 == i % 2)

                {

                    field[i].push_back(true);

                }else

                {

                    field[i].push_back(false);

                }

            }

        }

        if(n % 2 == 0)

        {

            field[n - 1][m - 1] = true;

            field[n - 1][m - 2] = true;

        }

        if(m % 2 == 0)

        {

            field[n - 1][m - 1] = true;

            field[n - 2][m - 1] = true;

        }

        for(auto elem: field)

        {

            for(auto elem2: elem)

            {

                cout << (elem2 ? "B" : "W");

            }

            cout << "\n";

        }

    }

    return 0;

}