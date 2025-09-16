# Stage 1: Build the React application
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY ide/frontend/package.json ide/frontend/package-lock.json ./
RUN npm install

# Copy the rest of the frontend source code
COPY ide/frontend/ .

# Build the application for production
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the built assets from the builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy the Nginx configuration file
COPY ide/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
