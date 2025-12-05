# Build stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS builder

# Set build-time environment
ENV NODE_ENV=production

WORKDIR /app

# Copy package files first for better layer caching
COPY package*.json ./

# Install dependencies (including dev dependencies for build)
RUN npm ci --include=dev && \
    npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM --platform=linux/amd64 node:24.6.0-alpine3.22 AS production

# Install dumb-init for proper signal handling
RUN apk add --no-cache dumb-init

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nuxt -u 1001

WORKDIR /app

# Copy built application from builder stage
COPY --from=builder --chown=nuxt:nodejs /app/.output ./.output
COPY --from=builder --chown=nuxt:nodejs /app/package*.json ./

# Install only production dependencies
RUN npm ci --omit=dev && \
    npm cache clean --force

# Switch to non-root user
USER nuxt

# Expose the port
EXPOSE 3000

# Set production environment
ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Start the application
CMD ["node", ".output/server/index.mjs"]