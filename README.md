# fileslice

## 文件切片和恢复

## 切片

```bash
go run main.go -t -f <filepath> -s 20
```

## 恢复

```bash
go run main.go -f <filepath>
```
## 字段说明

```bash
  -f string
        Filename
  -m    Print file hash values
  -s int
        Chunk size in MB (default 10)
  -t    Split the file into chunks
  -v    Show the app version
```
