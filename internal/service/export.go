package service

import (
	"bytes"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type ExportService struct {
	docService DocumentServiceInterface
}

func NewExportService(docService DocumentServiceInterface) *ExportService {
	return &ExportService{
		docService: docService,
	}
}

type DocumentData struct {
	Number      string
	Date        string
	Location    string
	CreatedBy   string
	ApprovedBy  string
	Items       []ItemData
	TotalItems  int
	TotalPrice  float64
	DateCreated string
}

type ItemData struct {
	Number       int
	Name         string
	SerialNumber string
	Quantity     int
	Price        float64
	TotalPrice   float64
}

func (s *ExportService) ExportDocument(docID uint) ([]byte, error) {
	// Получаем документ
	response := s.docService.GetDocument(docID)
	if response.Model == nil {
		return nil, fmt.Errorf("документ не найден")
	}

	doc := response.Model

	// Проверяем наличие необходимых связей
	if doc.Location == nil {
		return nil, fmt.Errorf("местоположение документа не найдено")
	}
	if doc.CreatedBy == nil {
		return nil, fmt.Errorf("создатель документа не найден")
	}

	// Создаем новый Excel файл
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// Устанавливаем имя листа
	sheetName := "Документ"
	f.SetSheetName("Sheet1", sheetName)

	// Устанавливаем заголовок документа
	docType := map[string]string{
		"inventory":  "АКТ ИНВЕНТАРИЗАЦИИ",
		"transfer":   "АКТ ПЕРЕМЕЩЕНИЯ",
		"write_off":  "АКТ СПИСАНИЯ",
		"acceptance": "АКТ ПРИЕМКИ",
	}[doc.Type]

	// Форматирование заголовка
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("%s №%s", docType, doc.Number))
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("от %s", doc.Date.Format("02.01.2006")))
	f.SetCellValue(sheetName, "A3", fmt.Sprintf("Местонахождение: %s", doc.Location.Name))
	f.SetCellValue(sheetName, "A4", fmt.Sprintf("Создал: %s", doc.CreatedBy.Username))
	if doc.ApprovedBy != nil {
		f.SetCellValue(sheetName, "A5", fmt.Sprintf("Утвердил: %s", doc.ApprovedBy.Username))
	}

	// Заголовки таблицы
	headers := []string{"№", "Наименование", "Серийный номер", "Количество", "Цена", "Сумма"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c7", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Данные таблицы
	var totalItems int
	var totalPrice float64
	for i, item := range doc.Items {
		// Проверяем наличие оборудования
		if item.Equipment.ID == 0 {
			continue
		}

		row := i + 8
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), item.Equipment.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), item.Equipment.SerialNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), item.Quantity)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), item.Price)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), item.TotalPrice)

		totalItems += item.Quantity
		totalPrice += item.TotalPrice
	}

	// Итоги
	lastRow := len(doc.Items) + 8
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", lastRow+1), "Итого:")
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", lastRow+1), totalItems)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", lastRow+1), totalPrice)

	// Подписи
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", lastRow+3), fmt.Sprintf("Создал: _____________ %s", doc.CreatedBy.Username))
	if doc.ApprovedBy != nil {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", lastRow+4), fmt.Sprintf("Утвердил: _____________ %s", doc.ApprovedBy.Username))
	}
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", lastRow+5), fmt.Sprintf("Дата составления: %s", time.Now().Format("02.01.2006")))

	// Устанавливаем ширину столбцов
	f.SetColWidth(sheetName, "A", "A", 5)  // №
	f.SetColWidth(sheetName, "B", "B", 30) // Наименование
	f.SetColWidth(sheetName, "C", "C", 15) // Серийный номер
	f.SetColWidth(sheetName, "D", "D", 10) // Количество
	f.SetColWidth(sheetName, "E", "E", 12) // Цена
	f.SetColWidth(sheetName, "F", "F", 12) // Сумма

	// Сохраняем в буфер
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, fmt.Errorf("ошибка при сохранении файла: %v", err)
	}

	return buf.Bytes(), nil
}

// Шаблоны документов
const inventoryActTemplate = `
                                        УТВЕРЖДАЮ
                                        {{if .ApprovedBy}}
                                        ______________________ {{.ApprovedBy}}
                                        {{end}}
                                        "___" ____________ 20__ г.

                          АКТ ИНВЕНТАРИЗАЦИИ №{{.Number}}
                                от {{.Date}}

Местонахождение: {{.Location}}

Комиссия в составе:
Председатель: {{.CreatedBy}}
{{if .ApprovedBy}}Члены комиссии: {{.ApprovedBy}}{{end}}

произвела инвентаризацию материальных ценностей, находящихся на балансе организации.

В результате инвентаризации установлено следующее:

+-----+------------------------+---------------+----------+------------+-------------+
| №   | Наименование          | Серийный №   | Кол-во   | Цена, руб. | Сумма, руб. |
+-----+------------------------+---------------+----------+------------+-------------+
{{range .Items}}
| {{printf "%-4d" .Number}} | {{printf "%-22s" .Name}} | {{printf "%-13s" .SerialNumber}} | {{printf "%-8d" .Quantity}} | {{printf "%-10.2f" .Price}} | {{printf "%-11.2f" .TotalPrice}} |
{{end}}
+-----+------------------------+---------------+----------+------------+-------------+

Итого наименований: {{len .Items}}
Общее количество: {{.TotalItems}}
Общая стоимость: {{printf "%.2f" .TotalPrice}} руб.

Председатель комиссии: _____________ {{.CreatedBy}}
{{if .ApprovedBy}}
Члены комиссии:       _____________ {{.ApprovedBy}}
{{end}}

Дата составления: {{.DateCreated}}
`

