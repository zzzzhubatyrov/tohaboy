package service

import (
    "bytes"
    "fmt"
    "time"

    "github.com/xuri/excelize/v2"
)

// ExportDocumentGOST создает Excel-файл по официальной форме (Гост) для указанного документа
func (s *ExportService) ExportDocumentGOST(docID uint) ([]byte, error) {
    // Получаем документ как обычно
    response := s.docService.GetDocument(docID)
    if response.Model == nil {
        return nil, fmt.Errorf("документ не найден")
    }
    doc := response.Model

    if doc.Location == nil {
        return nil, fmt.Errorf("местоположение документа не найдено")
    }
    if doc.CreatedBy == nil {
        return nil, fmt.Errorf("создатель документа не найден")
    }

    f := excelize.NewFile()
    defer func() {
        _ = f.Close()
    }()

    sheet := "Форма"
    f.SetSheetName("Sheet1", sheet)

    // Шапка в зависимости от типа документа
    var formName string
    switch doc.Type {
    case "inventory":
        formName = "Форма № ИНВ-1"
    case "transfer":
        formName = "Форма № МБ-1"
    case "write_off":
        formName = "Форма № МБ-8"
    case "acceptance":
        formName = "Форма № ОС-1"
    default:
        formName = "Форма документа"
    }

    f.SetCellValue(sheet, "A1", formName)
    f.MergeCell(sheet, "A1", "F1")
    f.SetCellValue(sheet, "A2", fmt.Sprintf("Акт № %s", doc.Number))
    f.MergeCell(sheet, "A2", "F2")
    f.SetCellValue(sheet, "A3", fmt.Sprintf("от %s", doc.Date.Format("02.01.2006")))
    f.MergeCell(sheet, "A3", "F3")

    // Заголовки таблицы (приближены к ГОСТ)
    headers := []string{"№", "Наименование", "Серийный номер", "Ед.изм.", "Кол-во", "Сумма"}
    for i, h := range headers {
        cell := fmt.Sprintf("%c5", 'A'+i)
        f.SetCellValue(sheet, cell, h)
    }

    // Данные таблицы
    var totalQty int
    var totalPrice float64
    for i, item := range doc.Items {
        if item.Equipment.ID == 0 {
            continue
        }
        row := i + 6
        f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
        f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.Equipment.Name)
        f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.Equipment.SerialNumber)
        f.SetCellValue(sheet, fmt.Sprintf("D%d", row), "шт.")
        f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Quantity)
        f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.TotalPrice)
        totalQty += item.Quantity
        totalPrice += item.TotalPrice
    }

    lastRow := len(doc.Items) + 6
    f.SetCellValue(sheet, fmt.Sprintf("A%d", lastRow), "Итого:")
    f.SetCellValue(sheet, fmt.Sprintf("E%d", lastRow), totalQty)
    f.SetCellValue(sheet, fmt.Sprintf("F%d", lastRow), totalPrice)

    // Подписи
    f.SetCellValue(sheet, fmt.Sprintf("A%d", lastRow+2), fmt.Sprintf("Создал: _____________ %s", doc.CreatedBy.Username))
    if doc.ApprovedBy != nil {
        f.SetCellValue(sheet, fmt.Sprintf("A%d", lastRow+3), fmt.Sprintf("Утвердил: _____________ %s", doc.ApprovedBy.Username))
    }
    f.SetCellValue(sheet, fmt.Sprintf("A%d", lastRow+4), fmt.Sprintf("Дата составления: %s", time.Now().Format("02.01.2006")))

    var buf bytes.Buffer
    if err := f.Write(&buf); err != nil {
        return nil, fmt.Errorf("ошибка при сохранении файла: %v", err)
    }
    return buf.Bytes(), nil
}
