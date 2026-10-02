FROM node:24.19.0-alpine3.24 AS build

WORKDIR /app

COPY package.json package-lock.json .npmrc ./
RUN npm ci

COPY index.html tsconfig.json vite.config.ts ./
COPY src/ ./src/

RUN npm run test -- --run
RUN npm run build

FROM nginx:1.30.5-alpine3.24

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist/ /usr/share/nginx/html/

EXPOSE 80
