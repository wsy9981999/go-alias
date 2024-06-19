package main

import (
	"context"
	"log"
	"os"

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
		log.Fatal("usage:alias [name] [value]")
	}
	name := os.Args[1]
	value := os.Args[2]

	newContent := gstr.Replace(content, "{{value}}", value)
	path := gfile.Join(gfile.SelfDir(), "src", name)
	if err := gfile.Mkdir(path); err != nil {
		log.Fatal(err.Error())
	}
	if err := gfile.PutContents(gfile.Join(path, "main.go"), newContent); err != nil {
		log.Fatal(err.Error())
	}
	if err := run(gctx.GetInitCtx(), "go mod init "+name, path); err != nil {
		log.Fatal(err.Error())
	}
	if err := run(gctx.GetInitCtx(), "go mod tidy", path); err != nil {
		log.Fatal(err.Error())
	}
	if err := run(gctx.GetInitCtx(), "go install", path); err != nil {
		log.Fatal(err.Error())
	}

}
func run(ctx context.Context, command string, path string) error {
	cmd := gproc.NewProcessCmd(command)
	cmd.Dir = path
	return cmd.Run(ctx)
}
