# Media Convertor GO

This tool helps you batch convert images and videos into smaller file sizes for efficient archiving.

## Features:

- Supports image formats: ImageMagick (hardcoded to AVIF output with 57 quality settings)
- Supports video formats: ffmpeg (hardcoded to MP4 output at 30 fps)
- Reduces file size while maintaining good quality
- Organizes processed files into separate folders for images and videos
- Renames processed files with creation date and resolution information

## Requirements:

- Go
- ExifTool
- ffmpeg
- ImageMagick

## Installation:

```sh
git clone https://github.com/lumpsoid/media-convertor-go.git
cd media-convertor-go
go build -o process_media
chmod +x ./process_media
./process_media -h
```

## How to Use:

```
./process_media -h
Usage of ./process_media:
  -ext string
    	File extension to process. In format: 'jpg,png,mov,mp4'
  -from_log_file string
    	Process files from previous failed run.
  -input_dir string
    	Input directory
  -min_video_dim int
    	Minimum video dimension convert to.
  -no_clean
    	Clean output directory before processing.
```

## Output:

Two new folders will be created within the source folder:
  - img - Contains compressed and renamed images.
  - mov - Contains compressed and renamed videos.

Filenames will be formatted as IMG/MOV_DateOfCreation_Resolution.EXT (e.g., 	
IMG_2018-01-21_13-11-59_342x480.avif or MOV_2017-12-20_18-57-17_853x480.mp4).

If a file fails to convert during the conversion process, the path to the file will be written to `.processError<unix time>`, after which you can run cli with the parameters `-from_log_file=<path/to/file> -no_clean` to try to convert only the failed files.

## Notes:

This cli currently uses hardcoded settings for output formats (AVIF for images and MP4 at 30 fps for videos). You can modify the `process_image.go` and `process_video.go` to customize these settings if needed.

## Disclaimer:

This script is provided as-is and the author is not responsible for any data loss or unexpected behavior. Always back up your original files before using this tool.
