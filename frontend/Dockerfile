FROM node:24-alpine

WORKDIR /app
COPY package.json package.json
COPY tsconfig.json tsconfig.json
COPY next.config.ts next.config.ts
COPY postcss.config.js postcss.config.js
COPY tailwind.config.ts tailwind.config.ts
COPY app app
COPY components components
COPY lib lib
COPY types types
COPY public public

RUN npm install
RUN npm run build

EXPOSE 3002
CMD ["npm", "run", "start"]
