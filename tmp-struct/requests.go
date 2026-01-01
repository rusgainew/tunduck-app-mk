package models

// CreateInvoiceRequest представляет запрос на создание счета-фактуры.
// Используется для POST запроса /api/command/invoice/{type}/create
//
// Содержит все необходимые поля для создания нового счета в системе:
//   - Основные реквизиты (номер, дата, сумма)
//   - Информацию о сторонах (поставщик и покупатель)
//   - Налоговую информацию (валюта, НДС)
//   - Детали товаров/услуг (строки счета)
//
// Все обязательные поля должны быть заполнены перед отправкой запроса.
type CreateInvoiceRequest struct {
	// InvoiceNumber - Номер счета-фактуры, присваиваемый поставщиком.
	// Обязательное поле. Должно быть уникальным для организации.
	InvoiceNumber string `json:"invoiceNumber"`

	// InvoiceDate - Дата составления счета (формат ISO 8601: YYYY-MM-DD).
	// Дата подписания документа. Может быть nil, если не установлена.
	InvoiceDate *string `json:"invoiceDate"`

	// DeliveryDate - Дата поставки товаров/оказания услуг (формат ISO 8601).
	// Может отличаться от InvoiceDate. Может быть nil.
	DeliveryDate *string `json:"deliveryDate"`

	// TotalAmount - Общая сумма счета (включая все налоги).
	// Рассчитывается как сумма всех деталей.
	TotalAmount float64 `json:"totalAmount"`

	// IsResident - Флаг резидентства поставщика.
	// true: резидент Кыргызской Республики, false: нерезидент.
	IsResident bool `json:"isResident"`

	// Note - Примечания/комментарии к счету.
	// Дополнительная информация для получателя.
	Note string `json:"note"`

	// LegalPerson - Реквизиты поставщика (юридического лица).
	// Обязательное поле. Содержит ИНН и наименование.
	LegalPerson *Party `json:"legalPerson"`

	// Contractor - Реквизиты покупателя/получателя услуг.
	// Обязательное поле. Содержит ИНН и наименование контрагента.
	Contractor *Party `json:"contractor"`

	// Currency - Валюта, в которой выражены все суммы счета.
	// Обязательное поле. Примеры: KGS (Сом), USD (Доллар).
	Currency *ReferenceItem `json:"currency"`

	// PaymentType - Вид оплаты (Наличная, Безналичная, Кредит и т.д.).
	// Обязательное поле. Определяет способ расчета.
	PaymentType *ReferenceItem `json:"paymentType"`

	// DeliveryType - Вид доставки (Самовывоз, Транспортом и т.д.).
	// Обязательное поле. Определяет способ передачи товаров.
	DeliveryType *ReferenceItem `json:"deliveryType"`

	// VATTaxType - Вид налога НДС и его ставка.
	// Обязательное поле. Определяет, какая ставка НДС применяется.
	VATTaxType *VATTaxType `json:"vatTaxType"`

	// Details - Массив строк счета-фактуры (товары/услуги).
	// Обязательное поле. Минимум одна деталь должна быть в счете.
	// Каждая деталь содержит информацию о товаре/услуге с количеством и ценой.
	Details []*InvoiceDetail `json:"details"`
}

// UpdateInvoiceRequest представляет запрос на обновление/редактирование счета-фактуры.
// Используется для PUT запроса /api/command/invoice/{type}/edit/{id}
//
// Позволяет изменять реквизиты уже созданного счета:
//   - Номер и дату счета
//   - Даты доставки
//   - Суммы и налоги
//   - Информацию о контрагентах (если разрешено)
//   - Детали товаров/услуг
//
// Примечание: Не все поля могут быть отредактированы после подписания счета.
type UpdateInvoiceRequest struct {
	// ID - Идентификатор счета в системе (обязательное поле).
	// Определяет, какой счет редактируется.
	ID int64 `json:"id"`

	// InvoiceNumber - Номер счета-фактуры (может быть изменен).
	InvoiceNumber string `json:"invoiceNumber"`

	// InvoiceDate - Дата составления счета (может быть изменена).
	InvoiceDate *string `json:"invoiceDate"`

	// DeliveryDate - Дата поставки товаров (может быть изменена).
	DeliveryDate *string `json:"deliveryDate"`

	// TotalAmount - Общая сумма счета (пересчитывается автоматически).
	TotalAmount float64 `json:"totalAmount"`

	// IsResident - Флаг резидентства (может быть изменен).
	IsResident bool `json:"isResident"`

	// Note - Примечания к счету (может быть изменено).
	Note string `json:"note"`

	// Currency - Валюта счета (может быть изменена в черновике).
	Currency *ReferenceItem `json:"currency"`

	// PaymentType - Вид оплаты (может быть изменен).
	PaymentType *ReferenceItem `json:"paymentType"`

	// DeliveryType - Вид доставки (может быть изменен).
	DeliveryType *ReferenceItem `json:"deliveryType"`

	// VATTaxType - Вид налога НДС (может быть изменен).
	VATTaxType *VATTaxType `json:"vatTaxType"`

	// Details - Обновленный список деталей счета.
	// Может быть добавлены, удалены или изменены строки.
	Details []*InvoiceDetail `json:"details"`
}

