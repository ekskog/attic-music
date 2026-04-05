FROM node:18-alpine
WORKDIR /app
COPY package.json package-lock.json* ./
# If package-lock.json is missing, `npm ci` fails; use `npm install` to be robust.
RUN npm install --omit=dev --no-audit --no-fund
COPY . .
EXPOSE 80
CMD ["node", "server/index.js"]
