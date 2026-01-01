package models

// ErrorResponse представляет единую ошибку в ответе API.
// Структурированное представление ошибки с кодом, сообщением и деталями.
type ErrorResponse struct {
	// Code - Код ошибки (машиночитаемый идентификатор).
	// Используется для программной обработки ошибок и локализации.
	// Примеры: "INVALID_REQUEST", "DOCUMENT_NOT_FOUND", "UNAUTHORIZED"
	Code string `json:"code"`

	// Message - Человеческое сообщение об ошибке на русском языке.
	// Короткое и понятное описание проблемы.
	// Примеры: "Документ не найден", "Неверный формат даты"
	Message string `json:"message"`

	// Details - Дополнительные детали об ошибке (опционально).
	// Может содержать более подробное объяснение, стек вызовов или
	// другую информацию для отладки (обычно только в режиме разработки).
	// Может быть пусто для ошибок в production.
	Details string `json:"details,omitempty"`
}

// ErrorCodes определяет константы кодов ошибок ESF API.
// Используется для унификации и программной обработки ошибок.
type ErrorCodes struct {
	// ========== ОШИБКИ ВАЛИДАЦИИ ==========

	// InvalidRequest - Запрос имеет неверный формат или структуру.
	// Примеры: неверный JSON, отсутствуют обязательные поля
	InvalidRequest string

	// MissingRequiredField - Отсутствует обязательное поле в запросе.
	// Указывает на то, что некоторые необходимые данные не предоставлены.
	MissingRequiredField string

	// InvalidFieldValue - Значение поля не соответствует ожидаемому формату.
	// Примеры: неверная дата, отрицательная сумма, недопустимый код
	InvalidFieldValue string

	// InvalidDocumentType - Неверный тип документа (не Реализация/Приобретение/Корректировочная).
	InvalidDocumentType string

	// InvalidDocumentStatus - Операция недопустима для текущего статуса документа.
	// Примеры: попытка отредактировать подписанный документ
	InvalidDocumentStatus string

	// InvalidParty - Неверные реквизиты юридического лица или контрагента.
	// Примеры: несуществующий ИНН, пустое наименование
	InvalidParty string

	// InvalidBankAccount - Неверные реквизиты банковского счета.
	// Примеры: неверный номер счета, несуществующий банк
	InvalidBankAccount string

	// ========== ОШИБКИ АВТОРИЗАЦИИ ==========

	// Unauthorized - Пользователь не авторизован (отсутствует токен).
	// HTTP 401. Требуется авторизация для доступа к ресурсу.
	Unauthorized string

	// Forbidden - Пользователь авторизован, но не имеет доступа к ресурсу.
	// HTTP 403. Недостаточно прав для выполнения операции.
	Forbidden string

	// InvalidToken - Токен авторизации неверен или истек.
	// Требуется новая авторизация.
	InvalidToken string

	// MissingAuthHeaders - Отсутствуют требуемые заголовки авторизации.
	// Примеры: отсутствует Authorization, X-Road-Client и т.д.
	MissingAuthHeaders string

	// ========== ОШИБКИ ДОКУМЕНТОВ ==========

	// DocumentNotFound - Документ с указанным ID/UUID не найден в системе.
	DocumentNotFound string

	// DocumentAlreadyExists - Документ с таким номером уже существует.
	// Нарушение уникальности номера в пределах периода.
	DocumentAlreadyExists string

	// DocumentCannotBeEdited - Документ находится в статусе, не позволяющем редактирование.
	// Примеры: подписанный, отозванный, отклоненный документ
	DocumentCannotBeEdited string

	// DocumentCannotBeDeleted - Документ не может быть удален в текущем состоянии.
	DocumentCannotBeDeleted string

	// DocumentCannotBeSigned - Документ не может быть подписан (неверный статус).
	// Примеры: черновик без деталей, уже подписанный документ
	DocumentCannotBeSigned string

	// DocumentCannotBeRevoked - Документ не может быть отозван.
	// Возможно, он уже подписан получателем.
	DocumentCannotBeRevoked string

	// InvalidDocumentUUID - UUID документа имеет неверный формат.
	InvalidDocumentUUID string

	// DuplicateDocumentNumber - Номер документа уже используется в другом счете.
	DuplicateDocumentNumber string

	// ========== ОШИБКИ БАНКОВ И СЧЕТОВ ==========

	// BankNotFound - Банк с указанным БИК/ID не найден в системе.
	BankNotFound string

	// BankAccountNotFound - Банковский счет с указанным ID не найден.
	BankAccountNotFound string

	// InvalidBankData - Данные о банке имеют неверный формат или неполные.
	InvalidBankData string

	// ========== ОШИБКИ СПРАВОЧНИКОВ ==========

	// CatalogNotFound - Элемент справочника с указанным ID не найден.
	CatalogNotFound string

	// InvalidCatalogData - Данные справочника имеют неверный формат или неполные.
	InvalidCatalogData string

	// DuplicateCatalogEntry - Элемент справочника с таким номером уже существует.
	DuplicateCatalogEntry string

	// ========== ОШИБКИ ИНОСТРАННЫХ КОМПАНИЙ ==========

	// ForeignCompanyNotFound - Иностранная компания с указанным ID не найдена.
	ForeignCompanyNotFound string

	// InvalidForeignCompanyData - Данные иностранной компании имеют неверный формат.
	InvalidForeignCompanyData string

	// ========== ОШИБКИ ВНУТРЕННИЕ ==========

	// InternalServerError - Внутренняя ошибка сервера (неожиданная ошибка).
	// HTTP 500. Указывает на ошибку в коде сервера.
	InternalServerError string

	// DatabaseError - Ошибка при работе с базой данных.
	// Примеры: потеря соединения, блокировка записи
	DatabaseError string

	// IntegrationError - Ошибка при интеграции с внешней системой (СМЭВ, ТУНДУК).
	// Внешняя система недоступна или вернула ошибку.
	IntegrationError string

	// ServiceUnavailable - Сервис недоступен (на обслуживании или перегружен).
	// HTTP 503.
	ServiceUnavailable string
}