// AcceptOrRejectRequest представляет запрос на принятие или отклонение счета-фактуры.
// Используется для POST запроса /api/command/invoice/{type}/accept-or-reject
//
// Позволяет получателю счета (покупателю) принять или отклонить счет.
// Это важный шаг в жизненном цикле документа перед его подписанием.
type AcceptOrRejectRequest struct {
	// DocumentUUID - Уникальный идентификатор счета-фактуры.
	// Обязательное поле. Определяет, какой счет принимается/отклоняется.
	DocumentUUID string `json:"documentUuid"`

	// Accept - Флаг принятия/отклонения.
	// true: Счет принят (одобрен)
	// false: Счет отклонен (требует исправления)
	Accept bool `json:"accept"`

	// RejectReason - Причина отклонения счета.
	// Заполняется, если Accept=false.
	// Содержит текстовое объяснение причины отклонения.
	// Может быть пусто, если счет принят (Accept=true).
	RejectReason string `json:"rejectReason"`
}

// RevokeRequest представляет запрос на отзыв (отмену) документа.
// Используется для POST запроса /api/command/invoice/{type}/revoke
//
// Позволяет отправителю счета отозвать уже отправленный документ
// (если получатель еще не подписал его).
type RevokeRequest struct {
	// DocumentUUID - Уникальный идентификатор счета для отзыва.
	// Обязательное поле. Определяет, какой документ отзывается.
	// После отзыва документ будет помечен как "Отозван" и больше не может использоваться.
	DocumentUUID string `json:"documentUuid"`
}

// SignRequest представляет запрос на подписание документа.
// Используется для POST запроса /api/command/invoice/{type}/sign
//
// Содержит электронную подпись (ЭЦП) документа.
// Подпись подтверждает подлинность и авторство документа.
type SignRequest struct {
	// DocumentUUID - Уникальный идентификатор счета для подписания.
	// Обязательное поле. Определяет, какой документ подписывается.
	DocumentUUID string `json:"documentUuid"`

	// Signature - Электронная цифровая подпись (ЭЦП) в формате base64.
	// Обязательное поле. Содержит цифровую подпись в кодировке base64.
	// Подпись создается с помощью средств ЭЦП (сертификат и ключ подписи).
	Signature string `json:"signature"`
}

// CreateBankAccountRequest представляет запрос на создание банковского счета.
// Используется для POST запроса /api/command/legal-person/bank-account/create
//
// Позволяет добавить новый банковский счет организации.
// Счет может использоваться в счетах-фактурах для указания реквизитов платежа.
type CreateBankAccountRequest struct {
	// AccountName - Название/описание счета в системе организации.
	// Например: "Основной расчетный счет", "Валютный счет USD".
	AccountName string `json:"accountName"`

	// BankAccount - Номер банковского счета.
	// Полный номер счета в банке (обычно 20-24 цифры).
	BankAccount string `json:"bankAccount"`

	// ContractorTIN - ИНН организации-владельца счета.
	// Обычно совпадает с ИНН текущей организации.
	ContractorTIN string `json:"contractorTin"`

	// Currency - Валюта счета (KGS, USD, EUR и т.д.).
	// Определяет, в какой валюте открыт счет.
	Currency *ReferenceItem `json:"currency"`

	// Bank - Информация о банке, в котором открыт счет.
	// Содержит название, БИК и адрес банка.
	Bank *Bank `json:"bank"`

	// IsResident - Флаг резидентства владельца счета.
	// true: счет резидента КР, false: счет нерезидента.
	IsResident bool `json:"isResident"`
}

