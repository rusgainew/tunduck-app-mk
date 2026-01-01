package models

// Route представляет API маршрут (endpoint) системы.
// Используется для описания доступных API endpoints и их параметров.
type Route struct {
	// Method - HTTP метод (GET, POST, PUT, DELETE, PATCH)
	Method string

	// Path - URL путь маршрута (например, "/api/query/invoice/realization/document-by-exchange-code")
	Path string

	// Name - Человеческое имя маршрута (например, "Get Realization Invoice by UUID")
	Name string

	// Description - Полное описание, что делает маршрут
	Description string

	// Category - Категория API (Query, Command, Dictionary, Auth)
	Category string

	// Authenticated - Требуется ли авторизация (обычно true для всех, кроме /api/account/auth)
	Authenticated bool

	// RequiredParams - Обязательные параметры запроса
	RequiredParams []string

	// OptionalParams - Опциональные параметры запроса
	OptionalParams []string

	// Response - Тип возвращаемого ответа (например, "Invoice", "[]*Invoice", "BankAccount")
	Response string
}

// APIRoutes содержит все доступные маршруты ESF API
var APIRoutes = []Route{
	// ========== QUERY API - СЧЕТА-ФАКТУРЫ ==========

	{
		Method:         "GET",
		Path:           "/api/query/invoice/realization/document-by-exchange-code",
		Name:           "Get Realization Invoice by UUID",
		Description:    "Получить счет-фактуру на реализацию по уникальному идентификатору (UUID)",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},
	{
		Method:         "GET",
		Path:           "/api/query/invoice/income/document-by-exchange-code",
		Name:           "Get Income Invoice by UUID",
		Description:    "Получить счет-фактуру на приобретение по уникальному идентификатору (UUID)",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},
	{
		Method:         "GET",
		Path:           "/api/query/invoice/correction/document-by-exchange-code",
		Name:           "Get Correction Invoice by UUID",
		Description:    "Получить корректировочный счет-фактуру по уникальному идентификатору (UUID)",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},
	{
		Method:         "GET",
		Path:           "/api/query/invoice/realization/all-documents",
		Name:           "Get All Realization Invoices",
		Description:    "Получить список счетов-фактур на реализацию с поддержкой фильтрации, сортировки и пагинации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"page", "size", "sort", "documentUuid", "invoiceNumber", "status"},
		Response:       "*InvoiceListResponse",
	},
	{
		Method:         "GET",
		Path:           "/api/query/invoice/income/all-documents",
		Name:           "Get All Income Invoices",
		Description:    "Получить список счетов-фактур на приобретение с поддержкой фильтрации, сортировки и пагинации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"page", "size", "sort", "documentUuid", "invoiceNumber", "status"},
		Response:       "*InvoiceListResponse",
	},
	{
		Method:         "GET",
		Path:           "/api/query/invoice/correction/all-documents",
		Name:           "Get All Correction Invoices",
		Description:    "Получить список корректировочных счетов-фактур с поддержкой фильтрации, сортировки и пагинации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"page", "size", "sort", "documentUuid", "invoiceNumber", "status"},
		Response:       "*InvoiceListResponse",
	},

	// ========== QUERY API - ДЕТАЛИ СЧЕТОВ ==========

	{
		Method:         "GET",
		Path:           "/api/query/detail/details-by-exchange-code",
		Name:           "Get Invoice Details (Paginated)",
		Description:    "Получить детали (строки) счета-фактуры с поддержкой пагинации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"invoiceUuid"},
		OptionalParams: []string{"page", "size", "sort"},
		Response:       "*DetailListResponse",
	},
	{
		Method:         "GET",
		Path:           "/api/query/detail/v2/details-by-exchange-code",
		Name:           "Get Invoice Details (Full List)",
		Description:    "Получить все детали счета-фактуры без пагинации (полный список)",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"invoiceUuid"},
		OptionalParams: []string{},
		Response:       "*DetailListResponse",
	},

	// ========== QUERY API - БАНКОВСКИЕ СЧЕТА ==========

	{
		Method:         "GET",
		Path:           "/api/query/legal-person/bank-accounts",
		Name:           "Get Bank Accounts",
		Description:    "Получить список всех банковских счетов организации с поддержкой пагинации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"page", "size", "sort"},
		Response:       "*BankAccountListResponse",
	},
	{
		Method:         "GET",
		Path:           "/api/query/legal-person/bank-account/{id}",
		Name:           "Get Bank Account by ID",
		Description:    "Получить информацию о конкретном банковском счете по его ID",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"id", "tin"},
		OptionalParams: []string{},
		Response:       "*BankAccount",
	},

	// ========== QUERY API - СПРАВОЧНИКИ ==========

	{
		Method:         "GET",
		Path:           "/api/query/legal-person/catalogs",
		Name:           "Get Catalogs",
		Description:    "Получить список справочников товаров/услуг организации",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"page", "size", "sort", "catalogType"},
		Response:       "*CatalogListResponse",
	},
	{
		Method:         "GET",
		Path:           "/api/query/legal-person/catalog/{id}",
		Name:           "Get Catalog by ID",
		Description:    "Получить информацию об элементе справочника по его ID",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"id", "tin"},
		OptionalParams: []string{},
		Response:       "*Catalog",
	},

	// ========== QUERY API - ИНОСТРАННЫЕ КОМПАНИИ ==========

	{
		Method:         "GET",
		Path:           "/api/query/legal-person/foreign-companies",
		Name:           "Get Foreign Companies",
		Description:    "Получить список иностранных компаний-контрагентов с фильтрацией и пагинацией",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"foreignTin", "fullName", "sort", "page", "size"},
		Response:       "[]*ForeignCompany",
	},
	{
		Method:         "GET",
		Path:           "/api/query/legal-person/foreign-company/{id}",
		Name:           "Get Foreign Company by ID",
		Description:    "Получить информацию об иностранной компании по ID",
		Category:       "Query",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{},
		Response:       "*ForeignCompany",
	},

	// ========== COMMAND API - СЧЕТА-ФАКТУРЫ ==========

	{
		Method:         "POST",
		Path:           "/api/command/invoice/create",
		Name:           "Create Invoice",
		Description:    "Создать новый счет-фактуру (Реализация, Приобретение или Корректировочная)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"invoiceNumber", "totalAmount", "legalPerson", "contractor", "details"},
		OptionalParams: []string{"invoiceDate", "deliveryDate", "note", "currency", "paymentType", "deliveryType"},
		Response:       "*Invoice",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/invoice/edit/{id}",
		Name:           "Edit Invoice",
		Description:    "Отредактировать существующий счет-фактуру (только в статусе черновика)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{"invoiceNumber", "invoiceDate", "deliveryDate", "totalAmount", "note", "details"},
		Response:       "*Invoice",
	},
	{
		Method:         "DELETE",
		Path:           "/api/command/invoice/delete/{id}",
		Name:           "Delete Invoice",
		Description:    "Удалить счет-фактуру (только в статусе черновика)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{},
		Response:       "map[string]interface{}",
	},
	{
		Method:         "POST",
		Path:           "/api/command/invoice/accept-or-reject",
		Name:           "Accept or Reject Invoice",
		Description:    "Принять или отклонить счет-фактуру (для получателя)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid", "accept"},
		OptionalParams: []string{"rejectReason"},
		Response:       "*Invoice",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/invoice/revoke",
		Name:           "Revoke Invoice",
		Description:    "Отозвать отправленный счет-фактуру (для отправителя, если еще не подписан)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},
	{
		Method:         "POST",
		Path:           "/api/command/invoice/correction/create/{primaryId}",
		Name:           "Create Correction Invoice",
		Description:    "Создать корректировочный счет для существующего счета-фактуры",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"primaryId", "details"},
		OptionalParams: []string{"correctionReason", "note"},
		Response:       "*Invoice",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/invoice/correction/edit/{id}",
		Name:           "Edit Correction Invoice",
		Description:    "Отредактировать корректировочный счет (только в статусе черновика)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{"details", "note"},
		Response:       "*Invoice",
	},

	// ========== COMMAND API - ПОДПИСАНИЕ ==========

	{
		Method:         "PUT",
		Path:           "/api/command/invoice/sign",
		Name:           "Sign Invoice",
		Description:    "Подписать счет-фактуру электронной подписью (ЭЦП)",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid", "signature"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/invoice/sign-with-infocom",
		Name:           "Sign Invoice with Infocom",
		Description:    "Подписать счет-фактуру с использованием сервиса Инфокома",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"documentUuid"},
		OptionalParams: []string{},
		Response:       "*Invoice",
	},

	// ========== COMMAND API - БАНКОВСКИЕ СЧЕТА ==========

	{
		Method:         "POST",
		Path:           "/api/command/legal-person/bank/create",
		Name:           "Create Bank Account",
		Description:    "Создать новый банковский счет для организации",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"accountName", "bankAccount", "currency", "bank"},
		OptionalParams: []string{"contractorTin", "isResident"},
		Response:       "*BankAccount",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/legal-person/bank/edit/{id}",
		Name:           "Edit Bank Account",
		Description:    "Отредактировать банковский счет",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{"accountName", "bankAccount", "currency", "bank", "isResident"},
		Response:       "*BankAccount",
	},
	{
		Method:         "DELETE",
		Path:           "/api/command/legal-person/bank/delete/{id}",
		Name:           "Delete Bank Account",
		Description:    "Удалить банковский счет",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{},
		Response:       "map[string]interface{}",
	},

	// ========== COMMAND API - СПРАВОЧНИКИ ==========

	{
		Method:         "POST",
		Path:           "/api/command/legal-person/catalog/create",
		Name:           "Create Catalog",
		Description:    "Создать новый элемент справочника товаров/услуг",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"name", "number", "unitClassification", "catalogType"},
		OptionalParams: []string{"tnvedCode", "gkedCode"},
		Response:       "*Catalog",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/legal-person/catalog/edit/{id}",
		Name:           "Edit Catalog",
		Description:    "Отредактировать элемент справочника",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{"name", "number", "tnvedCode", "gkedCode", "unitClassification", "catalogType"},
		Response:       "*Catalog",
	},
	{
		Method:         "DELETE",
		Path:           "/api/command/legal-person/catalog/delete/{id}",
		Name:           "Delete Catalog",
		Description:    "Удалить элемент справочника",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{},
		Response:       "map[string]interface{}",
	},

	// ========== COMMAND API - ИНОСТРАННЫЕ КОМПАНИИ ==========

	{
		Method:         "POST",
		Path:           "/api/command/legal-person/foreign-company/create",
		Name:           "Create Foreign Company",
		Description:    "Создать запись об иностранной компании-контрагента",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"pin", "fullName"},
		OptionalParams: []string{},
		Response:       "*ForeignCompany",
	},
	{
		Method:         "PUT",
		Path:           "/api/command/legal-person/foreign-company/edit/{id}",
		Name:           "Edit Foreign Company",
		Description:    "Отредактировать запись об иностранной компании",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{"pin", "fullName"},
		Response:       "*ForeignCompany",
	},
	{
		Method:         "DELETE",
		Path:           "/api/command/legal-person/foreign-company/delete/{id}",
		Name:           "Delete Foreign Company",
		Description:    "Удалить запись об иностранной компании",
		Category:       "Command",
		Authenticated:  true,
		RequiredParams: []string{"id"},
		OptionalParams: []string{},
		Response:       "map[string]interface{}",
	},

	// ========== DICTIONARY SERVICE ==========

	{
		Method:         "GET",
		Path:           "/api/dictionary-service/dictionary-types",
		Name:           "Get Dictionary Types",
		Description:    "Получить список типов справочников в системе",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"name", "shortName"},
		Response:       "[]*DictionaryType",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/dictionaries",
		Name:           "Get Dictionaries",
		Description:    "Получить элементы справочника по типу",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{"type"},
		OptionalParams: []string{"sort", "size", "page"},
		Response:       "[]*DictionaryItem",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/banks",
		Name:           "Get Banks",
		Description:    "Получить список всех банков в системе",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{},
		Response:       "[]*Bank",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/goods-classification",
		Name:           "Get Goods Classification",
		Description:    "Получить классификацию товаров по кодам ТНВЭД",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"size", "name", "sort"},
		Response:       "[]*GoodsClassification",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/countries",
		Name:           "Get Countries",
		Description:    "Получить список стран",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"size"},
		Response:       "[]*Country",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/document-statuses",
		Name:           "Get Document Statuses",
		Description:    "Получить возможные статусы документов в системе",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{},
		Response:       "[]*DocumentStatus",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/services-classification",
		Name:           "Get Services Classification",
		Description:    "Получить классификацию услуг по кодам ГКЭД",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"size"},
		Response:       "[]*ServicesClassification",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/taxes",
		Name:           "Get Taxes",
		Description:    "Получить налоговые ставки и виды налогов в системе",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"size"},
		Response:       "[]*TaxType",
	},
	{
		Method:         "GET",
		Path:           "/api/dictionary-service/legal-person-info",
		Name:           "Get Legal Person Info",
		Description:    "Получить информацию о юридических лицах из реестра",
		Category:       "Dictionary",
		Authenticated:  true,
		RequiredParams: []string{},
		OptionalParams: []string{"isResident", "fullName", "tin", "size"},
		Response:       "[]*Party",
	},

	// ========== AUTHENTICATION ==========

	{
		Method:         "POST",
		Path:           "/api/account/auth",
		Name:           "Authenticate",
		Description:    "Получить токен авторизации от сервиса Инфокома",
		Category:       "Auth",
		Authenticated:  false,
		RequiredParams: []string{"username", "password"},
		OptionalParams: []string{},
		Response:       "map[string]interface{}",
	},
}