// StatusCodes определяет константы кодов статусов документов в ESF API.
// Каждый документ проходит через разные статусы в своем жизненном цикле.
type StatusCodes struct {
	// Draft - Документ в статусе черновика.
	// Создан, но еще не отправлен. Может быть отредактирован.
	Draft string

	// Submitted - Документ отправлен получателю.
	// Ожидает принятия или отклонения.
	Submitted string

	// Approved - Документ одобрен получателем.
	// Контрагент согласен с содержимым документа.
	Approved string

	// Rejected - Документ отклонен получателем.
	// Требует исправления и переотправки.
	Rejected string

	// Signed - Документ подписан (электронной подписью).
	// Имеет юридическую силу. Финальный статус для валидного документа.
	Signed string

	// Revoked - Документ отозван отправителем.
	// Больше не действителен, исключен из оборота.
	Revoked string

	// AwaitingApproval - Документ ожидает одобрения/подписания.
	// Промежуточный статус перед финализацией.
	AwaitingApproval string

	// Error - Ошибка при обработке документа.
	// Возникла неожиданная ошибка, документ не может быть обработан.
	Error string
}

// CorrectionReasons определяет константы причин корректировки документов.
// Используется при создании корректировочных счетов для указания причины изменения.
type CorrectionReasons struct {
	// PriceError - Исправление ошибки в цене товара/услуги.
	// Цена была указана неверно, требуется корректировка.
	PriceError string

	// QuantityError - Исправление ошибки в количестве товара.
	// Количество было указано неверно.
	QuantityError string

	// TaxError - Исправление ошибки в расчете налогов (НДС, НСП).
	// Налоги были рассчитаны неверно.
	TaxError string

	// DetailsError - Исправление ошибки в описании или характеристиках товара.
	// Наименование, код ТНВЭД или другие реквизиты указаны неверно.
	DetailsError string

	// DataError - Исправление ошибки в общих данных счета.
	// Ошибка в датах, реквизитах сторон, валюте и т.д.
	DataError string

	// Other - Другая причина корректировки (не указанная выше).
	// Используется для прочих случаев изменения документа.
	Other string
}
