# Media Convertor GO

This tool helps you batch convert images and videos into smaller file sizes for efficient archiving.

## Features:

- Supports image formats: check ImageMagick 
- Supports video formats: check ffmpeg
- Safely converting files
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
go build -o process_media ./cmd/process-media/main.go
chmod +x ./process_media
./process_media -h
```

## How to Use:

```
./process_media -h
  -extensions string
        Comma-separated list of file extensions to process (e.g., 'jpg,png,mov,mp4')
  -fromLogFile string
        Path to the log file with file paths from a previous run
  -imageTargetFormat string
        Target format for processed images (e.g., 'jpg', 'png')
  -imageTargetQuality int
        Quality of the processed images (e.g., 80, 57)
  -inputDir string
        Directory containing input files. Must be flat
  -logFilePath string
        Path to the log file for recording processing details
  -overrideOutputDir
        Whether to override the output directory
  -videoMinDimension int
        Minimum dimension of the video to process
  -videoTargetFormat string
        Target format for processed videos (e.g., 'mp4')
  -videoTargetFps int
        Target frames per second for video processing
```

## Output:

Two new folders will be created within the source folder:
  - img - Contains compressed and renamed images, will be deleted at the end.
  - mov - Contains compressed and renamed videos, will be deleted at the end.
  - structured - Containing all media with Year/Month/Files structure.

Filenames will be formatted as IMG/MOV_DateOfCreation_Resolution.EXT (e.g., 	
`IMG_2018-01-21_13-11-59_342x480.avif` or `MOV_2017-12-20_18-57-17_853x480.mp4`).

If a file fails to convert during the conversion process, the path to the file will be written to `logfile` specified with `-logFilePath` flag, after which you can run CLI with the parameters `-fromLogFilePath=<path/to/file>` to try to convert only the failed files.

## Disclaimer:

This script is provided as-is and the author is not responsible for any data loss or unexpected behavior. Always back up your original files before using this tool.