// GetRoutesByCategory возвращает все роутеры определенной категории
func GetRoutesByCategory(category string) []Route {
	var result []Route
	for _, route := range APIRoutes {
		if route.Category == category {
			result = append(result, route)
		}
	}
	return result
}

// GetRoutesByMethod возвращает все роутеры определенного HTTP метода
func GetRoutesByMethod(method string) []Route {
	var result []Route
	for _, route := range APIRoutes {
		if route.Method == method {
			result = append(result, route)
		}
	}
	return result
}

// FindRoute ищет роутер по пути
func FindRoute(path string) *Route {
	for i, route := range APIRoutes {
		if route.Path == path {
			return &APIRoutes[i]
		}
	}
	return nil
}

// GetRouteCount возвращает общее количество роутеров
func GetRouteCount() int {
	return len(APIRoutes)
}

// GetRoutesByCategory возвращает количество роутеров по категориям
func GetRouteCountByCategory() map[string]int {
	counts := make(map[string]int)
	for _, route := range APIRoutes {
		counts[route.Category]++
	}
	return counts
}

// GetRouteCountByMethod возвращает количество роутеров по методам
func GetRouteCountByMethod() map[string]int {
	counts := make(map[string]int)
	for _, route := range APIRoutes {
		counts[route.Method]++
	}
	return counts
}
