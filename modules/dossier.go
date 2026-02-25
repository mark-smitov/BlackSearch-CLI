package modules

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type DossierData struct {
	Surname          string
	Name             string
	Patronymic       string
	Gender           string
	BirthDate        string
	BirthPlace       string
	Citizenship      string
	Nationality      string
	MaritalStatus    string
	Children         string
	Height           string
	Weight           string
	EyeColor         string
	HairColor        string
	Marks            string
	Passport         string
	PassportIssue    string
	PassportDate     string
	INN              string
	SNILS            string
	Phone            string
	Email            string
	Address          string
	ActualAddress    string
	Education        string
	Occupation       string
	Employer         string
	Income           string
	Languages        string
	SocialMedia      string
	EmergencyContact string
	Medical          string
	Criminal         string
	Military         string
	Photo            string
	Notes            string
	Date             string
	ID               string
}

func RunDossier() {
	d := &DossierData{}

	PrintSection("Личные данные")
	d.Surname = AskInput("Фамилия")
	d.Name = AskInput("Имя")
	d.Patronymic = AskInput("Отчество")
	d.Gender = AskInput("Пол")
	d.BirthDate = AskInput("Дата рождения")
	d.BirthPlace = AskInput("Место рождения")
	d.Citizenship = AskInput("Гражданство")
	d.Nationality = AskInput("Национальность")
	d.MaritalStatus = AskInput("Семейное положение")
	d.Children = AskInput("Дети")

	PrintSection("Внешность")
	d.Height = AskInput("Рост (см)")
	d.Weight = AskInput("Вес (кг)")
	d.EyeColor = AskInput("Цвет глаз")
	d.HairColor = AskInput("Цвет волос")
	d.Marks = AskInput("Особые приметы")

	PrintSection("Документы")
	d.Passport = AskInput("Паспорт (серия и номер)")
	d.PassportIssue = AskInput("Кем выдан")
	d.PassportDate = AskInput("Дата выдачи")
	d.INN = AskInput("ИНН")
	d.SNILS = AskInput("СНИЛС")

	PrintSection("Контакты")
	d.Phone = AskInput("Телефон")
	d.Email = AskInput("Email")
	d.Address = AskInput("Адрес регистрации")
	d.ActualAddress = AskInput("Адрес проживания")

	PrintSection("Социальное")
	d.Education = AskInput("Образование")
	d.Occupation = AskInput("Род занятий")
	d.Employer = AskInput("Место работы")
	d.Income = AskInput("Уровень дохода")
	d.Languages = AskInput("Иностранные языки")

	PrintSection("Сети и связи")
	d.SocialMedia = AskInput("Соцсети")
	d.EmergencyContact = AskInput("Контактное лицо")

	PrintSection("Дополнительно")
	d.Medical = AskInput("Медицинская информация")
	d.Criminal = AskInput("Судимости")
	d.Military = AskInput("Воинская обязанность")
	d.Photo = AskInput("Путь к фото или ссылка (Enter=пропустить)")
	d.Notes = AskInput("Примечания")

	now := time.Now()
	d.Date = now.Format("02.01.2006")
	d.ID = now.Format("20060102150405")

	sp := NewSpinner("досье")
	sp.Start()
	html := generateDossierHTML(d)
	sp.Stop()

	filename := fmt.Sprintf("dossier_%s_%s_%s.html", d.Surname, d.Name, d.ID)
	if err := os.WriteFile(filename, []byte(html), 0644); err != nil {
		PrintError("Ошибка записи: " + err.Error())
		return
	}

	PrintSuccess("Досье: " + filename)
}

func photoTag(photo string) string {
	if photo == "" {
		return `<div class="no-photo">НЕТ ФОТО</div>`
	}
	if strings.HasPrefix(photo, "http") {
		return fmt.Sprintf(`<img src="%s" class="photo-img" alt="Фото">`, photo)
	}
	absPath, _ := os.Getwd()
	return fmt.Sprintf(`<img src="file:///%s/%s" class="photo-img" alt="Фото">`, absPath, photo)
}

func dossierRow(label, value string) string {
	if value == "" {
		value = "—"
	}
	return fmt.Sprintf(`<tr><td class="td-label">%s</td><td class="td-value">%s</td></tr>`, label, value)
}

