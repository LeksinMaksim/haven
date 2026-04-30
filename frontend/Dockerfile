# Этап 1: Сборка
FROM node:22-alpine AS build

WORKDIR /app

# Копируем файлы зависимостей
COPY package*.json ./
RUN npm install

# Копируем остальной код и собираем проект
COPY . .
RUN npm run build

# Этап 2: Раздача статики через Nginx
FROM nginx:stable-alpine

# Копируем билд из первого этапа в папку Nginx
COPY --from=build /app/dist /usr/share/nginx/html

# Добавляем конфиг для Nginx (чтобы работал роутинг, если появится)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
