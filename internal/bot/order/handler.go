package bot_order

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

	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
	"github.com/Vovarama1992/go-bagdoor-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

var pendingPhotos = map[int64][]string{} // tgID → []photoURL
var awaitingPhotos = map[int64]int{}     // tgID → orderID

func HandleSetOrderID() tele.HandlerFunc {
	return func(c tele.Context) error {
		args := strings.Split(c.Message().Text, " ")
		if len(args) != 2 {
			return c.Send("Формат: /setorderid <ID>")
		}

		orderID, err := strconv.Atoi(args[1])
		if err != nil {
			return c.Send("❌ Некорректный ID")
		}

		tgID := c.Sender().ID
		awaitingPhotos[tgID] = orderID

		return c.Send(fmt.Sprintf("Готово. Прикрепите фото к заказу #%d", orderID))
	}
}
func HandlePhotoUpload(s *order.Service, bot *tele.Bot, uploader *storage.S3Uploader) tele.HandlerFunc {
	return func(c tele.Context) error {
		tgID := c.Sender().ID
		orderID, ok := awaitingPhotos[tgID]
		if !ok {
			return c.Send("❌ Сначала введите команду: /setorderid <номер>")
		}

		photo := c.Message().Photo
		if photo == nil {
			return c.Send("❌ Это не фото")
		}

		file, err := bot.FileByID(photo.FileID)
		if err != nil {
			log.Printf("Ошибка получения файла: %v", err)
			return c.Send("❌ Не удалось получить фото")
		}

		fileName := filepath.Base(file.FilePath)
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return c.Send("❌ Только .jpg, .jpeg и .png")
		}

		downloadURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
		resp, err := http.Get(downloadURL)
		if err != nil {
			log.Printf("Ошибка скачивания файла: %v", err)
			return c.Send("❌ Не удалось скачать файл")
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Ошибка чтения файла: %v", err)
			return c.Send("❌ Не удалось прочитать файл")
		}

		s3url, err := uploader.UploadOrderMedia(orderID, fileName, content)
		if err != nil {
			log.Printf("Ошибка загрузки в S3: %v", err)
			return c.Send("❌ Не удалось загрузить фото")
		}

		pendingPhotos[tgID] = append(pendingPhotos[tgID], s3url)
		count := len(pendingPhotos[tgID])

		if count > 5 {
			delete(pendingPhotos, tgID)
			delete(awaitingPhotos, tgID)
			return c.Send("❌ Слишком много фото. Допустимо не более 5. Начните заново с /setorderid.")
		}

		if count < 2 {
			return c.Send(fmt.Sprintf("📷 Фото %d/5 загружено. Ещё минимум %d", count, 2-count))
		}

		// Всё ок — сохраняем
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		urls := pendingPhotos[tgID]

		if err := s.AddMediaURLs(ctx, orderID, urls); err != nil {
			log.Printf("Ошибка при сохранении фото: %v", err)
			return c.Send("Ошибка при сохранении фото")
		}

		if err := s.UpdateModerationStatus(ctx, orderID, order.StatusApproved); err != nil {
			log.Printf("Ошибка при обновлении статуса: %v", err)
		}

		order, err := s.Repo.GetOrderByID(ctx, orderID)
		if err != nil {
			log.Printf("Ошибка получения заказа: %v", err)
			return c.Send("Ошибка при формировании подтверждения")
		}

		msg := fmt.Sprintf(`📦 %s
Название: %s
Описание: %s
Откуда: %s
Куда: %s
С %s по %s
Стоимость: %s
Вознаграждение: %.0f ₽

✅ Ваш заказ отправлен на модерацию! В скором времени он будет опубликован ✅`,
			order.OrderNumber,
			order.Title,
			order.Description,
			order.OriginCity,
			order.DestinationCity,
			order.StartDate.Format("02.01.2006"),
			order.EndDate.Format("02.01.2006"),
			formatNullablePrice(order.Cost),
			order.Reward,
		)

		delete(pendingPhotos, tgID)
		delete(awaitingPhotos, tgID)

		return c.Send(msg, &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{{Text: "↩ Вернуться назад в меню", Data: "menu_back"}},
			},
		})
	}
}

func formatNullablePrice(p *float64) string {
	if p == nil {
		return "—"
	}
	return fmt.Sprintf("%.0f ₽", *p)
}
