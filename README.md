# assets-uploader
![Go CI](https://github.com/wangyoucao577/assets-uploader/workflows/Go%20CI/badge.svg) ![Release Go Binaries](https://github.com/wangyoucao577/assets-uploader/workflows/Release%20Go%20Binaries/badge.svg)    
Command line tool to robustly upload Github release assets.  

## Features
- Upload file to Github Release Assets that identified by `tag`.    
- Allow `overwrite` if file exists.
- With the optional `-draft` flag an upload to an existing draft-release is supported (first pick)

## Build 
```bash
$ cd cmd/github-assets-uploader
$ go build
```

## Usage
- help
```bash
$ github-assets-uploader -h
Usage of ./github-assets-uploader:
  -f string
        File path to upload.
  -mediatype string
        E.g., 'application/zip'. (default "application/gzip")
  -overwrite
        Overwrite release asset if it's already exist.
  -repo string
        Github repo, e.g., 'wangyoucao577/assets-uploader'.
  -tag string
        Git tag to identify a Github Release in repo.
  -draft bool
        Upload asset to an existing draft release.
  -token string
        Github token to make changes.
  -version
        Print version and exit. 
```

- example    
```bash
$ github-assets-uploader -f vt2geojson-v0.1.5.1-testonly-linux-386.tar.gz -mediatype application/gzip -overwrite -repo wangyoucao577/vt2geojson -token *** -tag v0.1.5.1-testonly
```
