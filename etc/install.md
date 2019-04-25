## Развертывание БД

Запускаем постгрес, копипастим в терминал:
```bash
sudo -u postgres psql -a -f db/migrations/0.sql
```

Или, если авторизация через ОС выключена (тогда надо ввести пароль):
```bash
psql -U postgres -a -f db/migrations/0.sql
```

Затем нужно добавить содержимое конфига из conf/pg_hba.conf в /etc/postgresql/(version)/main/pg_hba.conf и перезапустить Postgres:
```bash
sudo cat misc/conf/pg_hba.conf >> /etc/postgresql/10/main/pg_hba.conf
# или для Postgresql 9.*
sudo cat misc/conf/pg_hba.conf >> /etc/postgresql/9/main/pg_hba.conf
# И затем
sudo service postgresql restart
```
