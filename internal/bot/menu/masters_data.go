package menu

type MasterInfo struct {
	Name    string
	Bio     string
	PhotoID string
}

// Список наших мастеров
var Masters = map[string]MasterInfo{
	"master_1": {
		Name: "Королева Елизавета",
		Bio: "👑 *Королева Елизавета — врач-косметолог*\n\n" +
			"Косметология для меня — это не просто процедуры, а искусство сохранения вашей природной уникальности и ресурса молодости.\n\n" +
			"• *Мой подход:* Доказательная медицина и бережная эстетика.\n\n" +
			"✨ *С чем ко мне приходят:*\n" +
			"— Экспертный уход за кожей и лечение акне;\n" +
			"— Безоперационное омоложение.",
		PhotoID: "AgACAgIAAxkBAAIERWoZ9wRca7-jwVh09lO6jaToCEOoAAI7GWsbx93RSL8dyZ-395zeAQADAgADeQADOwQ",
	},
	"master_2": {
		Name: "Лапина Полина",
		Bio: "✨ *Лапина Полина — врач-косметолог*\n\n" +
			"Красота в деталях, а здоровье кожи — в профессиональном подходе.\n\n" +
			"• *Мой подход:* Комфорт пациента и видимый результат.\n\n" +
			"🌟 _Моя цель — чтобы вы влюблялись в свое отражение каждый день!_",
		PhotoID: "AgACAgIAAxkBAAIERmoZ9wx1SUmbvKGUakoXZipmAAF2ZwACPBlrG8fd0UhMztycSLYbeAEAAwIAA3kAAzsE",
	},
}
