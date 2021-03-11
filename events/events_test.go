package events

import (
	"bytes"
	"kudago/models"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/require"
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

var EventWithNoImage = models.Events{{1, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", ""}}

var ExpectedEventImage = models.Event{1, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", "1.jpg"}

var EventCreate = models.Event{1, "Пример", "Место", "Пример без картинки", "12:00", "Примерное метро", "Примерная улица", "Музей", "1.jpg"}

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

	event, err := handler.GetOneEventByID(128)

	if !reflect.DeepEqual(event, ExpectedEvent) || err != nil {
		t.Errorf("expected: [%v], got: [%v]", ExpectedEvent, event)
	}
}

func TestGetOneEventByIDERROR(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}

	event, err := handler.GetOneEventByID(5)

	if !reflect.DeepEqual(event, models.Event{}) && err == nil {
		t.Errorf("expected: [%v], got: [%v]", models.Event{}, event)
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

func TestSaveImage(t *testing.T) {
	t.Parallel()
	handler := Handlers{
		Events: EventWithNoImage,
		Mu:     &sync.Mutex{},
	}
	img, err := os.Open("../2913aa38efbe34ebdb5a1d642dfa29d8.jpg")
	if err != nil {
		t.Error("No such file!")
	}
	handler.SaveImage(img, 1)
	if !reflect.DeepEqual(handler.Events[0], ExpectedEventImage) {
		t.Errorf("expected: [%v], got: [%v]", ExpectedEventImage, handler.Events[0])
		os.Remove("1.jpg")
	}
	os.Remove("1.jpg")
}

func setupEcho(t *testing.T, url, method string) (echo.Context,
	Handlers) {
	e := echo.New()
	var req *http.Request
	switch method {
	case http.MethodPost:
		f, _ := EventCreate.MarshalJSON()
		req = httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(f))
	case http.MethodGet:
		req = httptest.NewRequest(http.MethodGet, url, nil)
	}
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(url)

	uh := Handlers{
		Events: models.BaseEvents,
		Mu:     &sync.Mutex{},
	}
	return c, uh
}

func TestHandler_Main(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/", http.MethodGet)

	err := uh.GetAllEvents(c)
	require.Equal(t, nil, err)
}

func TestHandler_OneEventOK(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("127")

	err := uh.GetOneEvent(c)
	require.Equal(t, nil, err)
}

func TestHandler_OneEventERROR(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("aaa")

	err := uh.GetOneEvent(c)
	require.NotEqual(t, nil, err)
}

func TestHandler_OneEvent1ERROR(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/event/:id", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("545654")

	err := uh.GetOneEvent(c)
	require.NotEqual(t, nil, err)
}

/*func TestHandler_GetImageOK(t *testing.T) {
	c, uh := setupEcho(t, "api/v1/event/:id/image", http.MethodGet)
	c.SetParamNames("id")
	c.SetParamValues("127")

	err := uh.GetImage(c)
	require.NotEqual(t, nil, err)
}*/

func TestHandler_Events(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/event", http.MethodGet)
	c.QueryParams().Add("typeEvent", "Музей")
	err := uh.GetEvents(c)
	require.Equal(t, err, nil)
}

func TestHandler_Create(t *testing.T) {
	c, uh := setupEcho(t, "/api/v1/create", http.MethodPost)

	err := uh.Create(c)
	require.Equal(t, c.JSON(http.StatusOK, EventCreate), err)
}
