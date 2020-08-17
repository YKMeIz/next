# NEXT

[![Go Report Card](https://goreportcard.com/badge/github.com/YKMeIz/Pill?style=flat-square)](https://goreportcard.com/report/github.com/YKMeIz/Pill)
[![License](https://img.shields.io/github/license/YKMeIz/Pill.svg?color=%232b2b2b&style=flat-square)](https://github.com/YKMeIz/Pill/blob/master/LICENSE)

NEXT is a dead-simple IT automation tool enabling concurrent commands execution on multiple machines. It is inspired by Ansible and GNU Make. NEXT is agentless, and running its tasks remotely via temporarily SSH connections.

## Configuration File

Configuration file of NEXT is in YAML format. Here is the simplest configuration file example:

```yaml
run:
  target:
    - 10.0.0.155
  script:
    - echo "This is $(hostname):"
    - uname -a
```

It defines target machine IP address and commands need to run, then save it as `next.yml` and execute following command. You will get result similar as below:
```bash
$ next run
2020/08/17 15:30:05 10.0.0.155: This is mashiro:
2020/08/17 15:30:05 10.0.0.155: Linux mashiro 4.18.0-147.5.1.el8_1.x86_64 #1 SMP Wed Feb 5 02:00:39 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux
```

More configuration file examples can be found in `examples` folder.