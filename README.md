# { Key -> Value } Storage

## What is it?

Сервер, реализующий INSERT, UPDATE, SELECT и DELETE операции на хранилище типа STRING -> STRING. И клиент для работы с ним


## Start server

```
port := 1234
rwTimeoutSec := 10
storage := kvstorage.NewStorage()
server, _ := kvstorage.NewServer(port, rwTimeoutSec, storage)
server.Start()
```
## Insert example

```
client := kvstorage.NewClient(serverAddr, timeoutSec)
client.Insert("foo", "bar")
res, _ := client.Select("foo")
fmt.Println(res)
```

## Save and load

```
data, _ := storage.Dump()
fmt.Println(string(data))

newStorage := kvstorage.NewStorage()
newStorage.Load(data)
```
