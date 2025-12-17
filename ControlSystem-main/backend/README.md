## Авторизация

### **1. Login**

**Что делает:** Авторизация обычного пользователя.
**Адрес:** `POST {{baseURL}}/auth/login`
**Метод:** `POST`

**Параметры (Body JSON):**

```json
{
  "email": "test.menedzherov@controlsystem.ru",
  "password": "pXzvSX45wsQ8"
}
```

**Возвращает:**
`token` — JWT-токен пользователя.

---

### **2. Register**

**Что делает:** Регистрация нового пользователя (доступно только администратору).
**Адрес:** `POST {{baseURL}}/admin/register`
**Метод:** `POST`
**Авторизация:** Bearer Token (админ)

**Параметры (Body JSON):**

```json
{
  "first_name": "Тест",
  "middle_name": "Инженерович",
  "last_name": "Менеджеров",
  "email": "gleblobachev@gmail.com",
  "role": 1
}
```

---

##  Проекты

---

### **3. Create Project**

**Что делает:** Создает новый проект.
**Адрес:** `POST {{baseURL}}/projects`
**Метод:** `POST`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "name": "Cool project name",
  "description": "description"
}
```

---

### **4. Edit Project**

**Что делает:** Редактирует данные существующего проекта.
**Адрес:** `PATCH {{baseURL}}/projects/{id}`
**Метод:** `PATCH`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "name": "",
  "description": "fggh",
  "status": 2
}
```

---

### **5. Get Projects**

**Что делает:** Получает список всех проектов.
**Адрес:** `GET {{baseURL}}/projects/`
**Метод:** `GET`
**Авторизация:** Bearer Token

**Параметры:** отсутствуют.

---

### **6. Get Project**

**Что делает:** Получает данные конкретного проекта.
**Адрес:** `GET {{baseURL}}/projects/{id}`
**Метод:** `GET`
**Авторизация:** Bearer Token

---

### **7. Assign Engineer**

**Что делает:** Назначает инженера на проект.
**Адрес:** `POST {{baseURL}}/projects/{id}/assign`
**Метод:** `POST`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "engineer_id": 2
}
```

---

## Дефекты

---

### **8. Create Defect**

**Что делает:** Создает новый дефект в рамках проекта.
**Адрес:** `POST {{baseURL}}/defects/`
**Метод:** `POST`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "title": "dfvdfvb",
  "description": "fgvdfgfddf",
  "project_id": 1
}
```

---

### **9. Get Defects**

**Что делает:** Получает список дефектов по проекту.
**Адрес:** `GET {{baseURL}}/defects/`
**Метод:** `GET`
**Авторизация:** Bearer Token

**Query-параметры:**

| Параметр  | Тип    | Описание         |
| --------- | ------ | ---------------- |
| projectId | int    | ID проекта       |
| page      | int    | Номер страницы   |
| search    | string | Поисковая строка |

**Пример:**

```
{{baseURL}}/defects/?projectId=1&page=1&search=
```

---

### **10. Get Defect**

**Что делает:** Получает данные конкретного дефекта.
**Адрес:** `GET {{baseURL}}/defects/{id}`
**Метод:** `GET`
**Авторизация:** Bearer Token

---

### **11. Edit Defect**

**Что делает:** Редактирует существующий дефект.
**Адрес:** `PATCH {{baseURL}}/defects/{id}`
**Метод:** `PATCH`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "title": "",
  "description": "",
  "status": 3,
  "priority": 1
}
```

---

## Комментарии

---

### **12. Leave Comment**

**Что делает:** Добавляет комментарий к дефекту.
**Адрес:** `POST {{baseURL}}/defects/{id}/comments`
**Метод:** `POST`
**Авторизация:** Bearer Token

**Параметры (Body JSON):**

```json
{
  "content": "Tetsxv1654fcf"
}
```

---

### **13. Get Comments**

**Что делает:** Получает список комментариев к дефекту.
**Адрес:** `GET {{baseURL}}/defects/{id}/comments`
**Метод:** `GET`
**Авторизация:** Bearer Token

---

## Вложения

---

### **14. Add Attachment**

**Что делает:** Добавляет вложение (файл) к дефекту или проекту.
**Адрес:** `POST {{baseURL}}/attachments/`
**Метод:** `POST`
**Авторизация:** Bearer Token

**Параметры (FormData):**

| Поле      | Тип  | Пример | Описание                   |
| --------- | ---- | ------ | -------------------------- |
| projectId | text | 1      | ID проекта *(опционально)* |
| defectId  | text | 2      | ID дефекта                 |
| file      | file | —      | Файл для загрузки          |

---

## Пользователи (Админ)

---

### **15. Get Users**

**Что делает:** Получает список пользователей с фильтрацией.
**Адрес:** `GET {{baseURL}}/admin/get-users`
**Метод:** `GET`
**Авторизация:** Bearer Token

**Query-параметры:**

| Параметр   | Тип    | Описание             |
| ---------- | ------ | -------------------- |
| page       | int    | Номер страницы       |
| email      | string | Фильтр по email      |
| role       | int    | Фильтр по роли       |
| is_enabled | bool   | Фильтр по активности |

**Пример:**

```
{{baseURL}}/admin/get-users?page=1&email=&role=&is_enabled=
```

---

### **16. Edit User**

**Что делает:** Редактирует данные пользователя.
**Адрес:** `PATCH {{baseURL}}/admin/edit-user/{id}`
**Метод:** `PATCH`
**Авторизация:** Bearer Token (админ)

**Параметры (Body JSON):**

```json
{
  "first_name": "",
  "middle_name": "tesdfg",
  "last_name": "",
  "role": 2,
  "is_enabled": false
}
```

---

##  Переменные окружения

| Переменная    | Описание                               |
| ------------- | -------------------------------------- |
| `{{baseURL}}` | Базовый URL API сервера                |
| `{{token}}`   | JWT токен авторизованного пользователя |

