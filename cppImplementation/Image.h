#include <stdint.h>
#include <cstdio>

enum ImageType{
    PNG, JPG, BMP, TGA
};

struct Image
{
    uint8_t* data = NULL;
    size_t size = 0;
    int w;
    int h;
    int channels;
    
    Image(const char* filename);
    Image(int w, int h, int channels);
    Image(const Image& img);
    ~Image();

    bool read(const char* filename);
    bool write(const char* filename);

    ImageType getFileType(const char* filename);

    // Effects
    Image &darkenImage(const int percentage);
    Image &lightenImage(const int percentage);
    Image &grayscale_avg();
    Image &grayscale_lum();

    // Transformations
    Image &flipX();
    Image &flipY();
    Image &rotate(bool opposite = false);
    Image &rotateClockwise();
    Image &rotateCounterClockwise();
};