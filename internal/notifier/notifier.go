package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Vovarama1992/go-bagdoor-bot/internal/order"
)

type Notifier struct {
	BotToken string
}

func NewNotifier(token string) *Notifier {
	return &Notifier{BotToken: token}
}

func (n *Notifier) NotifyNewOrder(tgID int64, o *order.Order) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.BotToken)

	msg := fmt.Sprintf(`📦 %s
Название: %s
Описание: %s
Откуда: %s
Куда: %s
С %s по %s
Стоимость: %s
Вознаграждение: %.0f ₽

❗Ваш заказ в одном шаге от публикации. Пришлите 2–4 фото вашей посылки для завершения оформления.❗

/setorderid %d`,
		o.OrderNumber,
		o.Title,
		o.Description,
		o.OriginCity,
		o.DestinationCity,
		o.StartDate.Format("02.01.2006"),
		o.EndDate.Format("02.01.2006"),
		formatNullablePrice(o.Cost),
		o.Reward,
		o.ID,
	)

	body := map[string]interface{}{
		"chat_id": tgID,
		"text":    msg,
	}

	data, _ := json.Marshal(body)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram error: %d", resp.StatusCode)
	}

	return nil
}

func formatNullablePrice(p *float64) string {
	if p == nil {
		return "—"
	}
	return fmt.Sprintf("%.0f ₽", *p)
}