// UpdateBankAccountRequest представляет запрос на обновление банковского счета.
// Используется для PUT запроса /api/command/legal-person/bank-account/edit/{id}
//
// Позволяет изменять реквизиты существующего банковского счета:
//   - Название счета
//   - Валюту
//   - Информацию о банке
//   - Статус резидентства
type UpdateBankAccountRequest struct {
	// ID - Идентификатор счета для редактирования.
	// Обязательное поле. Определяет, какой счет обновляется.
	ID int64 `json:"id"`

	// AccountName - Новое название счета.
	AccountName string `json:"accountName"`

	// BankAccount - Номер банковского счета (может быть изменен).
	BankAccount string `json:"bankAccount"`

	// Currency - Валюта счета (может быть изменена).
	Currency *ReferenceItem `json:"currency"`

	// Bank - Информация о банке (может быть обновлена).
	Bank *Bank `json:"bank"`

	// IsResident - Флаг резидентства (может быть изменен).
	IsResident bool `json:"isResident"`
}

// CreateCatalogRequest представляет запрос на создание элемента справочника.
// Используется для POST запроса /api/command/legal-person/catalog/create
//
// Добавляет новый товар или услугу в справочник организации.
// Справочник используется для быстрого заполнения строк в счетах.
type CreateCatalogRequest struct {
	// Name - Наименование товара/услуги.
	// Обязательное поле. Полное описание товара для использования в счетах.
	Name string `json:"name"`

	// Number - Номер товара в справочнике организации.
	// Уникальный идентификатор в системе организации (например, SKU).
	Number string `json:"number"`

	// TNVEDCode - Код ТНВЭД для таможенной классификации.
	// Может быть пусто для товаров, не участвующих в импорте/экспорте.
	TNVEDCode string `json:"tnvedCode"`

	// GKEDCode - Код ГКЭД (классификация вида деятельности).
	// Определяет категорию товара/услуги.
	GKEDCode string `json:"gkedCode"`

	// UnitClassification - Единица измерения товара.
	// Обязательное поле. Определяет, в каких единицах измеряется товар.
	UnitClassification *ReferenceItem `json:"unitClassification"`

	// CatalogType - Тип элемента (товар или услуга).
	// Обязательное поле. Различает материальные товары от услуг.
	CatalogType *ReferenceItem `json:"catalogType"`
}

// UpdateCatalogRequest представляет запрос на обновление элемента справочника.
// Используется для PUT запроса /api/command/legal-person/catalog/edit/{id}
//
// Позволяет изменять информацию о товаре/услуге в справочнике.
type UpdateCatalogRequest struct {
	// ID - Идентификатор элемента справочника для редактирования.
	// Обязательное поле. Определяет, какой товар/услугу обновлять.
	ID int64 `json:"id"`

	// Name - Новое наименование товара/услуги.
	Name string `json:"name"`

	// Number - Новый номер товара.
	Number string `json:"number"`

	// TNVEDCode - Обновленный код ТНВЭД.
	TNVEDCode string `json:"tnvedCode"`

	// GKEDCode - Обновленный код ГКЭД.
	GKEDCode string `json:"gkedCode"`

	// UnitClassification - Обновленная единица измерения.
	UnitClassification *ReferenceItem `json:"unitClassification"`

	// CatalogType - Обновленный тип (товар/услуга).
	CatalogType *ReferenceItem `json:"catalogType"`
}

// CreateForeignCompanyRequest представляет запрос на создание иностранной компании.
// Используется для POST запроса /api/command/foreign-company/create
//
// Позволяет добавить в систему новую иностранную компанию-контрагента.
// Используется для операций с иностранными партнерами (экспорт, импорт, услуги).
type CreateForeignCompanyRequest struct {
	// PIN - Идентификационный номер иностранной компании в ее стране.
	// Может быть Tax ID, VAT Number или другой номер в зависимости от страны.
	PIN string `json:"pin"`

	// FullName - Полное наименование иностранной компании.
	// На языке страны происхождения или английском языке.
	FullName string `json:"fullName"`
}

// UpdateForeignCompanyRequest представляет запрос на обновление иностранной компании.
// Используется для PUT запроса /api/command/foreign-company/edit/{id}
//
// Позволяет изменять информацию об иностранной компании-контрагенте.
type UpdateForeignCompanyRequest struct {
	// ID - Идентификатор иностранной компании для редактирования.
	// Обязательное поле. Определяет, какую компанию обновлять.
	ID int64 `json:"id"`

	// PIN - Обновленный идентификационный номер компании.
	PIN string `json:"pin"`

	// FullName - Обновленное наименование компании.
	FullName string `json:"fullName"`
}
