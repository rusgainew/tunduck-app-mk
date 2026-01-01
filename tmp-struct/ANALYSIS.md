# Анализ ESF API документации

## Краткое описание сервиса

**ESF API (Электронные счета-фактуры)** - REST API для взаимодействия с информационной системой ЭСФ через СМЭВ "ТУНДУК".
Используется для управления счетами-фактурами: получение, создание, редактирование, удаление, подписание, корректировка.

## Основные Сущности (Entities)

### 1. **Invoice (Счет-фактура)**

- Три типа счетов: Realization (Реализация), Income (Приобретение), Correction (Корректировочная)
- Общие поля для всех типов

**Основные атрибуты:**

- `documentUuid` (UUID) - Уникальный идентификатор документа
- `invoiceNumber` (String) - Номер документа
- `invoiceDate` (Date) - Дата подписания
- `createdDate` (Date) - Дата создания
- `deliveryDate` (Date) - Дата поставки
- `totalAmount` (Decimal) - Общая сумма
- `isResident` (Boolean) - Статус резидентства
- `note` (String) - Комментарии
- `ownedCrmReceiptCode` (String) - Код учетной системы
- `number` (String) - Номер
- `correctedReceiptUuid` (String) - UUID корректируемого документа

**Финансовые поля:**

- `isPriceWithoutTaxes` (Boolean) - Цена без налогов
- `openingBalances` (Decimal) - Сальдо на начало
- `assessedContributionsAmount` (Decimal) - Начислено
- `paidAmount` (Decimal) - Оплачено
- `penaltiesAmount` (Decimal) - Штраф
- `finesAmount` (Decimal) - Пеня
- `closingBalances` (Decimal) - Сальдо на конец
- `amountToBePaid` (Decimal) - К оплате
- `personalAccountNumber` (String) - Лицевой счет

**Реквизиты счетов:**

- `legalPersonBankAccount` (String) - Счет юр. лица
- `contractorBankAccount` (String) - Счет контрагента

**Для корректировочных счетов:**

- `correctedReceiptCode` (String) - Номер корректируемого счета
- `correctionReasonCode` (String) - Код причины корректировки
- `correctionReasonName` (String) - Наименование причины
- `correctedReceiptCreationDate` (Date) - Дата создания корректируемого счета

**Для международных операций:**

- `exchangeRate` (Decimal) - Курс обмена валюты к сому
- `foreignName` (String) - Наименование иностранца
- `sellerBranchPin` (String) - ИНН поставщика

**Прочие поля:**

- `isIndustry` (Boolean) - Отраслевые
- `deliveryContractNumber` (String) - Номер договора поставки
- `deliveryContractDate` (Date) - Дата договора поставки

### 2. **ReferenceItem (Справочный элемент)**

Используется для различных справочников (статусы, типы, валюты, платежи и т.д.)

**Атрибуты:**

- `code` (String) - Код элемента
- `name` (String) - Наименование элемента

**Типы ReferenceItem:**

- PaymentType - Тип оплаты
- Currency - Валюта
- Status - Статус документа
- ReceiptType - Тип операции (счет-фактура)
- DeliveryType - Вид поставки
- VATDelivery - Вид поставки для НДС
- GoodsType - Тип товаров
- UnitClassification - Тип единиц измерения

### 3. **VATTaxType (Вид налога НДС)**

**Атрибуты:**

- `rate` (String) - Ставка (например, "12%", "0%")
- `name` (String) - Наименование ставки

### 4. **Party (Юридическое лицо / Контрагент)**

**Атрибуты:**

- `pin` (String) - ИНН
- `fullName` (String) - Полное наименование
- `mainFullName` (String) - Наименование головной организации
- `mainPin` (String) - ИНН головной организации

### 5. **Country (Страна)**

**Атрибуты:**

- `code` (String) - Код страны
- `name` (String) - Наименование страны

### 6. **InvoiceDetail (Строка детали счета)**

Описывает товары/услуги в счете-фактуре

**Атрибуты:**

- `id` (Long) - Идентификатор
- `invoiceUuid` (UUID) - UUID счета-фактуры
- `baseCount` (Decimal) - Количество
- `price` (Decimal) - Цена за единицу
- `amount` (Decimal) - Сумма (цена × количество)
- `amountWithoutVAT` (Decimal) - Сумма без НДС
- `amountVAT` (Decimal) - Сумма НДС
- `amountST` (Decimal) - Сумма НСП
- `fcdNumber` (String) - Номер ГТД (свидетельство о происхождении)
- `goodsName` (String) - Наименование товара/услуги
- `tnvedCode` (String) - Код ТНВЭД
- `gked` (String) - Код ГКЭД (государственный классификатор видов экономической деятельности)

**Справочные поля:**

- `unitClassification` (ReferenceItem) - Тип единиц измерения
- `goodsType` (ReferenceItem) - Тип товаров (товары, услуги)
- `stTaxType` (VATTaxType) - Вид налога НСП

### 7. **BankAccount (Банковский счет)**

**Атрибуты:**

