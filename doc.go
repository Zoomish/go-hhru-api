// Пакет hhru — типизированный HTTP-клиент для API HeadHunter (api.hh.ru).
//
// Вызовите [New] с заполненным [Options].HHUserAgent: имя приложения и контактный e-mail,
// как требует HH. Опциональный [Options.TokenSource] добавляет Bearer ко всем запросам.
// Для токена приложения используйте [ExchangeClientCredentials] и [AccessToken].
// Для пользовательского OAuth с refresh_token — [NewRefreshingTokenSource] или
// [NewRefreshingTokenSourceWithOptions] с [Clock] в тестах.
//
// Сгенерированные подклиенты (по разбиению публичной OpenAPI) лежат в пакетах:
//   - [github.com/Zoomish/go-hhru-api/gen/employer] — API работодателя;
//   - [github.com/Zoomish/go-hhru-api/gen/applicant] — авторизация и сценарии соискателя;
//   - [github.com/Zoomish/go-hhru-api/gen/public] — публичные справочники и подсказки;
//   - [github.com/Zoomish/go-hhru-api/gen/app] — эндпоинты в контексте приложения.
//
// На [Client] они доступны как Employer, Applicant, Public и App.
//
// Надёжность и наблюдаемость через [Options]: [Options.MaxRetries], пределы паузы
// [Options.RetryBackoffMin] и [Options.RetryBackoffMax], ограничение частоты
// [Options.MaxRequestsPerSecond], хуки [Options.RequestHook] и [Options.ResponseHook].
// Разбор JSON-ошибок API — [ParseAPIError].
//
// Запускаемые примеры — в каталоге examples/ (см. README репозитория).
// Живые HTTP-тесты — с тегом сборки "integration" в пакете integration/.
//
// Официальная документация OpenAPI: https://api.hh.ru/openapi/redoc
package hhru
