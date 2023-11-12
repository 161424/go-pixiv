package main

import (
	"context"
	"fmt"
	e "github.com/chen/download_pixiv_pic/common/err"
	//hd "github.com/chen/download_pixiv_pic/pkg/handle"
	//_ "github.com/chen/download_pixiv_pic/pkg/psgr"
	"github.com/spf13/cobra"
	godebug "runtime/debug"
)

var cmdRoot = &cobra.Command{
	Use:               "gpixiv",
	Short:             "Download Pixiv Picture",
	Long:              "null",
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func tweakGoGC() {
	// lower GOGC from 100 to 50, unless it was manually overwritten by the user
	oldValue := godebug.SetGCPercent(50)
	if oldValue != 100 {
		godebug.SetGCPercent(oldValue)
	}
}

var ctx context.Context

//var lg = log.NewDefaultSlog()

func main() {
	tweakGoGC()
	//hd.MainLoop()

	err := cmdRoot.ExecuteContext(ctx)
	//l.Send(slog.LevelDebug, "123", 3)
	if err != nil {
		err = fmt.Errorf("错误代码：%d, 错误内容：%s", e.RootError, err)
		fmt.Println(err.Error())
	}
	//ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	//defer stop()
	//
	//<-ctx.Done()

	//timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

}
