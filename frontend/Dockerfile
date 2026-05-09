# 构建阶段：安装依赖并打包 Vue 应用
FROM node:20-alpine AS builder

WORKDIR /app

# 优先复制锁文件，充分利用层缓存
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# 运行阶段：用 nginx 托管静态产物
FROM nginx:1.27-alpine

# 覆盖默认 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
