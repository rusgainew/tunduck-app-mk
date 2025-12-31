package models

import (
	"time"

	"github.com/google/uuid"
)

// Create document
// Запрос:
// METHOD: POST
// PATH: /api/command/invoice/create
type EsfCreateDocumentRequest struct {
	// false Наименование иностранца или Наименование на иностранном языке
	ForeignName string `json:"foreignName"`
	// false Отправить от имени филиала
	IsBranchDataSent bool `json:"isBranchDataSent"`
	// false Цена без налогов
	IsPriceWithoutTaxes bool `json:"isPriceWithoutTaxes"`
	// false ИНН филиала
	AffiliateTin string `json:"affiliateTin"`
	// false Отраслевые
	IsIndustry bool `json:"isIndustry"`
	// false Номер учетной системы
	OwnedCrmReceiptCode string `json:"ownedCrmReceiptCode"`
	// false Код вида операции
	OperationTypeCode string `json:"operationTypeCode"`
	// false Дата поставки
	DeliveryDate *time.Time `json:"deliveryDate"`
	// false Код типа поставки
	DeliveryTypeCode string `json:"deliveryTypeCode"`
	// false Субъект Кыргызской Республики
	IsResident bool `json:"isResident"`
	// false ИНН покупателя
	ContractorTin string `json:"contractorTin"`
	// false Номер банковского счета поставщика
	SupplierBankAccount string `json:"supplierBankAccount"`
	// false Номер банковского счета покупателя
	ContractorBankAccount string `json:"contractorBankAccount"`
	// false Код валюты
	CurrencyCode string `json:"currencyCode"`
	// false Код страны
	CountryCode string `json:"countryCode"`
	// false Курс валюты к сому
	CurrencyRate *float64 `json:"currencyRate"`
	// false Общая стоимость в валюте
	TotalCurrencyValue *float64 `json:"totalCurrencyValue"`
	// false Общая стоимость в валюте без налогов
	TotalCurrencyValueWithoutTaxes *float64 `json:"totalCurrencyValueWithoutTaxes"`
	// false Номер договора на поставку
	SupplyContractNumber string `json:"supplyContractNumber"`
	// false Дата договора на поставку
	ContractStartDate *time.Time `json:"contractStartDate"`
	// false Дата окончания договора на поставку Комментарий
	Comment string `json:"comment"`
	// false Код способа доставки
	DeliveryCode string `json:"deliveryCode"`
	// false Код формы оплаты
	PaymentCode string `json:"paymentCode"`
	// false Код ставки НДС
	TaxRateVATCode string `json:"taxRateVATCode"`
	// false Товары и услуги
	CatalogEntries []EsfEntriesModel `json:"catalogEntries"`
	// false Начальные остатки,сальдо на начало периода
	OpeningBalances *float64 `json:"openingBalances"`
	// false Начисленные взносы
	AssessedContributionsAmount *float64 `json:"assessedContributionsAmount"`
	// false Поступления,оплачено
	PaidAmount *float64 `json:"paidAmount"`
	// false Штрафы
	PenaltiesAmount *float64 `json:"penaltiesAmount"`
	// false пени
	FinesAmount *float64 `json:"finesAmount"`
	// false Конечные остатки,сальдо на конец периода
	ClosingBalances *float64 `json:"closingBalances"`
	// false Сумма к оплате
	AmountToBePaid *float64 `json:"amountToBePaid"`
	// false Лицевой счет
	PersonalAccountNumber string `json:"personalAccountNumber"`
}
type EsfCreateDocumentResponse struct {
	ResponseId   string `json:"responseId"`
	DocumentUuid string `json:"documentUuid"`
}

// Edit document
// Описание:
// Редактирование ЭСФ
// Запрос:
// METHOD: PUT
// PATH: /api/command/invoice/edit/{id}

type EsfEditDocumentRequest struct {
	ID uuid.UUID `json:"id" valid:"required"`
	EsfCreateDocumentRequest
}