const transferActTemplate = `
                                        УТВЕРЖДАЮ
                                        {{if .ApprovedBy}}
                                        ______________________ {{.ApprovedBy}}
                                        {{end}}
                                        "___" ____________ 20__ г.

                          АКТ ПЕРЕМЕЩЕНИЯ №{{.Number}}
                                от {{.Date}}

Местонахождение: {{.Location}}

Материально-ответственное лицо: {{.CreatedBy}}

Настоящий акт составлен о перемещении следующих материальных ценностей:

+-----+------------------------+---------------+----------+------------+-------------+
| №   | Наименование          | Серийный №   | Кол-во   | Цена, руб. | Сумма, руб. |
+-----+------------------------+---------------+----------+------------+-------------+
{{range .Items}}
| {{printf "%-4d" .Number}} | {{printf "%-22s" .Name}} | {{printf "%-13s" .SerialNumber}} | {{printf "%-8d" .Quantity}} | {{printf "%-10.2f" .Price}} | {{printf "%-11.2f" .TotalPrice}} |
{{end}}
+-----+------------------------+---------------+----------+------------+-------------+

Итого наименований: {{len .Items}}
Общее количество: {{.TotalItems}}
Общая стоимость: {{printf "%.2f" .TotalPrice}} руб.

Сдал:    _____________ {{.CreatedBy}}
{{if .ApprovedBy}}
Принял:  _____________ {{.ApprovedBy}}
{{end}}

Дата составления: {{.DateCreated}}
`

const writeOffActTemplate = `
                                        УТВЕРЖДАЮ
                                        {{if .ApprovedBy}}
                                        ______________________ {{.ApprovedBy}}
                                        {{end}}
                                        "___" ____________ 20__ г.

                          АКТ СПИСАНИЯ №{{.Number}}
                                от {{.Date}}

Местонахождение: {{.Location}}

Комиссия в составе:
Председатель: {{.CreatedBy}}
{{if .ApprovedBy}}Члены комиссии: {{.ApprovedBy}}{{end}}

произвела проверку состояния материальных ценностей и установила необходимость их списания:

+-----+------------------------+---------------+----------+------------+-------------+
| №   | Наименование          | Серийный №   | Кол-во   | Цена, руб. | Сумма, руб. |
+-----+------------------------+---------------+----------+------------+-------------+
{{range .Items}}
| {{printf "%-4d" .Number}} | {{printf "%-22s" .Name}} | {{printf "%-13s" .SerialNumber}} | {{printf "%-8d" .Quantity}} | {{printf "%-10.2f" .Price}} | {{printf "%-11.2f" .TotalPrice}} |
{{end}}
+-----+------------------------+---------------+----------+------------+-------------+

Итого наименований: {{len .Items}}
Общее количество: {{.TotalItems}}
Общая стоимость: {{printf "%.2f" .TotalPrice}} руб.

Заключение комиссии: Указанные материальные ценности подлежат списанию.

Председатель комиссии: _____________ {{.CreatedBy}}
{{if .ApprovedBy}}
Члены комиссии:       _____________ {{.ApprovedBy}}
{{end}}

Дата составления: {{.DateCreated}}
`

const acceptanceActTemplate = `
                                        УТВЕРЖДАЮ
                                        {{if .ApprovedBy}}
                                        ______________________ {{.ApprovedBy}}
                                        {{end}}
                                        "___" ____________ 20__ г.

                          АКТ ПРИЕМКИ №{{.Number}}
                                от {{.Date}}

Местонахождение: {{.Location}}

Комиссия в составе:
Председатель: {{.CreatedBy}}
{{if .ApprovedBy}}Члены комиссии: {{.ApprovedBy}}{{end}}

произвела приемку следующих материальных ценностей:

+-----+------------------------+---------------+----------+------------+-------------+
| №   | Наименование          | Серийный №   | Кол-во   | Цена, руб. | Сумма, руб. |
+-----+------------------------+---------------+----------+------------+-------------+
{{range .Items}}
| {{printf "%-4d" .Number}} | {{printf "%-22s" .Name}} | {{printf "%-13s" .SerialNumber}} | {{printf "%-8d" .Quantity}} | {{printf "%-10.2f" .Price}} | {{printf "%-11.2f" .TotalPrice}} |
{{end}}
+-----+------------------------+---------------+----------+------------+-------------+

Итого наименований: {{len .Items}}
Общее количество: {{.TotalItems}}
Общая стоимость: {{printf "%.2f" .TotalPrice}} руб.

Заключение комиссии: Материальные ценности соответствуют сопроводительным документам.

Председатель комиссии: _____________ {{.CreatedBy}}
{{if .ApprovedBy}}
Члены комиссии:       _____________ {{.ApprovedBy}}
{{end}}

Дата составления: {{.DateCreated}}
`
