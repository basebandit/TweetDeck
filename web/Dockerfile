#build stage
FROM node:12.2.0-alpine as build-stage
WORKDIR /app
COPY package*.json /app/
RUN npm install yarn -g
COPY . .
RUN yarn install && yarn run build

# production stage
FROM nginx:1.13.12-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
# Copy the respective nginx configuration files
COPY ./nginx/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
