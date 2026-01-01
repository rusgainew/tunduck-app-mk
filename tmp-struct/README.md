# Go Структуры для ESF API Document Service

## Описание

Сгенерированные Go структуры на основе анализа официального ESF API (Electronic Sales Invoices - Электронные счета-фактуры) Системы ТУНДУК Кыргызской Республики.

## Структура файлов

### Базовые справочные типы

#### `reference_item.go`

- **ReferenceItem** - Универсальный тип для всех справочных значений с кодом и названием
  - Используется для: PaymentType, Currency, Status, ReceiptType, DeliveryType, GoodsType, UnitClassification

#### `vat_tax_type.go`

- **VATTaxType** - Информация о налоге НДС (ставка и название)

#### `country.go`

- **Country** - Информация о стране (код и название)

#### `party.go`

- **Party** - Юридическое лицо или контрагент (ИНН, наименование, и т.д.)

### Основные сущности

#### `invoice.go`

- **Invoice** - Счет-фактура (Реализация, Приобретение, Корректировочная)
  - 30+ полей, включая идентификаторы, даты, финансовые данные, информацию о сторонах
  - Вложенные объекты: ReferenceItem, VATTaxType, Party, Country

#### `invoice_detail.go`

- **InvoiceDetail** - Строка детали счета (товары/услуги)
  - Содержит информацию о товаре: количество, цена, суммы НДС/НСП
  - Коды ТНВЭД, ГКЭД, справочные данные

#### `bank.go`

- **Bank** - Информация о банке (название, БИК, адрес)

#### `bank_account.go`

- **BankAccount** - Банковский счет юридического лица
  - Номер счета, валюта, информация о банке, статусы

#### `catalog.go`

- **Catalog** - Справочник товаров/услуг организации
  - Наименование, коды, классификация

#### `foreign_company.go`

- **ForeignCompany** - Иностранная компания/контрагент
  - Базовые реквизиты: ИНН, наименование

### Дополнительные типы

#### `page_info.go`

- **PageInfo** - Информация о пагинации
  - Номер страницы, размер, общее количество элементов, страниц

#### `api_response.go`

- **APIResponse** - Базовый обобщенный ответ API

  - RequestID, ResponseID, UserID, ClientID
  - Generic data поле для различных типов ответов

- **InvoiceListResponse** - Специализированный ответ со списком счетов
- **DetailListResponse** - Ответ со списком деталей
- **BankAccountListResponse** - Ответ со списком банковских счетов
- **CatalogListResponse** - Ответ со списком справочников

#### `requests.go`

Определены запросы для всех основных операций:

- **CreateInvoiceRequest** - Создание счета
- **UpdateInvoiceRequest** - Обновление счета
- **AcceptOrRejectRequest** - Принятие/отклонение счета
- **RevokeRequest** - Отзыв документа
- **SignRequest** - Подпись документа
- **CreateBankAccountRequest** - Создание счета
- **UpdateBankAccountRequest** - Обновление счета
- **CreateCatalogRequest** - Создание справочника
- **UpdateCatalogRequest** - Обновление справочника
- **CreateForeignCompanyRequest** - Создание иностранной компании
- **UpdateForeignCompanyRequest** - Обновление иностранной компании

#### `errors.go`

- **ErrorResponse** - Представление ошибки
- **ErrorCodes** - Константы кодов ошибок
- **StatusCodes** - Коды статусов документов
- **CorrectionReasons** - Причины корректировки

## Примечания по типам данных

### JSON теги

- Все структуры имеют JSON теги для сериализации/десериализации
- Используется camelCase в JSON (как в оригинальном API)

### Типы данных

- `string` - Текстовые данные, коды, ИНН, УУИДы
- `int64` - Целые числа (ID)
- `float64` - Decimal значения (суммы, курсы)
- `*string` - Опциональные строковые поля (даты, комментарии)
- `bool` - Логические флаги
- `*ReferenceItem`, `*VATTaxType`, `*Party`, `*Country` - Вложенные объекты со справочными данными
- `[]*InvoiceDetail` - Массивы (детали счета)

### Дата-время

- Все даты в API представлены строками в формате ISO 8601 (YYYY-MM-DD)
- В структурах используется `*string` для дат (опциональное поле)

## Использование в Document-Service

Эти структуры предназначены для использования в:

1. **Domain Layer** - Как основные доменные модели
2. **Application Layer** - Для обработки бизнес-логики
3. **Infrastructure Layer** - Для преобразования в/из внешних форматов
4. **Interfaces Layer** - Для HTTP запросов/ответов и gRPC сообщений

## Пример использования

```go
import "document-service/internal/models"

// Создание счета
invoice := &models.Invoice{
    DocumentUUID:   "uuid-12345",
    InvoiceNumber:  "INV-001",
    TotalAmount:    1000.00,
    IsResident:     true,
    Currency: &models.ReferenceItem{
        Code: "KGS",
        Name: "Сом",
    },
    LegalPerson: &models.Party{
        PIN:      "12345678901",
        FullName: "ООО Компания",
    },
}

// Получение списка счетов
response := &models.InvoiceListResponse{
    RequestID:     ptr("req-123"),
    ResponseID:    ptr("resp-456"),
    TotalElements: 100,
    TotalPage:     10,
    Invoices:      []*models.Invoice{invoice},
}
```

## Расширение структур

При необходимости добавить новые поля:

1. Добавить поле в соответствующую структуру
2. Добавить JSON тег с camelCase именем
3. Обновить комментарий документации
4. Если новое поле требует нового типа - создать отдельный файл

## Соответствие с ESF API

Все структуры скопированы напрямую из официальной документации ESF API:

- Названия полей соответствуют API (с преобразованием в Go convention)
- Типы данных соответствуют спецификации
- Обязательные/опциональные поля отражены в структурах
