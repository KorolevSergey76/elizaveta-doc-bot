package admin

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"cosmetologybotliza/internal/service"
)

/*
Этот файл — конструктор и входная точка для админ-панели.
Здесь вы связываете все необходимые компоненты для работы админа воедино.
*/

type Handler struct {
	Bot     *tgbotapi.BotAPI                //Сам объект бота, через который мы отправляем сообщения в Telegram.
	Service service.BookingServiceInterface //Использование интерфейса вместо конкретной реализации, поможет в будущем при желании заменить БД на другую или добавить дополнительные проверки в сервис, не придется менять код админки.
	User    *service.UserService            //Сервис пользователей, чтобы проверить, имеет ли право тот, кто написал /admin, действительно открывать админ-панель.
}

/*
NewHandler нужен, чтобы настроить и «узаконить» связь между админкой и остальными частями бота.
Это способ сказать: «Админка, вот тебе все, что нужно для работы. Теперь ты готова принимать команды».
Если убрать NewHandler и передачу параметров, получится «код-спагетти», в котором всё будет запутано друг с другом,
и малейшее изменение в одном месте будет ломать всё остальное.
*/
func NewHandler(
	bot *tgbotapi.BotAPI,
	svc service.BookingServiceInterface,
	userSvc *service.UserService,
) *Handler {
	return &Handler{
		Bot:     bot,
		Service: svc,
		User:    userSvc,
	}
}

/*
Это логический «фильтр»:
Здесь проверяется команда /admin.
Если она найдена — вызывается метод OpenAdminMenu
*/
func (h *Handler) HandleMessage(msg *tgbotapi.Message) {
	if msg.Text == "/admin" {
		h.OpenAdminMenu(msg.Chat.ID, 0)
	}
}
