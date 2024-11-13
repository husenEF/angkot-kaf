# Build stage
FROM oven/bun:1 AS builder

WORKDIR /app

# Copy package files
COPY package.json bun.lockb ./

# Install dependencies
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build the application
RUN bun build ./src/index.ts --outdir ./dist

# Final stage
FROM oven/bun:1-slim

WORKDIR /app

# Copy built files from builder
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./

# Create database directory
RUN mkdir -p database

# Install production dependencies only
RUN bun install --production --frozen-lockfile

# Run the application
CMD ["bun", "run", "dist/index.js"]
