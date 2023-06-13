#!/bin/sh
#jar名称
GO_NAME=xxxxxxx

#golang包路径
GO_PATH=`dirname $0`/

#PID  代表是PID文件
PID=$GO_NAME\.pid

#使用说明，用来提示输入参数
usage(){
  echo "使用: sh run.sh [start|stop|restart|status]"
  exit 1
}

#检查程序是否在运行
is_exist(){
  pid=`ps -ef|grep $GO_NAME|grep -v grep|awk '{print $2}'`
  #如果不存在返回1，存在返回0
  if [ -z "${pid}" ]; then
    return 1
  else
    return 0
  fi
}

#启动方法
start(){
  is_exist
  if [ $? -eq "0" ]; then
    echo ">>> $GO_NAME 正在运行中 PID:${pid} <<<"
  else
    nohup   $GO_PATH/$GO_NAME >> ./$GO_NAME.out 2>&1 &
    echo $! > $PID
    echo ">>> 启动 $GO_NAME 成功 PID:$! <<<"
   fi
  }

#停止方法
stop(){
  #is_exist
  pidf=$(cat $PID)
  #echo "$pidf"
  echo ">>> PID:$pidf 开始清除 $pidf <<<"
  kill $pidf
  rm -rf $PID
  sleep 2
  is_exist
  if [ $? -eq "0" ]; then
    echo ">>> 清除 PID:$pid 的进程  <<<"
    kill -9  $pid
    sleep 2
    echo ">>> $GO_NAME 进程已停止 <<<"
  else
    echo ">>> $GO_NAME 没有运行 <<<"
  fi
}

#输出运行状态
status(){
  is_exist
  if [ $? -eq "0" ]; then
    echo ">>> $GO_NAME 正在运行 PID:${pid} <<<"
  else
    echo ">>> $GO_NAME 没有运行 <<<"
  fi
}

#重启
restart(){
  stop
  start
}

#根据输入参数，选择执行对应方法，不输入则执行使用说明
case "$1" in
  "start")
    start
    ;;
  "stop")
    stop
    ;;
  "status")
    status
    ;;
  "restart")
    restart
    ;;
  *)
    usage
    ;;
esac
exit 0
