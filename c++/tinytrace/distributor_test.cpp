#include <gtest/gtest.h>

#include "distributor.h"

TEST(DistributorTest, MatchesLoop) {
  int nx = 200;
  int ny = 100;

  distributor *d = new distributor(nx, ny);

  int expected_count = 0;
  int c, x, y;
  for (int j = ny - 1; j >= 0; j--) {
    for (int i = 0; i < nx; i++) {
      bool res = d->next_pixel(c, x, y);

      EXPECT_EQ(res, true);
      EXPECT_EQ(c, expected_count);
      EXPECT_EQ(x, i);
      EXPECT_EQ(y, j);

      expected_count++;
    }
  }
}

TEST(DistributorTest, RandomizeMatchesCount) {
  int nx = 200;
  int ny = 100;

  auto counts = new std::tuple<int, int>[nx * ny + 1];

  int c = 0;
  for (int j = ny - 1; j >= 0; j--) {
    for (int i = 0; i < nx; i++) {
      counts[c] = std::make_tuple(i, j);
      c++;
    }
  }

  distributor *d = new distributor(nx, ny);
  d->set_randomize(true);

  int x, y;
  for (int i = 0; i < nx * ny; i++) {
    bool res = d->next_pixel(c, x, y);

    auto px = counts[c];

    // std::cout << c << " " << x << " " << y << " -- " << std::get<0>(px) << "
    // "
    //           << std::get<1>(px) << "\n";

    EXPECT_EQ(res, true);
    EXPECT_EQ(std::get<0>(px), x);
    EXPECT_EQ(std::get<1>(px), y);
  }
}
