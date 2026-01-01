# Структура tmp-struct

Код в директории `tmp-struct` был организован в логические блоки для улучшения читаемости и поддержки.

## Структура пакетов

```
tmp-struct/
├── entities/           # Основные бизнес-сущности
│   ├── invoice.go      # Счет-фактура
│   ├── invoice_detail.go  # Детали/строки счета
│   └── party.go        # Юридическое лицо/контрагент
│
├── dictionaries/       # Справочники и классификаторы
│   ├── bank.go         # Банк
│   ├── bank_account.go # Банковский счет
│   ├── catalog.go      # Каталог товаров/услуг
│   ├── country.go      # Страна
│   ├── foreign_company.go  # Иностранная компания
│   ├── reference_item.go   # Справочный элемент
│   └── vat_tax_type.go     # Тип НДС
│
├── api/                # API-специфичные типы
│   ├── api_response.go # Базовый ответ API
│   ├── page_info.go    # Пагинация
│   ├── requests.go     # Запросы (Create/Update/etc.)
│   └── routes.go       # Маршруты API
│
└── errors/             # Типы и коды ошибок
    └── errors.go       # Ошибки API
```

## Описание пакетов

### `entities`

Содержит основные бизнес-сущности системы электронных счетов-фактур:

- **Invoice** - главный документ (счет-фактура)
- **InvoiceDetail** - строки/детали счета с товарами/услугами
- **Party** - юридическое лицо (поставщик/покупатель)

### `dictionaries`

Справочные данные и классификаторы:

- **ReferenceItem** - универсальный справочный элемент
- **VATTaxType** - типы и ставки НДС
- **Bank**, **BankAccount** - банковская информация
- **Catalog** - каталог товаров/услуг организации
- **Country** - страны
- **ForeignCompany** - иностранные контрагенты

### `api`

API-специфичные структуры:

- **APIResponse** - стандартный формат ответа
- **PageInfo** - информация о пагинации
- **CreateInvoiceRequest**, **UpdateInvoiceRequest**, и др. - запросы к API
- **Route** - описание маршрутов API

### `errors`

Обработка ошибок:

- **ErrorResponse** - формат ошибки
- **ErrorCodes** - константы кодов ошибок

## Зависимости между пакетами

```
entities ───► dictionaries
    │
    ▼
   api  ────────► dictionaries
                       │
                       ▼
                    errors
```

- `entities` использует типы из `dictionaries`
- `api` использует типы из `entities` и `dictionaries`
- `errors` независим от других пакетов

## Импорты

При использовании типов из других пакетов необходимо импортировать соответствующий пакет:

```go
import (
    "github.com/rusgainew/tunduck-app-mk/tmp-struct/entities"
    "github.com/rusgainew/tunduck-app-mk/tmp-struct/dictionaries"
)
```
