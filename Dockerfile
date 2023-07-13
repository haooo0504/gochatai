# 開始於 Golang 基礎映像
FROM golang:latest

# 為我們的應用創建一個目錄
WORKDIR /gochatai

# 將 go.mod 和 go.sum 文件複製到工作目錄
COPY go.mod go.sum ./

# 下載所有依賴項
RUN go mod download

# 將源碼複製到工作目錄
COPY . .

# 建置應用程式
RUN go build -o main .

# 指定端口號
EXPOSE 8082

# 運行應用程式
CMD ["./main"]
