package main

import (
    "context"
    "fmt"
    "os/exec"
    "time"
)

type Result struct {
    err error
    output []byte
}

func main() {
   // 执行1个cmd,让它在一个协程里面去执行，让它执行2秒

   // 1秒的时候，我们杀死cmd
   var (
       ctx context.Context
       cancelFunc context.CancelFunc
       resultChan chan *Result
       res *Result
   )

   resultChan = make(chan *Result, 1000)

   ctx, cancelFunc = context.WithCancel(context.TODO())

   go func() {
       var (
           output []byte
           err error
           cmd *exec.Cmd
       )
    cmd = exec.CommandContext(ctx, "bash", "-c", "sleep 3; echo hello;")
    output, err = cmd.CombinedOutput()
    resultChan <- &Result{
        err: err,
        output:output,
    }

   }()

   time.Sleep(1 * time.Second)
   cancelFunc()
   res = <- resultChan
   fmt.Println(res.err, string(res.output))
}
