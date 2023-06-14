package main

import (
	"context"
	_ "github.com/chen/download_pixiv_pic/common/conf"
	hd "github.com/chen/download_pixiv_pic/pkg/handle"
	_ "github.com/chen/download_pixiv_pic/pkg/psgr"
	"os"
	"os/signal"
)

func main() {

	hd.MainLoop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()
	<-ctx.Done()

	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

}
