what i am learning through this ?

1. what's a database connection pool ? what does it do ?
what's the right size of pool to have in application based on number of requests it's processing ?
what would hapen if i have just one connection to db ?
ANS - https://chat.deepseek.com/a/chat/s/9bce0279-1e6b-4564-b9d5-82f9b7f048eb

2. how to insert a payload into DB ?
    refer to this commit https://github.com/memoryStack/library-management-system/commit/de7d7b0f8f31b8f0f1341a3e953af2f554cd197d
    - for json fields, we have to serialize them. else we will get errors like
        "SQLSTATE 42804 here means: Postgres column is jsonb, but GORM tried to insert related_images as a non-JSON value (record), so Postgres rejected it."
    Resources to read
        GORM models and field tags:
        https://gorm.io/docs/models.html

        GORM serializer tag (serializer:json):
        https://gorm.io/docs/serializer.html

        GORM customized data types (incl. JSON-oriented approaches):
        https://gorm.io/docs/data_types.html

        GORM + Postgres driver docs:
        https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL

3. why GORM soft deletes entries by default ?
    in this exercise i deleted a book and noticed that it's deleted_at column is updated, that's it.
    i checked in PgAdmin that the entry is still present. so below is the answer for the WHYs & HOWs.
    https://chat.deepseek.com/a/chat/s/8c1e69a6-7822-4e9a-8286-b648f4ee6656