- `id` (Long) - Идентификатор
- `accountName` (String) - Название счета
- `bankAccount` (String) - Номер банковского счета
- `currency` (ReferenceItem) - Валюта
- `bank` (Bank) - Информация о банке
- `isUsed` (Boolean) - Используется ли
- `contractorTin` (String) - ИНН контрагента
- `isResident` (Boolean) - Является ли резидентом

### 8. **Bank (Банк)**

**Атрибуты:**

- `id` (Long) - Идентификатор банка
- `name` (String) - Название банка
- `bik` (String) - БИК банка
- `addressText` (String) - Адрес банка

### 9. **Catalog (Справочник юр. лица)**

Справочник товаров/услуг организации

**Атрибуты:**

- `id` (Long) - Идентификатор
- `name` (String) - Наименование товара/услуги
- `number` (String) - Номер товара
- `tnvedCode` (String) - ТН ВЭД код
- `gkedCode` (String) - ГКЭД код
- `unitClassification` (ReferenceItem) - Тип единиц измерения
- `catalogType` (ReferenceItem) - Тип каталога (товары/услуги)

### 10. **ForeignCompany (Иностранная компания)**

**Атрибуты:**

- `id` (Long)
- `pin` (String)
- `fullName` (String)
- Другие поля для контакта и данных компании

### 11. **ApiResponse (Обобщенный ответ API)**

Используется для всех responses

**Атрибуты:**

- `requestId` (String) - ID запроса
- `responseId` (UUID) - ID ответа
- `userId` (String) - ID пользователя
- `clientId` (String) - ID клиента
- Специфичные для endpoint данные (invoices, details, banks, catalogs и т.д.)
- `totalElements` (Integer) - Общее количество элементов
- `totalPage` (Integer) - Общее количество страниц

## Основные Операции

### Документы (Documents)

1. **GET** `/api/query/invoice/{type}/document-by-exchange-code` - Получить счет по UUID
2. **GET** `/api/query/invoice/{type}/all-documents` - Получить список счетов (с фильтрами и пагинацией)
3. **POST** `/api/command/invoice/{type}/create` - Создать счет
4. **PUT** `/api/command/invoice/{type}/edit/{id}` - Редактировать счет
5. **DELETE** `/api/command/invoice/{type}/delete/{id}` - Удалить счет
6. **POST** `/api/command/invoice/{type}/accept-or-reject` - Принять/отклонить счет
7. **POST** `/api/command/invoice/{type}/revoke` - Отозвать счет
8. **POST** `/api/command/invoice/{type}/sign` - Подписать счет

### Детали (Details)

1. **GET** `/api/query/detail/details-by-exchange-code` - Получить детали счета (с пагинацией)
2. **GET** `/api/query/detail/v2/details-by-exchange-code` - Получить детали счета (без пагинации)

### Банковские счета (Bank Accounts)

1. **GET** `/api/query/legal-person/bank-accounts` - Получить список счетов
2. **GET** `/api/query/legal-person/bank-account/{id}` - Получить счет по ID
3. **POST** `/api/command/legal-person/bank-account/create` - Создать счет
4. **PUT** `/api/command/legal-person/bank-account/edit/{id}` - Редактировать счет
5. **DELETE** `/api/command/legal-person/bank-account/delete/{id}` - Удалить счет

### Справочники (Catalogs)

1. **GET** `/api/query/legal-person/catalogs` - Получить справочники
2. **GET** `/api/query/legal-person/catalog/{id}` - Получить справочник по ID
3. **POST** `/api/command/legal-person/catalog/create` - Создать справочник
4. **PUT** `/api/command/legal-person/catalog/edit/{id}` - Редактировать справочник
5. **DELETE** `/api/command/legal-person/catalog/delete/{id}` - Удалить справочник

### Иностранные компании (Foreign Companies)

1. **GET** `/api/query/foreign-company/...` - Получить иностранные компании
2. **POST** `/api/command/foreign-company/create` - Создать иностранную компанию
3. **PUT** `/api/command/foreign-company/edit/{id}` - Редактировать компанию
4. **DELETE** `/api/command/foreign-company/delete/{id}` - Удалить компанию

### Справочные данные (Dictionary Data)

1. **GET** `/api/query/dictionary-types` - Типы справочников
2. **GET** `/api/query/dictionaries` - Справочники
3. **GET** `/api/query/banks` - Список банков
4. **GET** `/api/query/goods-classification` - Классификация товаров
5. **GET** `/api/query/countries` - Страны
6. **GET** `/api/query/document-statuses` - Статусы документов
7. **GET** `/api/query/services-classification` - Классификация услуг
8. **GET** `/api/query/taxes` - Налоги
9. **GET** `/api/query/legal-person-info` - Информация о юр. лице

## Примечания

- Все дата-время в формате ISO 8601 (YYYY-MM-DD)
- Decimal поля используют точное представление для финансовых расчетов
- Все UUID формата string
- Пагинация с `page` (начиная с 0) и `size` параметрами
- Авторизация через Bearer Token
- Требуемые headers: X-Road-Client, ClientUUID, Authorization, USER-TIN
