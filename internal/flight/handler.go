package flight

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

var awaitingMap = map[int64]int{} // tgID → flightID

func HandleSetFlightID() tele.HandlerFunc {
	return func(c tele.Context) error {
		args := strings.Split(c.Message().Text, " ")
		if len(args) != 2 {
			return c.Send("Формат: /setflightid <ID>")
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return c.Send("❌ Неверный ID")
		}

		awaitingMap[c.Sender().ID] = id

		return c.Send(fmt.Sprintf(
			"✈️ Рейс #%d\nПришлите маршрутную карту в формате PDF.",
			id,
		))
	}
}

func HandlePdfUpload(service *Service, bot *tele.Bot, uploader *storage.S3Uploader) tele.HandlerFunc {
	return func(c tele.Context) error {
		tgID := c.Sender().ID
		flightID, ok := awaitingMap[tgID]
		if !ok {
			return c.Send("❌ Сначала введите: /setflightid <номер>")
		}

		doc := c.Message().Document
		if doc == nil || !strings.HasSuffix(strings.ToLower(doc.FileName), ".pdf") {
			return c.Send("❗Неверный формат. Пришлите маршрутную карту в PDF.", &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{{Text: "↩ Вернуться назад в меню", Data: "menu_back"}},
				},
			})
		}

		file, err := bot.FileByID(doc.FileID)
		if err != nil {
			log.Printf("FileByID error: %v", err)
			return c.Send("❌ Не удалось получить файл")
		}

		url := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Download error: %v", err)
			return c.Send("❌ Не удалось скачать файл")
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Read error: %v", err)
			return c.Send("❌ Не удалось прочитать файл")
		}

		if len(content) < 30*1024 || !strings.HasPrefix(http.DetectContentType(content), "application/pdf") {
			return c.Send("❗Неправильный формат или слишком маленький файл. Попробуйте ещё раз.", &tele.ReplyMarkup{
				InlineKeyboard: [][]tele.InlineButton{
					{{Text: "↩ Вернуться назад в меню", Data: "menu_back"}},
				},
			})
		}

		s3URL, err := uploader.UploadFlightMap(flightID, filepath.Base(file.FilePath), content)
		if err != nil {
			log.Printf("S3 upload error: %v", err)
			return c.Send("❌ Ошибка загрузки файла")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := service.SetMapURL(ctx, flightID, s3URL); err != nil {
			log.Printf("SetMapURL error: %v", err)
			return c.Send("❌ Ошибка сохранения карты")
		}

		if err := service.SetStatus(ctx, flightID, StatusApproved); err != nil {
			log.Printf("SetStatus error: %v", err)
		}

		delete(awaitingMap, tgID)

		return c.Send("✅ Маршрутная карта загружена. Рейс отправлен на модерацию.", &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "↩ Назад в меню", Data: "menu_back"}},
			},
		})
	}
}
