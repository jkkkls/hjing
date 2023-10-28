# procs

1. 进程管理
2. 进程日志收集

## 配置

```yaml
procs:
    -
        name: game
        cmd: ./game -app=game.yaml
        dir: ~/myapp/build
        log: ~/myapp/build/log/game_{date}.log
        stop_restart: true
        stop_signal: 2
        show_screen: true
net:
    port: 8888
```

## 命令

```shell
# 启动
procs -config=run.yaml
procs -config=run.yaml -daemon=true

#
procs shotdown
procs stop <proc_name>
procs restartall
procs restart <proc_name>
procs start <proc_name>
procs status

```