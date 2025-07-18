#include <bits/stdc++.h>
#include <bits/stdc++.h>

#define ll long long

using namespace std;

const int Mx = 1e5+5;

int A[10], B[Mx], ans[Mx], N, M, nc;

int razmer(char *ch)

{

	if(ch[0] == 'S')

            return 1;

	if(ch[0] == 'M')

            return 2;

	if(ch[0] == 'L')

            return 3;

	if(ch[1] == 'L')

            return 4;

	if(ch[2] == 'L')

            return 5;

    return 6;

}

char tab[6][5] = {"S", "M", "L", "XL", "XXL", "XXXL"};

struct P

{

	int x, id;

	bool operator < (const P &p) const

	{

		return x < p.x;

	}

}C[Mx];

int solve()

{

	if(M < N)

        return 0;

	sort(C, C + nc);

	for(int i = 1; i <= 6; ++i)

        if(A[i] < 0)

            return 0;

	for(int i = 0; i < nc; ++i)

	{

		int nx = C[i].x;

		if(A[nx])

		{

			A[nx]--; ans[C[i].id] = nx;

		}

        else

            if(A[nx+1])

            {

			A[nx+1]--; ans[C[i].id] = nx+1;

            }

            else

                return 0;

	}

	return 1;

}

int main()

{

	for(int i = 1; i <= 6; ++i)

        cin>>A[i], M += A[i];

	scanf("%d\n",&N);

	for(int i = 0; i < N; ++i)

    {

		char s[10]; fgets(s, sizeof(s), stdin);

		int sz1 = razmer(s), sz2 = 0;

		for(int j = 0; s[j]; ++j)

			if(s[j] == ',')

                sz2 = razmer(s+j+1);

		if(!sz2)

		{

			ans[i] = sz1; A[sz1]--;

		}

        else

        {

			C[nc].x = min(sz1, sz2);

			C[nc].id = i;

			nc++;

		}

	}

	if(solve())

	{

		cout<<"YES\n";

		for(int i = 0; i < N; ++i)

			cout<<tab[ans[i]-1]<<'\n';

	}

	else

		cout<<"NO";

}