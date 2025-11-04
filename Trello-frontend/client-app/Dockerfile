# Stage 1: Build the Angular application
FROM node:20 AS builder

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json for dependency installation
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . ./

# Build the Angular application
RUN npm run build --prod

# Stage 2: Serve the application using a lightweight server
FROM nginx:alpine

# Copy the built application from the builder stage
COPY --from=builder /app/dist/client-app/browser /usr/share/nginx/html

# Expose the port that NGINX is listening on
EXPOSE 80

# Start NGINX server
CMD ["nginx", "-g", "daemon off;"]
