<div align="center">

# BlackSearch V4

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://gin-gonic.com/)
[![OSINT](https://img.shields.io/badge/OSINT-Intelligence-red?style=for-the-badge)](https://github.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

**Твой Юнит в твоем самосуде!**

</div>

---

## 🔍 Возможности

| Модуль | Описание |
|--------|----------|
| 📱 **Анализ номера** | Многоисточниковая проверка телефона (API Ninjas, Numverify, Abstract API) |
| 🌐 **Анализ IP** | Геолокация, провайдер, ASN через 6 источников одновременно |
| 📧 **Анализ Email** | MX-записи, SMTP-верификация, проверка через API |
| 🌍 **Анализ домена** | WHOIS, DNS, геолокация IP |
| 🔎 **DNS-анализ** | A, AAAA, MX, NS, TXT, CNAME, SOA записи |
| 👤 **Поиск по нику** | Поиск аккаунтов на 28+ платформах одновременно |
| 📋 **Поиск по ФИО** | Поиск в базах ФНС через ofdata.ru |
| 🏢 **ИНН / ОГРН** | Поиск организаций и физлиц по реквизитам |
| 📂 **Локальная БД** | Поиск по CSV/TXT/JSON файлам с подсветкой |
| 🔵 **VK поиск** | Анализ профиля ВКонтакте |
| 🐙 **GitHub поиск** | Поиск пользователей, репозиториев, кода |
| 🔗 **Google Dork** | Генерация дорк-ссылок по категориям |
| 💳 **BIN-поиск** | Информация о банковской карте |
| 📍 **Гео-поиск** | Адрес ↔ Координаты через Nominatim |
| 🕸️ **Граф связей** | Веб-интерфейс на Gin + персистентное JSON хранилище |
| 📁 **Создание досье** | HTML-профиль с фото |
| 📊 **Создание отчёта** | Экспорт в HTML / TXT / Markdown / PDF-ready |

---

## 🏗️ Архитектура

```
BlackSearchV4/
├── main.go                    # Главный файл, меню
├── go.mod
├── go.sum
├── progress_bar/              # Утилиты вывода 
│   └── progress.go
├── html/                      # HTML-шаблоны (веб-интерфейс)
│   └── graph.html             # Граф связей (vis-network)
├── graph_data/                # Персистентное хранилище графов
│   └── graph.json             # (создаётся автоматически)
└── modules/
    ├── shared.go              # Общие утилиты, HTTP-клиент
    ├── phone.go               # Анализ номера телефона
    ├── ip.go                  # Анализ IP-адреса
    ├── email.go               # Анализ Email
    ├── domain.go              # Анализ домена
    ├── dns.go                 # DNS-анализ
    ├── nick.go                # Поиск по нику/username
    ├── fio.go                 # Поиск по ФИО/ИНН/ОГРН/компании/дорки/BIN
    ├── local_db.go            # Поиск по локальной базе
    ├── vk.go                  # VK поиск
    ├── github.go              # GitHub поиск (с красивым выводом репо)
    ├── geo.go                 # Гео-поиск
    ├── graph.go               # Граф связей (Gin сервер + JSON DB)
    ├── dossier.go             # Создание HTML-досье
    └── report.go              # Генерация отчётов
```

---

## 🔧 Зависимости

- [`charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss) — красивый TUI
- [`gin-gonic/gin`](https://github.com/gin-gonic/gin) — веб-фреймворк для граф-сервера
- [`likexian/whois`](https://github.com/likexian/whois) — WHOIS запросы
- [`likexian/whois-parser`](https://github.com/likexian/whois-parser) — парсинг WHOIS
- [`miekg/dns`](https://github.com/miekg/dns) — DNS-запросы
- [`vis-network`](https://visjs.github.io/vis-network/) — визуализация графа в браузере (CDN)

---

## 🚀 Инструкция по сборке

```bash
# 1. Клонируйте и войдите в папку
git clone https://github.com/yourname/BlackSearchV4.git
cd BlackSearchV4

# 2. Установите зависимости
go mod tidy

# 3. Запустите
go run .

# 4. Или соберите бинарник
go build -o blacksearch .
./blacksearch
```

> **Важно:** При запуске граф-сервера (пункт 17) программа запускает Gin HTTP-сервер на порту 8765.
> Откройте браузер: http://127.0.0.1:8765
> Данные графа сохраняются в `graph_data/graph.json` и переживают перезапуск.

---

## ⚙️ Переменные окружения

| Переменная | Описание |
|-----------|----------|
| `GITHUB_TOKEN` | GitHub Personal Access Token (60 → 5000 запросов/час) |

---

## 📁 Выходные файлы

Программа создаёт файлы в текущей директории:
- `nick_<ник>_<дата>.txt` — результаты поиска по нику
- `dorks_<запрос>.txt` — Google Dork ссылки
- `dossier_<фамилия>_<имя>_<id>.html` — HTML досье
- `report_<id>.html/txt/md` — отчёты
- `graph_data/graph.json` — данные графа связей


