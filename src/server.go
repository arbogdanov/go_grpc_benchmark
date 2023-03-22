package main

import (
  "google.golang.org/grpc"
  "go_grpc_benchmark/build"
  "log"
  "net"
  "os"
  "os/signal"
  "sync"
  "syscall"
)


type Server1 struct {
}

func (server Server1) Handler(stream echo.EchoService_HandlerServer) error {
  for {
    request, err := stream.Recv()
    if err != nil {
      log.Printf("Recv is failed")
      break
    }
    reply := echo.Reply{Message: request.GetMessage()}
    if err := stream.Send(&reply); err != nil {
      log.Printf("Send error %v", err)
      break
    }
  }

  return nil
}

type Server2 struct {
}

func (server Server2) Handler(stream echo.EchoService_HandlerServer) error {
  chanel := make(chan *echo.Request, 100)

  go func(ch <-chan *echo.Request) {
    for {
      request, ok := <-ch
      if !ok {
        break
      }
      reply := echo.Reply{Message: request.GetMessage()}
      if err := stream.Send(&reply); err != nil {
        log.Printf("Send error %v", err)
        break
      }
    }
  }(chanel)

  for {
    request, err := stream.Recv()
    if err != nil {
      log.Printf("Recv is failed")
      close(chanel)
      break
    } else {
      chanel <- request
    }
  }

  return nil
}


func main() {
  lis, err := net.Listen("tcp", ":7777")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  server := grpc.NewServer()
  echo.RegisterEchoServiceServer(server, Server2{})

  sigs := make(chan os.Signal, 1)
  signal.Notify(sigs,
    syscall.SIGHUP,
    syscall.SIGINT,
    syscall.SIGTERM,
    syscall.SIGQUIT)

  var waitGroup sync.WaitGroup
  waitGroup.Add(1)

  go func(waitGroup *sync.WaitGroup) {
    defer waitGroup.Done()
    log.Println("Start server")
    if err := server.Serve(lis); err != nil {
      log.Fatalf("Failed to serve: %v", err)
    }
    log.Println("Server is stoped")
  }(&waitGroup)

  <-sigs
  server.Stop()

  waitGroup.Wait()
}
