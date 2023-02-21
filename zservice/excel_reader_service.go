package zservice

import (
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
)

type ExcelReaderService struct {
}

func (*ExcelReaderService) Init(ctx context.Context) (bool, error) {
	var _, err = excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	return true, nil
}

func (*ExcelReaderService) GetName() string {
	return "excel"
}
