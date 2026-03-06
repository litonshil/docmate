# Development Dockerfile for Next.js
FROM node:20-alpine

WORKDIR /app

# Install pnpm and configure npm for better reliability
RUN npm config set fetch-retries 5 && \
    npm config set fetch-retry-mintimeout 20000 && \
    npm config set fetch-retry-maxtimeout 120000 && \
    npm install -g pnpm

# Copy package files first for better caching
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy the rest of the application
COPY . .

CMD ["pnpm", "run", "dev"]
