package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"go-study/app/database/tabtaba/internal/service"
	ctx2 "go-study/app/database/tabtaba/pkg/ctx"
	"go-study/app/database/tabtaba/pkg/util"
)

func (w *Window) NewForm() fyne.App {
	var dbEnv, excelPath, selectedProduct string
	// 创建应用
	myApp := w
	myWindow := myApp.NewWindow("黄乐东私人定制作案工具")
	// 创建日志输出框
	logOutput := widget.NewMultiLineEntry()
	logOutput.SetPlaceHolder("1.输入或选择你要作案的环境名称\n2.选择你要作案的项目\n3.浏览并选择你要作案的表格文件\n")

	// 环境选择
	env := []string{"dev", "public", "rc", "prod"}
	selectEnv := widget.NewSelectEntry(env)
	selectEnv.SetPlaceHolder("请选择环境") // 设置默认提示词
	selectEnv.OnChanged = func(selectedEnv string) {
		// 判断是否为空，或者是否为有效环境
		if selectedEnv == "" {
			logOutput.Append("请选择一个有效的环境。\n")
			return
		}

		dialog.ShowConfirm("确认选择", fmt.Sprintf("你确定要选择环境: %s 吗？", selectedEnv),
			func(confirmed bool) {
				if confirmed {
					dbEnv = selectedEnv
					logOutput.Append(fmt.Sprintf("已选择 %s 环境\n", dbEnv))
				} else {
					logOutput.SetText(fmt.Sprintf("已取消选择%s环境，请重新选择环境。\n", selectedEnv))
				}
			}, myWindow)
	}

	// 产品选择
	products := []string{"bbx", "bbz"}
	selectProduct := widget.NewSelect(products, func(product string) {
		dialog.ShowConfirm("确认选择", fmt.Sprintf("你确定要选择产品: %s 吗？", product),
			func(confirmed bool) {
				if confirmed {
					selectedProduct = product
					logOutput.Append(fmt.Sprintf("已选择产品: %s\n", selectedProduct))
				}

			}, myWindow)
	})
	// 文件选择按钮
	addFileButton := widget.NewButton("浏览或选择作案证据", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()
			excelPath = reader.URI().Path()
			logOutput.Append(fmt.Sprintf("文件已选择: %s\n请确认所有要素后按下启动按钮进行操作\n", excelPath))
		}, myWindow)
	})

	// 启动按钮
	startButton := widget.NewButton("启动！", func() {
		if dbEnv == "" || excelPath == "" || selectedProduct == "" {
			dialog.ShowInformation("参数错误", "请确认选择了环境、文件和产品", myWindow)
			return
		}

		logOutput.SetText(fmt.Sprintf("环境设置为: %s\n文件路径为: %s\n选择的产品: %s\n", dbEnv, excelPath, selectedProduct))

		// 解析 Excel 文件中的数据
		excelData, err := util.ParseExcel(excelPath, logOutput)
		if err != nil {
			dialog.ShowError(fmt.Errorf("解析 Excel 文件失败: %v", err), myWindow)
			logOutput.Append(fmt.Sprintf("解析 Excel 文件失败: %v\n", err))
			return
		}

		// 调用service 中的导入方法
		ctx := ctx2.NewCtx(selectedProduct, dbEnv)
		srv := service.NewService(dbEnv)
		if err = srv.ExportOfflineProject(ctx, excelData); err != nil {
			dialog.ShowError(fmt.Errorf("解析 Excel 文件失败: %v", err), myWindow)
			logOutput.Append(fmt.Sprintf("解析 Excel 文件失败: %v\n", err))
			return
		}

		// 成功信息
		logOutput.Append("数据导入完成！\n")
	})

	// 布局设置：左侧包含环境选择、产品选择和文件选择按钮
	leftLayout := container.NewVBox(
		widget.NewLabel("选择作案环境:"),
		selectEnv,
		widget.NewLabel("选择作案产品:"),
		selectProduct,
		addFileButton,
	)

	// 底部包含启动按钮
	bottomLayout := container.NewVBox(startButton)

	// 总体布局：左侧按钮 + 右侧日志框 + 底部按钮
	content := container.NewHSplit(leftLayout, logOutput)
	content.SetOffset(0.3) // 调整左右比例

	// 创建最终的布局，将底部按钮加到窗口底部
	finalLayout := container.NewBorder(nil, bottomLayout, nil, nil, content)

	// 设置窗口内容并显示
	myWindow.SetContent(finalLayout)
	myWindow.Resize(fyne.NewSize(680, 640)) // 初始窗口大小
	myWindow.ShowAndRun()
	myWindow.Show()

	return myApp
}
