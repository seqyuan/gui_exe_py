//go build -ldflags="-H windowsgui -w -s"

package main

import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "log"
    "os/exec"

    "os"
    "bytes"
    "fmt"
    "strings"
)

func main(){
    mw := &MyMainWindow{}
    err := MainWindow{
        AssignTo: &mw.MainWindow, //窗口重定向至mw，重定向后可由重定向变量控制控件
        Title:   "选择python脚本并执行",
        MinSize:  Size{Width: 400, Height: 200},    //最小 尺寸
        MaxSize:  Size{Width: 600, Height: 200},   //最大尺寸
        Size:     Size{400, 400}, // 显示尺寸
        Layout:   VBox{}, //样式，纵向

        Children: []Widget{ //控件组
            TextEdit{
                AssignTo: &mw.Script2exe,
                MinSize:  Size{Width: 300, Height: 40},
                MaxSize:  Size{Width: 300, Height: 40},
            },
            PushButton{
                Text:      "选择要执行哪个程序",
                OnClicked: mw.SelectFile, //点击事件响应函数
            },
            PushButton{
                Text:    "执行上边选择的程序",
                OnClicked: mw.Exescript,
            },
            TextEdit{
                AssignTo: &mw.Start,
                MinSize:  Size{Width: 400, Height: 20},
                MaxSize:  Size{Width: 400, Height: 20},
            },
            ListBox{
                AssignTo: &mw.process,
                Row:      5,
            },
        },
    }.Create() //创建

    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    mw.Run() //运行
}


func (mw *MyMainWindow) Exescript() {
    model := []string{}
    mw.Start.AppendText("开始执行程序了 ...")

    model_s, stderr, err := mw.CatchPyResult()
    for _, i := range strings.Split(model_s, "\n") {
        model = append(model, i)
    }
    mw.process.SetModel(model)

    if err != nil {
        mw.Start.SetText("程序执行错误")
        for _, ii := range strings.Split(string(stderr), "\n") {
            model = append(model, ii)
        }
    }else{
        mw.Start.SetText("程序执行完成")
    }

    mw.process.SetModel(model)    
}

func (mw *MyMainWindow) SelectFile() {
    dlg := new(walk.FileDialog)
    dlg.Title = "选择要执行的程序"
    dlg.Filter = "可执行文件 (*.py)|*.py|所有文件 (*.*)|*.*"

    mw.Script2exe.SetText("") 

    if ok, err := dlg.ShowOpen(mw); err != nil {
        mw.Script2exe.AppendText("Error : File Open\r\n")
        return
    } else if !ok {
        mw.Script2exe.AppendText("Cancel\r\n")
        return
    }

    s := fmt.Sprintf("%s", dlg.FilePath)
    mw.Script2exe.AppendText(s)
    mw.ScriptPath = s
}

type MyMainWindow struct {
    *walk.MainWindow
    process   *walk.ListBox
    Start    *walk.TextEdit
    Script2exe *walk.TextEdit
    ScriptPath  string
}

func (mw *MyMainWindow) CatchPyResult() (string, string, error){
    cmd := exec.Command("python", mw.ScriptPath)
    // 获取输出对象，可以从该对象中读取输出结果
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Start()
    if err != nil {
        log.Fatal(string(stderr.Bytes()), err)
    }
    _ = cmd.Wait()

    return string(stdout.Bytes()), string(stderr.Bytes()), err
}

