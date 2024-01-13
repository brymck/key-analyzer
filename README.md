key-analyzer
============

Work in progress for analyzing [my keylogger][keylogger]'s output.

Usage
-----

You should be able to build this in macOS with only system dependencies:

```bash
git clone https://github.com/brymck/key-analyzer.git
cd key-analyzer
go build .
./key-analyzer ~/keyboard.bin
```

[keylogger]: https://github.com/brymck/c-macos-keylogger
