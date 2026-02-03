# sshsync

Simple utility which I have written for myself for synchronization `authorized_keys` file on multiple machines.

It takes contents of file which stores on GitHub's Gist and just pushes into user's `authorized_keys`. Since ssh keys are public keys, it's totally safe to store them in open place.

## Usage

1. Download binary automatically

```bash
wget -O - https://raw.githubusercontent.com/holy-filipp/sshsync/refs/heads/main/install.sh | bash
```

2. Set URL of source file directly or using special flag

```bash
./sshsync do -u https://gist.github.com/holy-filipp/064b9318240c56ce4cb0c5f7908ff226/raw/authorized_keys
./sshsync url https://gist.github.com/holy-filipp/064b9318240c56ce4cb0c5f7908ff226/raw/authorized_keys
```

`sshsync url` saves URL to config file, program prints location of this file on screen.

### Other commands

1. `sshsync crontab [optional url]` - if URL specified, sets it to config, adds `sshsync do` to crontab for running each minute.

```bash
./sshsync crontab https://gist.github.com/holy-filipp/064b9318240c56ce4cb0c5f7908ff226/raw/authorized_keys
```

## Disclaimer

This utility I wrote in two days, if you really want to use this, be careful. I didn't even test it at the moment.