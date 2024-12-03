# advent-of-code (2024)

Base repository for the Advent of Code tasks (usually there are 2 parts per task per day).

## how to run (Go)

+ Go 1.23 runtime

```shell
go run 02/part_02.go -inputFile 02/input02.txt
```

Or using the Dcoker engine and official Go image:

```shell
docker run --rm -v "./:/opt/aoc" golang:1.23 go run /opt/aoc/02/part_02.go -inputFile /opt/aoc/02/input02.txt
```
