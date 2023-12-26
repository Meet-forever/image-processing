#include "Image.h"

int main(int argc, char **argv)
{
    Image temp_img("test.png");

    // test_img.write("new.png");
    // Image temp_img = test_img;

    // temp_img = temp_img.rotate();
    // temp_img = temp_img.rotate(true);

    // temp_img = temp_img.rotateClockwise();
    // temp_img = temp_img.rotateCounterClockwise();

    // temp_img = temp_img.darkenImage(50);
    // temp_img = temp_img.lightenImage(50);

    // temp_img = temp_img.flipX();
    // temp_img = temp_img.flipY();

    // temp_img = temp_img.grayscale_avg();
    // temp_img = temp_img.grayscale_lum();

    temp_img.write("temp_img.png");

    // Image blank(100, 100, 3);
    // blank.write("blank.jpg");
    return 0;
}