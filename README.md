# Go multiplexer example

Это репозиторий с минипроектом, написанном на go в качестве тестового задания.

## Сама задача в неизменном виде:

Тестовое задание HTTP-мультиплексор:
приложение представляет собой http-сервер с одним хендлером, хендлер на вход
получает POST-запрос со списком url в json-формате сервер запрашивает данные по
всем этим url и возвращает результат клиенту в json-формате если в процессе
обработки хотя бы одного из url получена ошибка, обработка всего списка
прекращается и клиенту возвращается текстовая ошибка Ограничения: для реализации
задачи следует использовать Go 1.13 или выше использовать можно только
компоненты стандартной библиотеки Go сервер не принимает запрос если количество
url в в нем больше 20 сервер не обслуживает больше чем 100 одновременных
входящих http-запросов для каждого входящего запроса должно быть не больше 4
одновременных исходящих таймаут на запрос одного url - секунда обработка запроса
может быть отменена клиентом в любой момент, это должно повлечь за собой
остановку всех операций связанных с этим запросом сервис должен поддерживать
'graceful shutdown'.

## Как запускать

Как обычно: задать env переменные (можно руками, можно через
`export $(cat .env | xargs)`), далее go run ./src

## комментарии к решению

* В задаче не сказано, какие именно данные будут переданы от url'a, текстовые
  или бинарные. я предполагаю, что текстовые, но возможно, что необходимо
  сделать дополнительную проверку на content-type и при необходимости,
  кодировать ответ в base64.

* Мы не говорили о том, что нам нужно возможно необходимо помимо самих данных
  еще и отдавать статус код, определнные хедеры, и пр. (так же может нам нужно
  помимо самих урлов указывать какие-то креденшелы, чтобы получить к ним доступ)
  Это я продумал, в таком случае мы можем изменить структуры RequestParameters и
  RequestResponse как нам это действительно необходимо.

* в задаче говорилось о том, что нельзя использовать нестандартные пакеты, но
  в идеале нужно тестировать апи, для этого бы хорошо подошел
  [resty](https://github.com/go-resty/resty), но раз нельзя то нельзя ¯\_(ツ)_/¯

* возможно есть вопрос, почему роутер и оператор решены через интерфейсы:
  время выполнения на это усложнение не сильно много тратится, однако зато мы
  уменьшаем наш техдолг и при необходимости можем оперативно доработать
  реализацию с сохранением обратной совместимости

* т.к. задача "скачать данные с кучи урлов" все таки тяжеловесная, возможно
  правильнее было бы спроектировать апи так, что бы при запросе выдавать клиенту
  "тикет" (уникальный иднетификатор задачи) и предлагать ему когда необходимо
  получать статус выполнения задачи. Это 100% необходимо, если клиент попросил
  загрузить данные с 100+ эндпоинтов. Этот момент не столько комментарий,
  сколько мысли вслух по поводу "чего здесь не хватает, если бы мы делали
  полноценный инструмент"
