# Линтер для архитектуры проекта и ее визуализатор

Установка

```shell
go install github.com/gbh007/goarchlint/cmd/goarchlint@latest
```

Генерация документации в локальной директории

```shell
goarchlint generate
```

Генерация документации вне локальной директории

```shell
goarchlint generate -p ~/projects/hgraber/hgraber-next -o ~/projects/hgraber/hgraber-next/docs/arch
```
