// prask3.cpp : This file contains the 'main' function. Program execution begins and ends there.

//



#include <iostream>

#include <vector>

#include <set>

using namespace std;

int main()

{

    int t;

    cin >> t;

    while (t--)

    {

        int n;

        cin >> n;

        if (n == 1)

        {

            cout << 1 << endl;

            cout << "1 2" << endl;

        }

        if (n == 2)

        {

            cout << 1 << endl;

            cout << "2 6" << endl;

        }

        else if (n > 2)

        {

            cout << n / 2 + n % 2 << endl;

            for (int i = 1, j = 3 * n; i <= j; i += 3, j -= 3)

            {

                cout << i << ' ' << j << endl;

            }

        }

    }

    return 0;

    ///std::cout << "Hello World!\n";

}



// Run program: Ctrl + F5 or Debug > Start Without Debugging menu

// Debug program: F5 or Debug > Start Debugging menu



// Tips for Getting Started: 

//   1. Use the Solution Explorer window to add/manage files

//   2. Use the Team Explorer window to connect to source control

//   3. Use the Output window to see build output and other messages

//   4. Use the Error List window to view errors

//   5. Go to Project > Add New Item to create new code files, or Project > Add Existing Item to add existing code files to the project

//   6. In the future, to open this project again, go to File > Open > Project and select the .sln file