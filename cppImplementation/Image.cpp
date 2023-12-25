#define STB_IMAGE_IMPLEMENTATION
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "Image.h"
#include "stb_image.h"
#include "stb_image_write.h"

Image::Image(const char *filename)
{
    if (!read(filename))
    {
        printf("Failed to read %s \n", filename);
    }
    size = w * h * channels;
    printf("Read %s \n", filename);
}

Image::Image(int w, int h, int channels) : w(w), h(h), channels(channels)
{
    size = w * h * channels;
    data = new uint8_t[size];
}

Image::Image(const Image &img) : Image(img.w, img.h, img.channels)
{
    memcpy(data, img.data, size);
}

Image::~Image()
{
    stbi_image_free(data);
}

bool Image::read(const char *filename)
{
    data = stbi_load(filename, &w, &h, &channels, 0);
    return data != NULL;
}

bool Image::write(const char *filename)
{
    ImageType type = getFileType(filename);
    int success;
    switch (type)
    {
    case PNG:
        success = stbi_write_png(filename, w, h, channels, data, w * channels);
        break;
    case JPG:
        success = stbi_write_jpg(filename, w, h, channels, data, 100);
        break;
    case BMP:
        success = stbi_write_bmp(filename, w, h, channels, data);
        break;
    case TGA:
        success = stbi_write_tga(filename, w, h, channels, data);
        break;
    }
    return success;
}

ImageType Image::getFileType(const char *filename)
{
    const char *ext = strrchr(filename, '.');
    if (ext != nullptr)
    {
        if (strcmp(ext, ".png") == 0)
        {
            return PNG;
        }
        else if (strcmp(ext, ".jpg") == 0)
        {
            return JPG;
        }
        else if (strcmp(ext, ".bmp") == 0)
        {
            return JPG;
        }
        else if (strcmp(ext, ".tga") == 0)
        {
            return JPG;
        }
    }
    return PNG;
}

// Effects
Image &Image::darkenImage(const int percentage)
{
    float shade = (100.0 - fmin(fmax((float)percentage, 0.0), 100.0)) / 100.0;
    for (int i = 0; i < size; ++i)
    {
        data[i] *= shade;
    }
    return *this;
}

Image &Image::lightenImage(const int percentage)
{
    float tint = (100.0 - fmin(fmax((float)percentage, 0.0), 100.0)) / 100.0;
    for (int i = 0; i < size; ++i)
    {
        data[i] = data[i] + (255 - data[i]) * tint;
    }
    return *this;
}

// Transformations
Image &Image::flipX()
{
    uint8_t temp[4];
    uint8_t *p1;
    uint8_t *p2;
    for (int i = 0; i < h; ++i)
    {
        for (int j = 0; j < w / 2; ++j)
        {
            p1 = &data[(i * w + j) * channels];
            p2 = &data[(i * w + w - 1 - j) * channels];
            memcpy(temp, p1, channels);
            memcpy(p1, p2, channels);
            memcpy(p2, temp, channels);
        }
    }
    return *this;
}

Image &Image::flipY()
{
    int M = h, N = w * channels;
    for (int i = 0; i < h / 2; ++i)
    {
        for (int j = 0; j < N; ++j)
        {
            data[i * N + j] ^= data[(M - i - 1) * N + j];
            data[(M - i - 1) * N + j] ^= data[i * N + j];
            data[i * N + j] ^= data[(M - i - 1) * N + j];
        }
    }
    return *this;
}
