# 🎯 START HERE

**Дата:** 1 января 2026  
**Статус:** ✅ Консолидация документации завершена

---

## 📌 ГЛАВНЫЙ ДОКУМЕНТ

# 👉 Прочитайте: [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md)

Это **единственный документ для работы** с проектом.

**Содержит:**

- ⚡ Quickstart (10 минут)
- 🏗️ Архитектура DDD (полная)
- 📋 Все доменные области (Auth, Company, Document)
- 🔌 Service Communication (gRPC + RabbitMQ)
- 📅 План разработки (6 фаз)
- 🛠️ Инструменты и технологии
- 🧪 Testing strategy
- 🚀 Getting Started для разработчиков
- ❓ FAQ & Troubleshooting

**Время:** 30-45 минут на чтение

---

## 📚 Дополнительные документы (если нужны детали)

| Файл                          | Для чего                                            |
| ----------------------------- | --------------------------------------------------- |
| **ANALYSIS_AND_PLAN.md**      | 👈 **НОВЫЙ** План разработки всех 6 фаз (13 недель) |
| **DEVELOPMENT_PLAN.md**       | 👈 **НОВЫЙ** Детальный план с кодом и тестами       |
| `PHASE1_COMPLETION_REPORT.md` | Что было сделано в Phase 1                          |
| `PHASE2_AUTH_SERVICE_PLAN.md` | Детальный план разработки Auth-Service              |
| `PHASE2_COMPLETION_REPORT.md` | Статус Phase 2                                      |
| `TEAM_BRIEFING.md`            | Информация для команды                              |
| `api/proto/README.md`         | Работа с proto файлами                              |

---

## 🚀 Быстрый старт (3 команды)

```bash
# 1. Клонировать и перейти
git clone https://github.com/rusgainew/tunduck-app-mk.git
cd tunduck-app-mk

# 2. Запустить инфраструктуру
docker-compose -f docker-compose.microservices.yml up -d

# 3. Компилировать proto
cd api/proto && make proto
```

Готово! 🎉

---

## 📖 Для каких ролей

**👨‍💼 Менеджер/Архитектор:**  
→ [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md) (раздел "План разработки")

**👨‍💻 Разработчик Go:**  
→ [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md) (раздел "Getting Started")

**🏗️ Разработчик Auth-Service:**  
→ [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md) + [PHASE2_AUTH_SERVICE_PLAN.md](PHASE2_AUTH_SERVICE_PLAN.md)

**🔧 DevOps инженер:**  
→ [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md) (раздел "Deployment Architecture")

**📡 Фронтенд разработчик:**  
→ [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md) (раздел "Service Communication")

---

## ❓ Вопросы?

Смотрите раздел **FAQ & Troubleshooting** в [PROJECT_MASTER_GUIDE.md](PROJECT_MASTER_GUIDE.md)
