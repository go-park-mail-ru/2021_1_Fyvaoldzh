package events

import (
	"kudago/models"
	"reflect"
	"sync"
	"testing"
)

var ExpectedEvents = models.Events{
	{126, "Музей Эмоций", "Музей Эмоций", "Музей Эмоций — авторский проект художника Алексея Сергиенко, который раскрывает его взгляд на мир эмоций и переживаний. В пространстве Музея соседствуют несколько тематических зон, которые соединены между собой специальными коридорами. Каждая зона посвящена определённой эмоции, а коридор символизирует переход от одной эмоции к другой. Погружаться в разные эмоциональные состояния помогают запахи, звуки и тактильные ощущения. Миссия Музея — помочь людям обратить более пристальное внимание на собственные эмоции и их разнообразие, дать толчок к развитию эмоциональной грамотности и эмоционального интеллекта.", "ежедневно 11:00–21:00", "Курская", "пер. Нижний Сусальный, БК «Арма», д. 5, стр. 18", "Музей", "945e08955c7685816acaf8a7cf99ac5b.png"},
	{127, "Проект VR Gallery", "Центр современного искусства «МАРС»", "Студия «АртДинамикс» представляет проект VR Gallery («Виртуальная галерея»), который станет вашим порталом в мир классического и современного искусства. Новейшие технологии переносят в иное пространство и обеспечивают полное погружение в мир творцов и их идей. Посетителям предлагаются четыре VR-путешествия на выбор. В рамках первого — Deep Immersive («Глубокое погружение») — вы взглянете на «Крик» Эдварда Мунка глазами Сандры Погам и Шарля Аятса, переосмыслите творчество Сальвадора Дали и полюбуетесь нежными кувшинками Клода Моне. Второй вариант — виртуально посетить Galactic Gallery («Галактическая галерея») и выставку Rone («Рон») и открыть для себя творчество талантливых мастеров стрит-арта и диджитал-художников. Третий вариант Beyond the Glass («За стеклом») даёт возможность узнать тайны картины «Мона Лиза» и прогуляться по собору Парижской Богоматери, увидеть храм изнутри в оригинальном интерьере XVIII века. Четвертая возможность — заглянуть в VR-музей Кремера, где собраны шедевры датского и фламандского искусства XVII века. Данный VR-проект проходит по сеансам — смотрите расписание на сайте и заранее приобретайте электронные билеты.", "1 октября 2020 – 4 апреля 2021	12:00–22:00", "Сухаревская", "пер. Пушкарёв, д. 5", "Галлерея", "fc7a55e3897601089ddc88bd3f13d2ac.jpg"},
	{128, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", ""},
}

var ExpectedEvent = models.Event{128, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", ""}

var ExpectedEventChangeID = models.Event{128, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", ""}

var ExpectedEventType = models.Events{
	{127, "Проект VR Gallery", "Центр современного искусства «МАРС»", "Студия «АртДинамикс» представляет проект VR Gallery («Виртуальная галерея»), который станет вашим порталом в мир классического и современного искусства. Новейшие технологии переносят в иное пространство и обеспечивают полное погружение в мир творцов и их идей. Посетителям предлагаются четыре VR-путешествия на выбор. В рамках первого — Deep Immersive («Глубокое погружение») — вы взглянете на «Крик» Эдварда Мунка глазами Сандры Погам и Шарля Аятса, переосмыслите творчество Сальвадора Дали и полюбуетесь нежными кувшинками Клода Моне. Второй вариант — виртуально посетить Galactic Gallery («Галактическая галерея») и выставку Rone («Рон») и открыть для себя творчество талантливых мастеров стрит-арта и диджитал-художников. Третий вариант Beyond the Glass («За стеклом») даёт возможность узнать тайны картины «Мона Лиза» и прогуляться по собору Парижской Богоматери, увидеть храм изнутри в оригинальном интерьере XVIII века. Четвертая возможность — заглянуть в VR-музей Кремера, где собраны шедевры датского и фламандского искусства XVII века. Данный VR-проект проходит по сеансам — смотрите расписание на сайте и заранее приобретайте электронные билеты.", "1 октября 2020 – 4 апреля 2021	12:00–22:00", "Сухаревская", "пер. Пушкарёв, д. 5", "Галлерея", "fc7a55e3897601089ddc88bd3f13d2ac.jpg"},
}

var ExpectedNewEvent = models.Event{1, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", ""}

func TestDeleteByIDOK(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	handler.DeleteByID(125)

	if !reflect.DeepEqual(ExpectedEvents, handler.Events) {
		t.Errorf("expected: [%v], got: [%v]", ExpectedEvents, handler.Events)
	}
}

func TestDeleteByIDERROR(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	if handler.DeleteByID(5) {
		t.Error("Your func just deleted non-existent event")
	}
}

func TestGetOneEventByIDOK(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	if !reflect.DeepEqual(handler.GetOneEventByID(128), ExpectedEvent) {
		t.Errorf("expected: [%v], got: [%v]", ExpectedEvent, handler.GetOneEventByID(128))
	}
}

func TestGetOneEventByIDERROR(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	if !reflect.DeepEqual(handler.GetOneEventByID(5), models.Event{}) {
		t.Errorf("expected: [%v], got: [%v]", models.Event{}, handler.GetOneEventByID(5))
	}
}

func TestGetEventsByTypeOK(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: ExpectedEvents,
		Mu:     &sync.Mutex{},
	}

	if !reflect.DeepEqual(handler.GetEventsByType("Галлерея"), ExpectedEventType) {
		t.Errorf("expected: [%v], got: [%v]", ExpectedEventType, handler.GetEventsByType("Галлерея"))
	}
}

func TestGetEventsByTypeERROR(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	var none models.Events

	if !reflect.DeepEqual(handler.GetEventsByType("NOTYPE"), none) {
		t.Errorf("expected: [%v], got: [%v]", none, handler.GetEventsByType("NOTYPE"))
	}
}

func TestCreateNewEventOK(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.Events{},
		Mu:     &sync.Mutex{},
	}

	handler.CreateNewEvent(&ExpectedNewEvent)

	if !reflect.DeepEqual(handler.Events[0], ExpectedNewEvent) {
		t.Errorf("expected: [%v], got: [%v]", ExpectedNewEvent, handler.Events[0])
	}

}

func TestCreateNewEventERROR(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.Events{},
		Mu:     &sync.Mutex{},
	}

	id := ExpectedEventChangeID.ID

	handler.CreateNewEvent(&ExpectedEventChangeID)

	if reflect.DeepEqual(handler.Events[0].ID, id) {
		t.Error("Your id didn't generate correct!")
	}

}
