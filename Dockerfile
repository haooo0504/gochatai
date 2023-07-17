# 開始於 Golang 基礎映像
FROM golang:latest AS builder

# 為我們的應用創建一個目錄
WORKDIR /gochatai

COPY config/app.yml /config/

# 將 go.mod 和 go.sum 文件複製到工作目錄
COPY go.mod go.sum ./

# 下載所有依賴項
RUN go mod download

# 將源碼複製到工作目錄
COPY . .

# 建置應用程式
RUN CGO_ENABLED=0 go build -o main .

# 第二階段：只需要一個最小化的基本映像即可
FROM alpine:latest

WORKDIR /root/

# 從編譯器階段複製可執行檔到我們的最小映像
COPY --from=builder /gochatai/main .

EXPOSE 8082

# 運行應用程式
CMD ["./main"]
