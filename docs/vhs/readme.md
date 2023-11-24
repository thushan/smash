# VHS Demo for `smash`

This directory contains the demo [vhs tapes](https://github.com/charmbracelet/vhs) for recording `smash`.

* [Themes](https://github.com/charmbracelet/vhs/blob/main/THEMES.md)

# Setup & Recording

## Windows via Docker

```bash
$ MSYS_NO_PATHCONV=1 docker run --rm -v $PWD:/vhs ghcr.io/charmbracelet/vhs demo.tape
```

## Linux

```bash
$ vhs install.tape
```

### Flush buffers before demo
```bash
 free && sync && echo 3 > /proc/sys/vm/drop_caches && free
```
