#include <bits/stdc++.h>
#include <stdio.h>

#include <vector>

#include <algorithm>

#include <math.h>



using namespace std;



int maxx(int fd, int sd)

{

  return (fd > sd) ? fd : sd;

}



int minn(int fd, int sd)

{

  return (fd < sd) ? fd : sd;

}



vector<long long> x;

vector<long long> y;



long long SS(int a, int b, int c)

{

  long long result = (x[a] - x[b]) * (y[a] + y[b]) + (x[b] - x[c]) * (y[b] + y[c]) + (x[c] - x[a]) * (y[c] + y[a]);

  if (result < 0)

    result = -result;

  return result;

}



int main()

{

  int n;

  long long S;

  int count = 0;

  scanf("%d %lld", &n, &S);

  x.resize(n);

  y.resize(n);

  for (int i = 0; i < n; ++i)

  {

    scanf("%lld %lld", &x[i], &y[i]);

  }

  int a = 0;

  int b = 1;

  int c = 2;

  while (true)

  {

    long long curSS = SS(a, b, c);

    int newa = a, newb = b, newc = c;

    for (int i = 0; i < n; ++i)

    {

      if (SS(a, b, i) > curSS)

      {

        newc = i;

        curSS = SS(a, b, newc);

      }

    }

    for (int i = 0; i < n; ++i)

    {

      if (SS(a, i, newc) > curSS)

      {

        newb = i;

        curSS = SS(a, newb, newc);

      }

    }

    for (int i = 0; i < n; ++i)

    {

      if (SS(i, newb, newc) > curSS)

      {

        newa = i;

        curSS = SS(newa, newb, newc);

      }

    }

    if (((newa == a) && (newb == b) && (newc == c)) || (count > 1000))

    {

      break;

    }

    a = newa;

    b = newb;

    c = newc;

    count++;

  }

  printf("%lld %lld\n", x[a] + x[b] - x[c], y[a] + y[b] - y[c]);

  printf("%lld %lld\n", x[a] + x[c] - x[b], y[a] + y[c] - y[b]);

  printf("%lld %lld", x[c] + x[b] - x[a], y[c] + y[b] - y[a]);

  return 0;

}