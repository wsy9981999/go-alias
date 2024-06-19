package main

import (
	"context"
	"log"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
)

const content = `package main
import (
	"strings"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
)
func main(){
    ctx:=gctx.GetInitCtx()
	pwd:=gfile.Pwd()
	commandline:="{{value}} " + strings.Join(os.Args[1:]," ")
    cmd:=gproc.NewProcessCmd(commandline)
	cmd.Dir=pwd
	g.Log().Infof(ctx, "run \"%s\" in \"%s\"",commandline,pwd)
	cmd.Run(gctx.GetInitCtx())
}
`

func main() {
	if len(os.Args) < 3 {
		log.Fatal("usage:alias `name` `value` [-r]")
	}

	p, err := gcmd.Parse(g.MapStrBool{
		"r,remove": false,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	name := p.GetArg(1).String()
	value := p.GetArg(2).String()

	newContent := gstr.Replace(content, "{{value}}", value)
	path := gfile.Join(gfile.SelfDir(), "src", name)
	if err := gfile.Mkdir(path); err != nil {
		log.Fatal(err.Error())
	}
	ctx := gctx.GetInitCtx()
	g.Log().Info(ctx, "start write main.go")
	if err := gfile.PutContents(gfile.Join(path, "main.go"), newContent); err != nil {
		log.Fatal(err.Error())
	}
	g.Log().Info(ctx, "write success")
	g.Log().Info(ctx, "start run `go mod init`")
	if err := run(gctx.GetInitCtx(), "go mod init "+name, path); err != nil {
		log.Fatal(err.Error())
	}
	g.Log().Info(ctx, "run success")

	g.Log().Info(ctx, "start run `go mod tidy`")
	if err := run(gctx.GetInitCtx(), "go mod tidy", path); err != nil {
		log.Fatal(err.Error())
	}
	g.Log().Info(ctx, "run success")

	g.Log().Info(ctx, "start run `go intall`")
	if err := run(gctx.GetInitCtx(), "go install", path); err != nil {
		log.Fatal(err.Error())
	}
	g.Log().Info(ctx, "run success")

	rm := p.GetOpt("remove")
	if rm != nil {
		g.Log().Info(ctx, "start remove src")
		if err := gfile.Remove(path); err != nil {
			log.Fatal(err.Error())

		}
		g.Log().Info(ctx, "remove success")

	}
}
func run(ctx context.Context, command string, path string) error {
	cmd := gproc.NewProcessCmd(command)
	cmd.Dir = path
	return cmd.Run(ctx)
}
