#build stage
FROM node:12.2.0-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# production stage
FROM nginx:1.13.12-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
# Copy the respective nginx configuration files
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]