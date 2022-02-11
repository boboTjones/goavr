# GoAVR
This is the early, pre-librification version of the AVR emulator, written in Go, for StarFighter. 

## How to make it do something

```
go build .
./goavr -i -f testcases/abs.elf
```

It works roughly like gdb after that. Type s to step, n for next, q to quit, and (I think) b for break and c for continue. 

# Historical context
[Starfighter, Summer 2015](https://sockpuppet.org/blog/2015/07/13/starfighter/)
