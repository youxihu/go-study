package util

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
	"go-study/app/database/tabtaba/internal/entity"
	"strconv"
	"time"
)

// ParseExcel reads data from an Excel file and returns populated ProjectRepo and FinanceRepo instances
func ParseExcel(excelPath string, log *widget.Entry) ([]*entity.ExcelParseData, error) {

	// 解析 Excel 文件
	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		log.Append(err.Error())
		return nil, err
	}
	// 延迟执行函数，确保在函数结束时关闭Excel文件
	defer func(f *excelize.File) {
		// 确保正确处理关闭错误
		if err := f.Close(); err != nil {
			log.Append(fmt.Sprintf("关闭Excel文件失败: %v\n", err))
		}
	}(f)
	var reply = make([]*entity.ExcelParseData, 0)

	// Assume the data is on the first sheet
	sheets := f.GetSheetList() // Get the name of the first sheet
	log.Append(fmt.Sprintf("文件中共有 %d 个工作表\n", len(sheets)))
	// 逐行输出每个工作表名称
	for i, sheet := range sheets {
		log.Append(fmt.Sprintf("工作表 %d: %s\n", i+1, sheet))
	}
	// 获取工作表名称
	sheetName := f.GetSheetName(0)
	log.Append(fmt.Sprintf("正在处理工作表: %s\n", sheetName))
	log.Append("数据导入中.......\n")
	// 获取当前工作表的所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Append(fmt.Sprintf("获取工作表 '%s' 数据行失败: %v\n", sheetName, err))
		return reply, nil
	}
	// 遍历所有行
	for i, row := range rows {
		// 跳过标题行
		if i < 1 {
			continue
		}
		// 如果当前行的元素数小于 20，则打印错误信息并跳过当前行
		if len(row) < 20 {
			log.Append(fmt.Sprintf("第 %d 行数据插入失败：元素数量不足，只有 %d 个元素\n", i+1, len(row)))
			continue // 跳过当前行，继续处理下一行
		}
		// Parse data from the Excel row
		accountId, _ := strconv.Atoi(row[2])                                                     //  row[2] is accountId
		amount, _ := strconv.ParseFloat(row[16], 64)                                             //  row[16] is amount
		createTime, _ := time.Parse("2006-01-02", row[1])                                        //  row[1] is createTime
		expectDeliverTime, _ := time.Parse("2006-01-02", row[9])                                 //  row[9] is expectDeliverTime
		remark := fmt.Sprintf("在 %s 导入了项目 %s", time.Now().Format("2006-01-02 15:04:05"), row[7]) //  row[7] is remark
		transactAt, _ := time.Parse("2006-01-02", row[19])                                       //  row[19] is transactAt
		title := row[7]
		projectType := FormatProjectType(amount)
		summary := fmt.Sprintf("支付【%s】款项", title)
		workerAccountId, _ := strconv.Atoi(row[5])

		reply = append(reply, &entity.ExcelParseData{
			AccountId:         int32(accountId),
			Amount:            amount,
			CreateTime:        createTime,
			ExpectDeliverTime: expectDeliverTime,
			Remark:            remark,
			TransactAt:        transactAt,
			Title:             title,
			ProjectType:       int8(projectType),
			Summary:           summary,
			WorkerAccountId:   int32(workerAccountId),
		})
	}
	return reply, nil
}
