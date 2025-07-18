#include<bits/stdc++.h>

using namespace std;

typedef long long ll;

#define forn(i, n)  for(int i = 0;i<n;i++)

#define for1(i, n)  for(int i = 1;i<=n;i++)

#define MASK(i)  (1<<(i))

int l[5], r[5];

int main()

{

    //freopen("bike.inp","r",stdin);

    int n;

    scanf("%d",&n);

    forn(i, n) scanf("%d %d",&l[i], &r[i]);

    long double ans = 0;

    long double dem = 0;

    setprecision(30);



    for(int tap = 1;tap <MASK(n);tap ++)

    {

        int L = -1e9, R = 1e9;

        forn(i, n)

        {

            if(tap & MASK(i))

            {

                L = max(L, l[i]);

                R = min(R, r[i]);

            }

        }

        if(L > R) continue;



        ///case 1 : only one other > b and not in the set

        for(int b = L;b<=R;b++)

        {

            forn(win, n)

            {

                if(tap & MASK(win)) continue;



                if(r[win] <= b) continue;



                long double c = r[win] - max(l[win] - 1, b);///so cac cai lon hon han



                forn(j, n)

                {

                    if(j == win) continue;

                    if(tap & MASK(j)) continue;

                    if(l[j] >= b)

                    {

                        c = 0;

                        break;

                    }

                    c = c * (min(r[j] + 1, b) - l[j]);///do dai no day

                }

                ans = ans + c * b;

                dem += c;

            }

            //cout<<'c' << c << endl;



            long double c = 1;

            ///TH 2 : la

            if(__builtin_popcount(tap) ==1 ) continue;///neu chi co dung 1 phan tu thi thoi



            forn(other, n)

            {

                if(tap & MASK(other)) continue;

                if(l[other] >= b)

                {

                    c = 0;

                    break;

                }

                c = c * (min(r[other] + 1, b) - l[other]);

            }

            ans = ans + c * b;

            dem += c;

        }

    }

    cout<<setprecision(10)<< fixed << ans/(dem + 0.0000000) << endl;

}