#include <iostream>



using namespace std;



int k, x, n, m;



int ncount(int fs1, int ls1, int fs2, int ls2, int num1, int num2, int it = 2)

{

    if (it == k)

        return num2;

    if (num2 > x)

        return -1;

    return ncount(fs2, ls2, fs1, ls2, num2, (ls1 == 0 && fs2 == 2) ? (num1 + num2 + 1) : (num1 + num2), it + 1);

}



int main()

{

    ios_base::sync_with_stdio(false), cin.tie(nullptr);



    cin >> k >> x >> n >> m;



    for (int fs1 = 0; fs1 < 3; fs1++)

    {

        for (int ls1 = 0; ls1 < 3; ls1++)

        {

            for (int fs2 = 0; fs2 < 3; fs2++)

            {

                for (int ls2 = 0; ls2 < 3; ls2++)

                {



                    for (int num1 = 0; num1 <= n / 2; num1++)

                    {

                        for (int num2 = 0; num2 <= m / 2; num2++)

                        {

                            bool check = true;



                            if (n == 1 && fs1 != ls1) check = false;

                            if (m == 1 && fs2 != ls2) check = false;

                            if (n > 1 && n % 2 == 0 && num1 == n / 2 && (!(fs1 == 0) || !(ls1 == 2))) check = false;

                            if (m > 1 && m % 2 == 0 && num2 == m / 2 && (!(fs2 == 0) || !(ls2 == 2))) check = false;

                            if (n > 1 && n % 2 == 1 && num1 == n / 2 && (!(fs1 == 0) && !(ls1 == 2))) check = false;

                            if (m > 1 && m % 2 == 1 && num2 == m / 2 && (!(fs2 == 0) && !(ls2 == 2))) check = false;

                            if (n == 2 && fs1 == 0 && ls1 == 2 && num1 != 1) check = false;

                            if (m == 2 && fs2 == 0 && ls2 == 2 && num2 != 1) check = false;



                            if (check && ncount(fs1, ls1, fs2, ls2, num1, num2) == x)

                            {

                                string s;

                                if (fs1 == 0)

                                    s = "A";

                                else if (fs1 == 1)

                                    s = "X";

                                else

                                    s = "C";

                                for (int i = 1; i < n - 1; i++)

                                {

                                    if (num1)

                                    {

                                        if (s[i - 1] == 'A')

                                        {

                                            s += "C";

                                            num1--;

                                        }

                                        else

                                            s += "A";

                                    }

                                    else

                                        s += "X";

                                }

                                if (n > 1)

                                {

                                    if (ls1 == 0)

                                        s += "A";

                                    else if (ls1 == 1)

                                        s += "X";

                                    else

                                        s += "C";

                                }

                                cout << s << "\n";



                                s = "";

                                if (fs2 == 0)

                                    s = "A";

                                else if (fs2 == 1)

                                    s = "X";

                                else

                                    s = "C";

                                for (int i = 1; i < m - 1; i++)

                                {

                                    if (num2)

                                    {

                                        if (s[i - 1] == 'A')

                                        {

                                            s += "C";

                                            num2--;

                                        }

                                        else

                                            s += "A";

                                    }

                                    else

                                        s += "X";

                                }

                                if (m > 1)

                                {

                                    if (ls2 == 0)

                                        s += "A";

                                    else if (ls2 == 1)

                                        s += "X";

                                    else

                                        s += "C";

                                }

                                cout << s << "\n";

                                return 0;

                            }

                        }

                    }



                }

            }

        }

    }



    cout << "Happy new year!";

    return 0;

}