func generateDossierHTML(d *DossierData) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="ru">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Досье: %s %s</title>
<style>
  @import url('https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap');
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body { background: #080808; color: #c8c8c8; font-family: 'Inter', -apple-system, sans-serif; padding: 40px 20px; }
  .container { max-width: 960px; margin: 0 auto; }
  .header {
    display: flex; gap: 32px; align-items: flex-start;
    padding: 32px; background: #111; border: 1px solid #1e1e1e;
    border-radius: 8px; margin-bottom: 24px;
  }
  .photo-img { width: 160px; height: 200px; object-fit: cover; border-radius: 4px; border: 1px solid #2a2a2a; }
  .no-photo {
    width: 160px; height: 200px; border: 1px solid #222;
    background: #0d0d0d; display: flex; align-items: center;
    justify-content: center; color: #333; font-size: 11px;
    border-radius: 4px; letter-spacing: 2px; flex-shrink: 0;
  }
  .header-info { flex: 1; }
  .header-info h1 { font-size: 30px; font-weight: 700; letter-spacing: -1px; color: #fff; margin-bottom: 8px; }
  .badges { display: flex; flex-wrap: wrap; gap: 6px; margin-bottom: 16px; }
  .badge {
    display: inline-block; background: #161616; border: 1px solid #2a2a2a;
    color: #555; font-size: 10px; padding: 3px 10px; border-radius: 3px; letter-spacing: 1px;
  }
  .badge.accent { border-color: #6272a4; color: #6272a4; }
  .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
  .card { background: #111; border: 1px solid #1a1a1a; border-radius: 6px; overflow: hidden; }
  .card-title {
    background: #141414; border-bottom: 1px solid #1a1a1a;
    padding: 10px 16px; font-size: 9px;
    text-transform: uppercase; letter-spacing: 2px; color: #444;
  }
  table { width: 100%%; border-collapse: collapse; }
  .td-label { color: #555; width: 45%%; padding: 6px 14px; font-size: 12px; border-bottom: 1px solid #161616; }
  .td-value { color: #aaa; padding: 6px 14px; font-size: 12px; border-bottom: 1px solid #161616; word-break: break-word; }
  tr:last-child .td-label, tr:last-child .td-value { border-bottom: none; }
  .notes-card {
    margin-top: 16px; background: #111; border: 1px solid #1a1a1a;
    border-radius: 6px; padding: 20px;
  }
  .notes-content { color: #777; font-size: 13px; line-height: 1.8; white-space: pre-wrap; }
  .footer { text-align: center; margin-top: 24px; color: #2a2a2a; font-size: 10px; letter-spacing: 2px; text-transform: uppercase; }
  @media print {
    body { background: white; color: #222; }
    .card, .header, .notes-card { border-color: #ddd; background: white; }
    .card-title { background: #f5f5f5; color: #888; }
    .td-label { color: #888; }
    .td-value { color: #333; }
  }
</style>
</head>
<body>
<div class="container">
  <div class="header">
    %s
    <div class="header-info">
      <h1>%s %s %s</h1>
      <div class="badges">
        <span class="badge accent">ID: %s</span>
        <span class="badge">%s</span>
        <span class="badge">BLACKSEARCH</span>
        <span class="badge">КОНФИДЕНЦИАЛЬНО</span>
      </div>
    </div>
  </div>
  <div class="grid">
    <div class="card">
      <div class="card-title">Личные данные</div>
      <table>
        %s%s%s%s%s%s
      </table>
    </div>
    <div class="card">
      <div class="card-title">Документы</div>
      <table>
        %s%s%s%s%s
      </table>
    </div>
    <div class="card">
      <div class="card-title">Внешность</div>
      <table>
        %s%s%s%s%s
      </table>
    </div>
    <div class="card">
      <div class="card-title">Контакты</div>
      <table>
        %s%s%s%s
      </table>
    </div>
    <div class="card">
      <div class="card-title">Социальное</div>
      <table>
        %s%s%s%s%s
      </table>
    </div>
    <div class="card">
      <div class="card-title">Прочее</div>
      <table>
        %s%s%s%s
      </table>
    </div>
  </div>
  <div class="notes-card">
    <div class="card-title" style="margin: -20px -20px 16px; padding: 10px 20px;">Примечания</div>
    <div class="notes-content">%s</div>
  </div>
  <div class="footer">BlackSearch Intelligence · Создано: %s</div>
</div>
</body>
</html>`,
		d.Surname, d.Name,
		photoTag(d.Photo),
		d.Surname, d.Name, d.Patronymic,
		d.ID, d.Date,
		dossierRow("Пол", d.Gender),
		dossierRow("Дата рождения", d.BirthDate),
		dossierRow("Место рождения", d.BirthPlace),
		dossierRow("Гражданство", d.Citizenship),
		dossierRow("Национальность", d.Nationality),
		dossierRow("Семейное положение", d.MaritalStatus),
		dossierRow("Паспорт", d.Passport),
		dossierRow("Кем выдан", d.PassportIssue),
		dossierRow("Дата выдачи", d.PassportDate),
		dossierRow("ИНН", d.INN),
		dossierRow("СНИЛС", d.SNILS),
		dossierRow("Рост", orEmpty(d.Height, " см")),
		dossierRow("Вес", orEmpty(d.Weight, " кг")),
		dossierRow("Цвет глаз", d.EyeColor),
		dossierRow("Цвет волос", d.HairColor),
		dossierRow("Особые приметы", d.Marks),
		dossierRow("Телефон", d.Phone),
		dossierRow("Email", d.Email),
		dossierRow("Адрес регистрации", d.Address),
		dossierRow("Адрес проживания", d.ActualAddress),
		dossierRow("Образование", d.Education),
		dossierRow("Род занятий", d.Occupation),
		dossierRow("Место работы", d.Employer),
		dossierRow("Доход", d.Income),
		dossierRow("Языки", d.Languages),
		dossierRow("Соцсети", d.SocialMedia),
		dossierRow("Медицинская", d.Medical),
		dossierRow("Судимости", d.Criminal),
		dossierRow("Воинская обязанность", d.Military),
		func() string {
			if d.Notes == "" {
				return "—"
			}
			return d.Notes
		}(),
		d.Date,
	)
}

func orEmpty(val, suffix string) string {
	if val == "" {
		return ""
	}
	return val + suffix
}
