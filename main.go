package main

import (
	"fmt"
	"os"
	"path/filepath"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func main() {
	// Директория с HTML файлами
	htmlDir := "./html_files"
	tempDir := "./temp_files"

	// Создаем временную директорию
	err := os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка при создании временной директории:", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Получаем список файлов в директории
	files, err := os.ReadDir(htmlDir)
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		return
	}

	// Создаем новый PDF генератор
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		fmt.Println("Ошибка при создании генератора PDF:", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".html" {
			filePath := filepath.Join(htmlDir, file.Name())

			// Чтение HTML файла
			htmlContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println("Ошибка при чтении файла:", err)
				continue
			}

			// Создание временного HTML файла
			tempFilePath := filepath.Join(tempDir, file.Name())
			err = os.WriteFile(tempFilePath, htmlContent, 0644)
			if err != nil {
				fmt.Println("Ошибка при создании временного файла:", err)
				continue
			}

			// Создание новой страницы из временного HTML файла
			page := wkhtml.NewPage(tempFilePath)
			pdfg.AddPage(page)
		}
	}

	// Установка общих опций для документа
	pdfg.PageSize.Set(wkhtml.PageSizeA4)
	pdfg.Dpi.Set(300)

	// Создание PDF документа
	err = pdfg.Create()
	if err != nil {
		fmt.Println("Ошибка при создании PDF:", err)
		return
	}

	// Сохранение PDF в файл
	outputFileName := "output.pdf"
	err = pdfg.WriteFile(outputFileName)
	if err != nil {
		fmt.Println("Ошибка при сохранении PDF:", err)
		return
	}

	fmt.Printf("PDF файл успешно создан %s\n", outputFileName)
}
